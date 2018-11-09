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
		trustorPublicKey := DestAddr
		asset, err := xdrBuilder.SetAsset("ABC", "GAEBJVQJJO5ZPRJ2ZPNSDJLMNN64REZO7S5VUZAMNLI34B5XUQVD3URR")
		Expect(err).NotTo(HaveOccurred())
		authorize := true
		opB64, err = xdrBuilder.AllowTrust(SourceAddr, trustorPublicKey, asset, authorize)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return a correct xdr transaction string", func() {
		By("adding one allow trust operation")
		var (
			tB64 string
		)
		transactionBuilder := xdrSigner.GetBuilderInstance(&tx)
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		expected := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAQAAAABbADzR8HMtmT2H7RBKCPnQD69GQTfr7l1CjxCnKtpprQAAAAcAAAAAdfeOkY7szblb3J8lDy2i0o1ssnDcDOkFjxjwFx/sV+gAAAABQUJDAAAAAAEAAAAA"
		Expect(err).NotTo(HaveOccurred())
		Expect(tB64).Should(Equal(expected))
	})

	It("should return a correct unmarshalled bytes and operation", func() {
		By("adding one allow trust operation")
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
