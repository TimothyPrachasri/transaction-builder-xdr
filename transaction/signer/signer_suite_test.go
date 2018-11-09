package xdrsigner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"

	"testing"
)

func TestXdr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Signer Suite")
}

var (
	SourceSeed     string
	SourceAddr     string
	DestSeed       string
	DestAddr       string
	PassPhrase     string
	Skp            keypair.KP
	Dkp            keypair.KP
	Source         xdr.AccountId
	Memo           xdr.Memo
	err            error
	tx             xdr.Transaction
	DefaultBaseFee int
)

var _ = BeforeSuite(func() {
	SourceSeed = "SBMRYER7UW2KHXAUKFZC6YVRDX6TLBUL55RE37N7CNAGVRSVRW5GA6I4"
	SourceAddr = "GBNQAPGR6BZS3GJ5Q7WRASQI7HIA7L2GIE36X3S5IKHRBJZK3JU22QOK"
	DestSeed = "SAQSP6VPR5VY3SQTIXFHL7DISWSBG4ISXGR4HBTVSI2DIFCFMIHWGLNY"
	DestAddr = "GB27PDURR3WM3OK33SPSKDZNULJI23FSODOAZ2IFR4MPAFY75RL6Q6QR"
	Skp = keypair.MustParse(SourceSeed)
	Dkp = keypair.MustParse(DestAddr)
	PassPhrase = "Test SDF Network ; September 2015"
	Memo, err = xdr.NewMemo(xdr.MemoTypeMemoNone, nil)
	DefaultBaseFee = 100
	err = Source.SetAddress(Skp.Address())
	Expect(err).NotTo(HaveOccurred())
	err = Source.SetAddress(Skp.Address())
	Expect(err).NotTo(HaveOccurred())
	tx = xdr.Transaction{
		SourceAccount: Source,
		Fee:           10,
		SeqNum:        xdr.SequenceNumber(1),
		Memo:          Memo,
	}
})
