package repo

import (
	"os"

	"git.s8k.top/SeraphJACK/beanbot/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const Path = "./beancount"

var repo *git.Repository

func Init() error {
	if _, err := os.Stat(Path); err != nil {
		// repo does not exist, clone it
		if err := clone(); err != nil {
			return err
		}
	} else {
		// repo exists, just open it
		repo, err = git.PlainOpen(Path)
		if err != nil {
			return err
		}
	}
	return pull()
}

func clone() error {
	var err error
	repo, err = git.PlainClone(Path, false, &git.CloneOptions{
		URL: config.Cfg.Repo,
		Auth: &http.BasicAuth{
			Username: config.Cfg.Username,
			Password: config.Cfg.Password,
		},
	})
	return err
}

func pull() error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	return w.Pull(&git.PullOptions{Auth: &http.BasicAuth{
		Username: config.Cfg.Username,
		Password: config.Cfg.Password,
	}})
}
