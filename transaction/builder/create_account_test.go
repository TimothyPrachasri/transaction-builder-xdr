package transactionbuilder_test

import (
	"encoding/base64"
	"strings"
	builder "transaction-builder-xdr/transaction/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating transaction XDR with payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		destinationPublicKey := "GDKV36XRERL7HVQ5GKRAV47ZLEPIZMFM7MMLEO4NKQOWFPL5NCIEW3GR"
		var startingBalance uint64
		startingBalance = 100
		var destination xdr.AccountId
		err = destination.SetAddress(destinationPublicKey)
		Expect(err).NotTo(HaveOccurred())
		body := xdr.CreateAccountOp{
			Destination:     destination,
			StartingBalance: xdr.Int64(startingBalance) * 10000000,
		}
		op := xdr.Operation{}
		op.Body, err = xdr.NewOperationBody(xdr.OperationTypeCreateAccount, body)
		Expect(err).NotTo(HaveOccurred())
		opB64, err = xdr.MarshalBase64(op)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("Adding One Payment Operation")
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
		Expect(tB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA1V368SRX89YdMqIK8/lZHoywrPsYsjuNVB1ivX1okEsAAAAAO5rKAAAAAAA="))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("Adding One Payment Operation")
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
		Expect(bytesRead).Should(Equal(116))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))
	})
})
