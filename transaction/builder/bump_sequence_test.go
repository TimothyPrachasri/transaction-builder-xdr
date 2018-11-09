package transactionbuilder_test

import (
	"encoding/base64"
	"strings"

	builder "github.com/TimothyPrachasri/transaction-builder-xdr/transaction/builder"

	xdrBuilder "gitlab.com/lightnet-thailand/xdr-builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating transaction XDR with payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		var sequenceNumber uint64
		sequenceNumber = 1008097543847999
		opB64, err = xdrBuilder.BumpSequence(SourceAddr, sequenceNumber)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("adding one build sequence operation")
		var (
			tB64 string
		)
		transactionBuilder := builder.GetInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		expected := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAQAAAABbADzR8HMtmT2H7RBKCPnQD69GQTfr7l1CjxCnKtpprQAAAAsAA5TcAAAAPwAAAAA="
		Expect(err).NotTo(HaveOccurred())
		Expect(tB64).Should(Equal(expected))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("adding one build sequence operation")
		var (
			tB64           string
			unmarshalledTx xdr.Transaction
			bytesRead      int
		)
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
