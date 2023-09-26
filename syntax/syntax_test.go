package syntax

import (
	"fmt"
	"strings"
	"testing"
)

var Cfg = &Config{
	Currencies: []string{"CNY", "USD"},
	Accounts: map[string]string{
		"zfb":   "Assets::Digital::Alipay",
		"wx":    "Assets::Digital::Wechat",
		"dt":    "Expenses::Travel::Train",
		"lunch": "Expenses::Food::Lunch",
	},
	DefaultCurrency: "CNY",
}

func TestParse(t *testing.T) {
	for _, str := range []string{
		"地铁 dt 3 zfb -1 wx",
		"午饭 lunch 11.20 wx",
		"2023-01-01 转账 zfb 100 wx",
	} {
		t.Run(str, func(t *testing.T) {
			txn, err := Parse(strings.Split(str, " "), Cfg)
			if err != nil {
				t.Fatalf("%v", err)
			}
			fmt.Print(txn.ToBeanLanguageSyntax())
		})
	}
}
