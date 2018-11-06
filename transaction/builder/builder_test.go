package transactionbuilder_test

import (
	"encoding/base64"
	"strings"
	builder "transaction-builder-xdr/transaction/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating Transaction XDR", func() {
	var transactionBuilder builder.TransactionBuilder
	var sourceSeed string
	var destAddr string
	var passPhrase string
	var skp keypair.KP
	var dkp keypair.KP
	BeforeEach(func() {
		sourceSeed = "SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME"
		destAddr = "GCICVEBF5JYDBCTR3TXFGN56WGYBAKKWVHUQYPM72F6ZEQ7BDQZT4NFZ"
		skp = keypair.MustParse(sourceSeed)
		dkp = keypair.MustParse(destAddr)
		passPhrase = "Test SDF Network ; September 2015"
	})
	JustBeforeEach(func() {
		memo, _ := xdr.NewMemo(xdr.MemoTypeMemoNone, nil)
		var source xdr.AccountId
		_ = source.SetAddress(skp.Address())
		tx := xdr.Transaction{
			SourceAccount: source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          memo,
		}
		transactionBuilder = builder.GetInstance(&tx)
	})

	By("Adding Payment Operation")

	asset, _ := xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
	var destination xdr.AccountId
	_ = destination.SetAddress(dkp.Address())
	op := xdr.PaymentOp{
		Destination: destination,
		Asset:       asset,
		Amount:      50 * 10000000,
	}
	transactionBuilder.MakeOperation(xdr.OperationTypePayment, op)
	_ = transactionBuilder.Sign(sourceSeed, passPhrase)
	txeB64, _ := transactionBuilder.ToBase64()
	rawr := strings.NewReader(txeB64)
	b64r := base64.NewDecoder(base64.StdEncoding, rawr)

	var tx xdr.TransactionEnvelope
	bytesRead, _ := xdr.Unmarshal(b64r, &tx)

	It("should return a correct xdr enveloped transaction", func() {
		Expect(txeB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAkCqQJepwMIpx3O5TN76xsBApVqnpDD2f0X2SQ+EcMz4AAAAAAAAAAB3NZQAAAAAAAAAAAan/sYsAAABA3LBOukhEwmXdSSi7zTUC5MfZM8opRWI0SZ8cscGO/av0aJXfZzRtGXiWE2cK7voF1GswJUiEpNo/EEWHYvzjCA=="))
	})

	It("should unmarshalled to a correct bytes and number of operations", func() {
		Expect(bytesRead).Should(Equal(196))
		Expect(len(tx.Tx.Operations)).Should(Equal(1))
	})
})
