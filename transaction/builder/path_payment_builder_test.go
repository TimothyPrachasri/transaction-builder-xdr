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

var _ = Describe("Creating transaction XDR with path payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		sendAsset, err := xdrBuilder.SetAsset("ABC", "GAEBJVQJJO5ZPRJ2ZPNSDJLMNN64REZO7S5VUZAMNLI34B5XUQVD3URR")
		Expect(err).NotTo(HaveOccurred())
		var sendMax uint64
		sendMax = 100
		destinationPublicKey := "GCIQJ3JRXEEAKFL22C43X66B4NKACPWZ27WIMNXGA5CIEHOYWNXD3EQR"
		destAsset, err := xdrBuilder.SetAsset("CDF", "GAEBJVQJJO5ZPRJ2ZPNSDJLMNN64REZO7S5VUZAMNLI34B5XUQVD3URR")
		Expect(err).NotTo(HaveOccurred())
		var destAmount uint64
		destAmount = 30
		tempAsset1, err := xdrBuilder.SetNativeAsset()
		Expect(err).NotTo(HaveOccurred())
		var path xdrBuilder.Path
		path.XDRAsset = append(path.XDRAsset, tempAsset1.XDRAsset)
		var destination xdr.AccountId
		err = destination.SetAddress(destinationPublicKey)
		Expect(err).NotTo(HaveOccurred())
		body := xdr.PathPaymentOp{
			SendAsset:   sendAsset.XDRAsset,
			SendMax:     xdr.Int64(sendMax) * 10000000,
			Destination: destination,
			DestAsset:   destAsset.XDRAsset,
			DestAmount:  xdr.Int64(destAmount) * 10000000,
			Path:        path.XDRAsset,
		}
		op := xdr.Operation{}
		op.Body, err = xdr.NewOperationBody(xdr.OperationTypePathPayment, body)
		Expect(err).NotTo(HaveOccurred())
		opB64, err = xdr.MarshalBase64(op)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("adding One Payment Operation")
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
		Expect(tB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAIAAAABQUJDAAAAAAAIFNYJS7uXxTrL2yGlbGt9yJMu/LtaZAxq0b4Ht6QqPQAAAAA7msoAAAAAAJEE7TG5CAUVetC5u/vB41QBPtnX7IY25gdEgh3Ys249AAAAAUNERgAAAAAACBTWCUu7l8U6y9shpWxrfciTLvy7WmQMatG+B7ekKj0AAAAAEeGjAAAAAAEAAAAAAAAAAA=="))

		By("adding two or more operation")
		tx = xdr.Transaction{
			SourceAccount: Source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          Memo,
		}
		transactionBuilder = builder.GetInstance(&tx)
		operationsArr := []string{opB64, opB64, opB64}
		transactionBuilder.MakeAllOperations(operationsArr)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())
		Expect(tB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAADAAAAAAAAAAIAAAABQUJDAAAAAAAIFNYJS7uXxTrL2yGlbGt9yJMu/LtaZAxq0b4Ht6QqPQAAAAA7msoAAAAAAJEE7TG5CAUVetC5u/vB41QBPtnX7IY25gdEgh3Ys249AAAAAUNERgAAAAAACBTWCUu7l8U6y9shpWxrfciTLvy7WmQMatG+B7ekKj0AAAAAEeGjAAAAAAEAAAAAAAAAAAAAAAIAAAABQUJDAAAAAAAIFNYJS7uXxTrL2yGlbGt9yJMu/LtaZAxq0b4Ht6QqPQAAAAA7msoAAAAAAJEE7TG5CAUVetC5u/vB41QBPtnX7IY25gdEgh3Ys249AAAAAUNERgAAAAAACBTWCUu7l8U6y9shpWxrfciTLvy7WmQMatG+B7ekKj0AAAAAEeGjAAAAAAEAAAAAAAAAAAAAAAIAAAABQUJDAAAAAAAIFNYJS7uXxTrL2yGlbGt9yJMu/LtaZAxq0b4Ht6QqPQAAAAA7msoAAAAAAJEE7TG5CAUVetC5u/vB41QBPtnX7IY25gdEgh3Ys249AAAAAUNERgAAAAAACBTWCUu7l8U6y9shpWxrfciTLvy7WmQMatG+B7ekKj0AAAAAEeGjAAAAAAEAAAAAAAAAAA=="))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("adding one Operation")
		var (
			tB64               string
			unmarshalledTx     xdr.Transaction
			bytesRead          int
			transactionBuilder builder.TransactionBuilder
		)
		tx := xdr.Transaction{
			SourceAccount: Source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          Memo,
		}
		transactionBuilder = builder.GetInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())
		rawr := strings.NewReader(tB64)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		bytesRead, err = xdr.Unmarshal(b64r, &unmarshalledTx)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesRead).Should(Equal(220))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))

		By("adding two or more Operation")
		tx = xdr.Transaction{
			SourceAccount: Source,
			Fee:           10,
			SeqNum:        xdr.SequenceNumber(1),
			Memo:          Memo,
		}
		transactionBuilder = builder.GetInstance(&tx)
		operationsArr := []string{opB64, opB64, opB64}
		transactionBuilder.MakeAllOperations(operationsArr)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())
		rawr = strings.NewReader(tB64)
		b64r = base64.NewDecoder(base64.StdEncoding, rawr)
		bytesRead, err = xdr.Unmarshal(b64r, &unmarshalledTx)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesRead).Should(Equal(532))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(3))
	})
})
