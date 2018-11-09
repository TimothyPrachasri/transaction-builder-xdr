package xdrsigner_test

import (
	"strconv"

	xdrSigner "github.com/TimothyPrachasri/transaction-builder-xdr/transaction/signer"

	"github.com/stellar/go/xdr"
	xdrBuilder "gitlab.com/lightnet-thailand/xdr-builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/clients/horizon"
)

var _ = Describe("Sending enveloped transaction to Horizon server", func() {
	var (
		opB64    string
		opB64Arr []string
	)

	BeforeEach(func() {
		var asset xdrBuilder.Asset
		opB64, err = xdrBuilder.Payment(SourceAddr, DestAddr, asset, "0.01")
		for i := 0; i <= 99; i++ {
			opB64Arr = append(opB64Arr, opB64)
		}
		Expect(err).NotTo(HaveOccurred())
	})

	It("should successfully sent to Horizon", func() {
		var (
			transactionBuilder  xdrSigner.TransactionBuilder
			tB64                string
			txeB64              string
			transactionEnvelope xdrSigner.TransactionEnvelope
		)
		By("Adding Single Operation")
		acct, err := horizon.DefaultTestNetClient.LoadAccount(SourceAddr)
		Expect(err).NotTo(HaveOccurred())
		sequenceNum, err := strconv.ParseInt(acct.Sequence, 10, 64)
		Expect(err).NotTo(HaveOccurred())
		transactionBuilder = xdrSigner.GetBuilderInstance(&tx)
		transactionBuilder.MakeSequenceNumber(uint64(sequenceNum + 1))
		transactionBuilder.MakeFee(xdr.Uint32(DefaultBaseFee))
		transactionBuilder.MakeOperation(opB64)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())

		transactionEnvelope = xdrSigner.GetEnveloperInstance(tB64)
		err = transactionEnvelope.Sign(SourceSeed, PassPhrase)
		Expect(err).NotTo(HaveOccurred())
		txeB64, err = transactionEnvelope.ToBase64()

		_, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
		Expect(err).NotTo(HaveOccurred())

		By("Adding 100 Operations")
		acct, err = horizon.DefaultTestNetClient.LoadAccount(SourceAddr)
		Expect(err).NotTo(HaveOccurred())
		sequenceNum, err = strconv.ParseInt(acct.Sequence, 10, 64)
		Expect(err).NotTo(HaveOccurred())
		transactionBuilder = xdrSigner.GetBuilderInstance(&tx)
		transactionBuilder.MakeSequenceNumber(uint64(sequenceNum + 1))
		transactionBuilder.MakeFee(xdr.Uint32(int(100) * 100))
		transactionBuilder.MakeAllOperations(opB64Arr)
		tB64, err = transactionBuilder.ToBase64()
		Expect(err).NotTo(HaveOccurred())

		transactionEnvelope = xdrSigner.GetEnveloperInstance(tB64)
		err = transactionEnvelope.Sign(SourceSeed, PassPhrase)
		Expect(err).NotTo(HaveOccurred())
		txeB64, err = transactionEnvelope.ToBase64()
		_, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
		Expect(err).NotTo(HaveOccurred())
	})
})
