package transactionbuilder_test

import (
	"encoding/base64"
	"strings"
	builder "transaction-builder-xdr/transaction/builder"

	xdrBuilder "github.com/Kafakk/xdr-builder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating transaction XDR with payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		asset, err := xdrBuilder.SetAsset("ABC", "GAEBJVQJJO5ZPRJ2ZPNSDJLMNN64REZO7S5VUZAMNLI34B5XUQVD3URR")
		if err != nil {
			panic(err)
		}
		var limit uint64
		limit = 1000
		body := xdr.ChangeTrustOp{
			Line:  asset.XDRAsset,
			Limit: xdr.Int64(limit),
		}
		op := xdr.Operation{}
		op.Body, err = xdr.NewOperationBody(xdr.OperationTypeChangeTrust, body)
		Expect(err).NotTo(HaveOccurred())
		opB64, err = xdr.MarshalBase64(op)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("adding one change trust operation")
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
		Expect(tB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABQUJDAAAAAAAIFNYJS7uXxTrL2yGlbGt9yJMu/LtaZAxq0b4Ht6QqPQAAAAAAAAPoAAAAAA=="))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("adding one change trust operation")
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
		Expect(bytesRead).Should(Equal(124))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))
	})
})
