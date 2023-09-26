package repo

import (
	"bytes"
	"os"
	"path"
	"time"

	"git.s8k.top/SeraphJACK/beanbot/config"
	"git.s8k.top/SeraphJACK/beanbot/syntax"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CommitTransaction(txn *syntax.Transaction) error {
	beanPath, err := calcPathForTxn(txn)
	if err != nil {
		return err
	}

	fullBeanPath := path.Join(Path, beanPath)

	// Write txn bean language to bean file
	f, err := os.OpenFile(fullBeanPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(txn.ToBeanLanguageSyntax())
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	// Commit bean txn
	// First pull the repository so *hopefully* we won't conflict with someone else
	if err := pull(); err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Stage the changes
	_, err = w.Add(beanPath)
	if err != nil {
		return err
	}

	// Commit the staged changes
	_, err = w.Commit("Beanbot auto commit txn "+time.Now().Format("2006-01-02 15:04:05"), &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Cfg.CommitAuthor,
			Email: config.Cfg.CommitAuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// Push changes to the remote repository
	return repo.Push(&git.PushOptions{Auth: &http.BasicAuth{
		Username: config.Cfg.Username,
		Password: config.Cfg.Password,
	}})
}

func calcPathForTxn(txn *syntax.Transaction) (string, error) {
	buf := &bytes.Buffer{}
	err := config.TxnBeanPathTmpl.Execute(buf, map[string]string{
		"YYYY": txn.Date[0:4],
		"MM":   txn.Date[5:7],
		"DD":   txn.Date[8:10],
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
