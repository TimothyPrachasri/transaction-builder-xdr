package transactionbuilder_test

import (
	builder "transaction-builder-xdr/transaction/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Creating Transaction", func() {

	It("should return a correct xdr transaction fields", func() {
		var (
			transactionBuilder builder.TransactionBuilder
			accountId          xdr.AccountId
		)
		By("adding source account")
		transactionBuilder = builder.GetInstance()
		transactionBuilder.MakeSourceAccount(SourceSeed)
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.SourceAccount).Should(BeAssignableToTypeOf(accountId))

		By("adding fee")
		transactionBuilder = builder.GetInstance()
		transactionBuilder.MakeFee(100)
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.Fee).Should(BeEquivalentTo(100))

		By("adding sequence number")
		transactionBuilder = builder.GetInstance()
		transactionBuilder.MakeSequenceNumber(uint64(1))
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.SeqNum).Should(BeEquivalentTo(1))

		By("adding memo")
		transactionBuilder = builder.GetInstance()
		transactionBuilder.MakeMemo(xdr.MemoTypeMemoNone, nil)
		accountId = xdr.AccountId{}
		Expect(transactionBuilder.TransactionXDR.Memo.Type).Should(BeEquivalentTo(xdr.MemoTypeMemoNone))
	})
})
