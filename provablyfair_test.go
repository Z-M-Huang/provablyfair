package provablyfair

import (
	"encoding/hex"
	"testing"
)

var PFClient *Client

func TestMain(m *testing.M) {
	serverSeed, _ := hex.DecodeString("1676a2e695dd49480bfef863b608545afdbf84d3ae2bc06eafc3afb368c6a114")
	PFClient = &Client{
		ServerSeed: serverSeed,
	}
}

func TestGenerate(t *testing.T) {
	num, _, _, _ := PFClient.GenerateFromString("dfdbbe2fa6a17e076fc4096d6193b8921a015bf99e7a41957646bcef58729472")

	if num != 90.44 {
		t.Errorf("Expecting 90.44, got %f", num)
	}
}

func TestVerify(t *testing.T) {
	pass, err := VerifyFromString("dfdbbe2fa6a17e076fc4096d6193b8921a015bf99e7a41957646bcef58729472",
		"1676a2e695dd49480bfef863b608545afdbf84d3ae2bc06eafc3afb368c6a114", 1, 90441)
	if err != nil {
		t.Error(err)
	} else if !pass {
		t.Error("Verify failed")
	}
}
