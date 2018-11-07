package transactionbuilder_test

import (
	"fmt"
	builder "transaction-builder-xdr/transaction/builder"
	enveloper "transaction-builder-xdr/transaction/envelope"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

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

var _ = Describe("Creating Transaction XDR", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		asset, err := xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
		Expect(err).NotTo(HaveOccurred())
		var destination xdr.AccountId
		err = destination.SetAddress(Dkp.Address())
		Expect(err).NotTo(HaveOccurred())
		body := xdr.PaymentOp{
			Destination: destination,
			Asset:       asset,
			Amount:      50 * 10000000,
		}
		op := xdr.Operation{}
		op.Body, err = xdr.NewOperationBody(xdr.OperationTypePayment, body)
		if err != nil {
			panic(err)
		}
		opB64, err = xdr.MarshalBase64(op)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr enveloped transaction", func() {
		By("Adding Payment Operation")
		var (
			tB64   string
			txeB64 string
		)
		tx := xdr.Transaction{
			SourceAccount: Source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          Memo,
		}
		fmt.Println(opB64, "base 64 bitch")
		transactionBuilder := builder.GetInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		fmt.Println(transactionBuilder.TransactionXDR)
		tB64, err = transactionBuilder.ToBase64()
		fmt.Println("pass0")
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("pass1")
		transactionEnvelope := enveloper.GetInstance(tB64)
		err = transactionEnvelope.Sign("SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME", "Test SDF Network ; September 2015")
		fmt.Println("pass2")
		Expect(err).NotTo(HaveOccurred())
		txeB64, err = transactionEnvelope.ToBase64()
		Expect(txeB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAkCqQJepwMIpx3O5TN76xsBApVqnpDD2f0X2SQ+EcMz4AAAAAAAAAAB3NZQAAAAAAAAAAAan/sYsAAABA3LBOukhEwmXdSSi7zTUC5MfZM8opRWI0SZ8cscGO/av0aJXfZzRtGXiWE2cK7voF1GswJUiEpNo/EEWHYvzjCA=="))
	})
})
