package bind

import (
	"fmt"
	"testing"

	"github.com/evan-forbes/ethq/contracts/erc20/dai"
)

func TestBind(t *testing.T) {
	code, err := Bind([]string{"dai"}, []string{dai.DaiABI}, []string{dai.DaiBin}, "dai")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(code)
}
