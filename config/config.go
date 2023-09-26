package config

import (
	"os"
	"text/template"

	"git.s8k.top/SeraphJACK/beanbot/syntax"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Repo               string        `yaml:"repo"`
	Username           string        `yaml:"username"`
	Password           string        `yaml:"password"`
	CommitAuthor       string        `yaml:"commitAuthor"`
	CommitAuthorEmail  string        `yaml:"commitAuthorEmail"`
	Syntax             syntax.Config `yaml:"syntax"`
	TxnBeanPathPattern string        `yaml:"txnBeanPathPattern"`
}

var Cfg = Config{
	Repo:              "https://git.xxx.com/xxx/xxx.git",
	Username:          "someone",
	Password:          "password",
	CommitAuthor:      "Beanbot",
	CommitAuthorEmail: "beanbot@example.com",
	Syntax: syntax.Config{
		Currencies:      []string{"CNY", "USD"},
		Accounts:        map[string]string{"zfb": "Assets::Digital::Alipay"},
		DefaultCurrency: "CNY",
	},
	TxnBeanPathPattern: "txs/{{.YYYY}}-{{.MM}}.bean",
}

var TxnBeanPathTmpl *template.Template

func Load(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&Cfg); err != nil {
		return err
	}
	TxnBeanPathTmpl, err = template.New("txn_bean_path").Parse(Cfg.TxnBeanPathPattern)
	return err
}
