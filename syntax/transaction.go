package syntax

import (
	"bytes"
	_ "embed"
	"errors"
	"text/template"
)

//go:embed transaction.tmpl
var txnTmpl string

var parsedTmpl *template.Template

func init() {
	parsedTmpl = template.Must(template.New("bean-language-txn").Parse(txnTmpl))
}

type Posting struct {
	Account  string
	Amount   float64
	Currency string
	Flag     string
}

type Transaction struct {
	Date      string
	Payee     string
	Narration string
	Flag      string
	Postings  []*Posting
}

func (txn *Transaction) Validate() error {
	// currency: amount
	sum := map[string]float64{}
	zeros := map[string]int{}
	zerop := map[string]*Posting{}

	for _, p := range txn.Postings {
		sum[p.Currency] += p.Amount
		if p.Amount == 0 {
			zeros[p.Currency]++
			zerop[p.Currency] = p
		}
	}

	for cur, s := range sum {
		if zeros[cur] > 1 {
			return errors.New("multiple zero amount postings")
		}
		if zeros[cur] == 1 {
			zerop[cur].Amount = -s
		}
		if zeros[cur] == 0 && s != 0 {
			return errors.New("sum of amount is non-zero")
		}
	}

	return nil
}

func (txn *Transaction) ToBeanLanguageSyntax() string {
	buf := &bytes.Buffer{}
	err := parsedTmpl.Execute(buf, txn)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
