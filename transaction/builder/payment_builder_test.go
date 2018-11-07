package transactionbuilder_test

import (
	"encoding/base64"
	"strings"
	builder "transaction-builder-xdr/transaction/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

var (
	SourceSeed string
	DestAddr   string
	PassPhrase string
	Skp        keypair.KP
	Dkp        keypair.KP
	Source     xdr.AccountId
	Memo       xdr.Memo
	err        error
)

var _ = BeforeSuite(func() {
	SourceSeed = "SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME"
	DestAddr = "GCICVEBF5JYDBCTR3TXFGN56WGYBAKKWVHUQYPM72F6ZEQ7BDQZT4NFZ"
	Skp = keypair.MustParse(SourceSeed)
	Dkp = keypair.MustParse(DestAddr)
	PassPhrase = "Test SDF Network ; September 2015"
	Memo, err = xdr.NewMemo(xdr.MemoTypeMemoNone, nil)
	err = Source.SetAddress(Skp.Address())
	if err != nil {
		panic(err)
	}
})

var _ = Describe("Creating transaction XDR with payment operation", func() {
	var (
		opB64 string
	)

	BeforeEach(func() {
		asset, err := xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
		Expect(err).NotTo(HaveOccurred())
		var destination xdr.AccountId
		err = destination.SetAddress(Dkp.Address())
		Expect(err).NotTo(HaveOccurred())
		body := xdr.PaymentOp{
			Destination: destination,
			Asset:       asset,
			Amount:      50 * 10000000,
		}
		op := xdr.Operation{}
		op.Body, err = xdr.NewOperationBody(xdr.OperationTypePayment, body)
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
		Expect(tB64).Should(Equal("AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAkCqQJepwMIpx3O5TN76xsBApVqnpDD2f0X2SQ+EcMz4AAAAAAAAAAB3NZQAAAAAA"))
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
		Expect(bytesRead).Should(Equal(120))
		Expect(len(unmarshalledTx.Operations)).Should(Equal(1))
	})
})
