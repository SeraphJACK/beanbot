package syntax

import (
	"errors"
	"strconv"
	"time"
)

var ErrMissingPayee = errors.New("missing payee")
var ErrUnknownAccount = errors.New("unknown account")
var ErrEmptyTransaction = errors.New("no postings defined")

type Config struct {
	Currencies      []string          `yaml:"currencies"`
	Accounts        map[string]string `yaml:"accounts"`
	DefaultCurrency string            `yaml:"defaultCurrency"`
}

func (c *Config) hasCurrency(cur string) bool {
	for _, v := range c.Currencies {
		if v == cur {
			return true
		}
	}
	return false
}

func (c *Config) account(abbr string) (string, bool) {
	if res, ok := c.Accounts[abbr]; ok {
		return res, true
	}
	return "", false
}

type tokenReader struct {
	tokens []string
	cur    int
}

func (r *tokenReader) HasNext() bool {
	return r.cur < len(r.tokens)
}

func (r *tokenReader) Peek() string {
	return r.tokens[r.cur]
}

func (r *tokenReader) Next() string {
	res := r.tokens[r.cur]
	r.cur++
	return res
}

func Parse(tokens []string, cfg *Config) (*Transaction, error) {
	r := &tokenReader{tokens: tokens}
	txn := &Transaction{}

	// optional transaction date
	if r.HasNext() {
		t, err := time.Parse("2006-01-02", r.Peek())
		if err == nil {
			txn.Date = t.Format("2006-01-02")
			r.Next()
		} else {
			// default date: current date
			txn.Date = time.Now().Format("2006-01-02")
		}
	}

	// optional txn mark
	if r.HasNext() {
		if r.Peek() == "*" || r.Peek() == "!" {
			txn.Flag = r.Next()
		} else {
			// default flag is * (complete txn)
			txn.Flag = "*"
		}
	}

	// mandatory payee
	if !r.HasNext() {
		return nil, ErrMissingPayee
	}
	txn.Payee = r.Next()

	p, err := readPosting(r, cfg)
	if err != nil {
		// account does not exist, token is optional narration
		if r.HasNext() {
			txn.Narration = r.Next()
		}
	} else {
		txn.Postings = append(txn.Postings, p)
	}

	// remaining postings
	for r.HasNext() {
		p, err := readPosting(r, cfg)
		if err != nil {
			return nil, err
		}
		txn.Postings = append(txn.Postings, p)
	}

	if len(txn.Postings) == 0 {
		return nil, ErrEmptyTransaction
	}

	return txn, txn.Validate()
}

func readPosting(r *tokenReader, cfg *Config) (*Posting, error) {
	p := &Posting{}

	// mandatory account
	if !r.HasNext() {
		return nil, ErrUnknownAccount
	}
	account, ok := cfg.account(r.Peek())
	if !ok {
		return nil, ErrUnknownAccount
	}
	p.Account = account
	r.Next()

	// optional flag
	if r.HasNext() && (r.Peek() == "!" || r.Peek() == "*") {
		p.Flag = r.Next()
	} else {
		p.Flag = "*"
	}

	// optional amount
	if r.HasNext() {
		amount, err := strconv.ParseFloat(r.Peek(), 64)
		if err == nil {
			p.Amount = amount
			r.Next()
		}
	}

	// optional currency
	if r.HasNext() && cfg.hasCurrency(r.Peek()) {
		p.Currency = r.Next()
	} else {
		p.Currency = cfg.DefaultCurrency
	}

	return p, nil
}
