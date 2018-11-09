package xdrsigner_test

import (
	"encoding/base64"
	"strings"

	xdrSigner "github.com/TimothyPrachasri/transaction-builder-xdr/transaction/signer"

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
		var asset xdrBuilder.Asset
		opB64, err = xdrBuilder.Payment(SourceAddr, DestAddr, asset, "100")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("Adding One Payment Operation")
		var (
			tB64 string
		)
		expected := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAQAAAABbADzR8HMtmT2H7RBKCPnQD69GQTfr7l1CjxCnKtpprQAAAAEAAAAAdfeOkY7szblb3J8lDy2i0o1ssnDcDOkFjxjwFx/sV+gAAAAAAAAAADuaygAAAAAA"
		transactionBuilder := xdrSigner.GetBuilderInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())
		Expect(tB64).Should(Equal(expected))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("Adding One Payment Operation")
		var (
			tB64           string
			unmarshalledTx xdr.Transaction
			bytesRead      int
		)
		transactionBuilder := xdrSigner.GetBuilderInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		rawr := strings.NewReader(tB64)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		bytesRead, err = xdr.Unmarshal(b64r, &unmarshalledTx)
		Expect(bytesRead).Should(Equal(156))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))
	})
})
