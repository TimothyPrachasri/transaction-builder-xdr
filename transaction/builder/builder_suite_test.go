package transactionbuilder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"

	"testing"
)

func TestXdr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builder Suite")
}

var (
	SourceSeed string
	DestAddr   string
	PassPhrase string
	Skp        keypair.KP
	Dkp        keypair.KP
	Source     xdr.AccountId
	Memo       xdr.Memo
	err        error
)

var _ = BeforeSuite(func() {
	SourceSeed = "SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME"
	DestAddr = "GCICVEBF5JYDBCTR3TXFGN56WGYBAKKWVHUQYPM72F6ZEQ7BDQZT4NFZ"
	Skp = keypair.MustParse(SourceSeed)
	Dkp = keypair.MustParse(DestAddr)
	PassPhrase = "Test SDF Network ; September 2015"
	Memo, err = xdr.NewMemo(xdr.MemoTypeMemoNone, nil)
	err = Source.SetAddress(Skp.Address())
	if err != nil {
		panic(err)
	}
})
