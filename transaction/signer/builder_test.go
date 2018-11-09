package xdrsigner_test

import (
	xdrSigner "transaction-builder-xdr/transaction/signer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating Transaction", func() {

	It("should return a correct xdr transaction fields", func() {
		var (
			transactionBuilder xdrSigner.TransactionBuilder
			accountId          xdr.AccountId
		)
		By("adding source account")
		transactionBuilder = xdrSigner.GetBuilderInstance()
		transactionBuilder.MakeSourceAccount(SourceSeed)
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.SourceAccount).Should(BeAssignableToTypeOf(accountId))

		By("adding fee")
		transactionBuilder = xdrSigner.GetBuilderInstance()
		transactionBuilder.MakeFee(100)
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.Fee).Should(BeEquivalentTo(100))

		By("adding sequence number")
		transactionBuilder = xdrSigner.GetBuilderInstance()
		transactionBuilder.MakeSequenceNumber(uint64(1))
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.SeqNum).Should(BeEquivalentTo(1))

		By("adding memo")
		transactionBuilder = xdrSigner.GetBuilderInstance()
		transactionBuilder.MakeMemo(xdr.MemoTypeMemoNone, nil)
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.Memo.Type).Should(BeEquivalentTo(xdr.MemoTypeMemoNone))
	})
})
