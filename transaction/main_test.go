package test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	builder "transaction-builder-xdr/transaction/builder"
	envelope "transaction-builder-xdr/transaction/envelope"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

func ExampleUsingTransactionBuilder() {
	skp := keypair.MustParse("SBMRYER7UW2KHXAUKFZC6YVRDX6TLBUL55RE37N7CNAGVRSVRW5GA6I4")
	dkp := keypair.MustParse("GB27PDURR3WM3OK33SPSKDZNULJI23FSODOAZ2IFR4MPAFY75RL6Q6QR")

	asset, err := xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
	if err != nil {
		panic(err)
	}

	var destination xdr.AccountId
	err = destination.SetAddress(dkp.Address())
	if err != nil {
		panic(err)
	}

	body := xdr.PaymentOp{
		Destination: destination,
		Asset:       asset,
		Amount:      1000000000,
	}
	op := xdr.Operation{}
	op.Body, err = xdr.NewOperationBody(xdr.OperationTypePayment, body)
	if err != nil {
		panic(err)
	}

	memo, err := xdr.NewMemo(xdr.MemoTypeMemoNone, nil)

	var source xdr.AccountId
	err = source.SetAddress(skp.Address())
	if err != nil {
		panic(err)
	}
	var defaultBaseFee uint64 = 100
	tx := xdr.Transaction{
		SourceAccount: source,
		Fee:           xdr.Uint32(int(defaultBaseFee) * 1),
		SeqNum:        xdr.SequenceNumber(2517825793032198),
		Memo:          memo,
	}

	transactionBuilder := builder.GetInstance(&tx)
	opB64, err := xdr.MarshalBase64(op)
	if err != nil {
		panic(err)
	}
	transactionBuilder.MakeOperation(opB64)

	tB64, err := transactionBuilder.ToBase64()
	if err != nil {
		panic(err)
	}
	transactionEnvelope := envelope.GetInstance(tB64)
	err = transactionEnvelope.Sign("SBMRYER7UW2KHXAUKFZC6YVRDX6TLBUL55RE37N7CNAGVRSVRW5GA6I4", "Test SDF Network ; September 2015")
	if err != nil {
		panic(err)
	}
	res, _ := json.Marshal(transactionEnvelope)
	fmt.Println(string(res))
	txeB64, err := transactionEnvelope.ToBase64()
	if err != nil {
		panic(err)
	}
	fmt.Printf("tx base64: %s", txeB64)
	// Output: tx base64: AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAAZAAI8fMAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAdfeOkY7szblb3J8lDy2i0o1ssnDcDOkFjxjwFx/sV+gAAAAAAAAAAAAAAGQAAAAAAAAAASraaa0AAABAQVkwp5giAPb+it3dZlGMDzC7YmbaC3ESwUF1ud440EYHSTCsa2lW547GSTZMFfRgb2gOAmuKooD7lH4qYC0dAg==
}

func ExampleDecodeTransaction() {
	data := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAAZAAI8fMAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAdfeOkY7szblb3J8lDy2i0o1ssnDcDOkFjxjwFx/sV+gAAAAAAAAAAAAAAGQAAAAAAAAAASraaa0AAABAQVkwp5giAPb+it3dZlGMDzC7YmbaC3ESwUF1ud440EYHSTCsa2lW547GSTZMFfRgb2gOAmuKooD7lH4qYC0dAg=="
	rawr := strings.NewReader(data)
	b64r := base64.NewDecoder(base64.StdEncoding, rawr)

	var tx xdr.TransactionEnvelope
	bytesRead, err := xdr.Unmarshal(b64r, &tx)

	fmt.Printf("read %d bytes\n", bytesRead)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("This tx has %d operations\n", len(tx.Tx.Operations))
	// Output: read 196 bytes
	// This tx has 1 operations
}
