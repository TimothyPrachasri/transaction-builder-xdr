package transactionbuilder_test

import (
	"encoding/base64"
	"strings"
	builder "transaction-builder-xdr/transaction/builder"

	xdrBuilder "github.com/Kafakk/xdr-builder"
	"github.com/stellar/go/price"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating transaction XDR with payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		sellingAsset, err := xdrBuilder.SetAsset("ABC", "GAEBJVQJJO5ZPRJ2ZPNSDJLMNN64REZO7S5VUZAMNLI34B5XUQVD3URR")
		Expect(err).NotTo(HaveOccurred())
		buyingAsset, err := xdrBuilder.SetNativeAsset()
		Expect(err).NotTo(HaveOccurred())
		var amount uint64
		amount = 1000
		priceString := "7.2"
		price, err := price.Parse(priceString)
		Expect(err).NotTo(HaveOccurred())
		body := xdr.CreatePassiveOfferOp{
			Selling: sellingAsset.XDRAsset,
			Buying:  buyingAsset.XDRAsset,
			Amount:  xdr.Int64(amount) * 10000000,
			Price:   price,
		}
		op := xdr.Operation{}
		op.Body, err = xdr.NewOperationBody(xdr.OperationTypeCreatePassiveOffer, body)
		Expect(err).NotTo(HaveOccurred())
		opB64, err = xdr.MarshalBase64(op)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("adding one create passive offer operation")
		var (
			tB64 string
		)
		tx := xdr.Transaction{
			SourceAccount: Source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          Memo,
		}
		transactionBuilder := builder.GetInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())
		Expect(tB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAQAAAABQUJDAAAAAAAIFNYJS7uXxTrL2yGlbGt9yJMu/LtaZAxq0b4Ht6QqPQAAAAAAAAACVAvkAAAAACQAAAAFAAAAAA=="))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("adding one create passive offer operation")
		var (
			tB64           string
			unmarshalledTx xdr.Transaction
			bytesRead      int
		)
		tx := xdr.Transaction{
			SourceAccount: Source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          Memo,
		}
		transactionBuilder := builder.GetInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		rawr := strings.NewReader(tB64)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		bytesRead, err = xdr.Unmarshal(b64r, &unmarshalledTx)
		Expect(bytesRead).Should(Equal(136))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))
	})
})
