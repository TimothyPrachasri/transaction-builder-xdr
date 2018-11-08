package transactionbuilder_test

import (
	"encoding/base64"
	"strings"
	builder "transaction-builder-xdr/transaction/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
	xdrBuilder "gitlab.com/lightnet-thailand/xdr-builder"
)

var _ = Describe("Creating transaction XDR with payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		inflationDestinationPublicKey := "GCIQJ3JRXEEAKFL22C43X66B4NKACPWZ27WIMNXGA5CIEHOYWNXD3EQR"
		opB64, err = xdrBuilder.SetOptionInflation(SourceAddr, inflationDestinationPublicKey)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("adding one set options operation")
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
		expected := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAQAAAABbADzR8HMtmT2H7RBKCPnQD69GQTfr7l1CjxCnKtpprQAAAAUAAAABAAAAAJEE7TG5CAUVetC5u/vB41QBPtnX7IY25gdEgh3Ys249AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
		Expect(err).NotTo(HaveOccurred())
		Expect(tB64).Should(Equal(expected))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("adding one set options operation")
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
		Expect(bytesRead).Should(Equal(180))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))
	})
})
