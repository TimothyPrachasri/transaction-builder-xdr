package test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	builder "transaction-builder-xdr/transaction"

	"github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// ExampleBuildTransaction creates and signs a simple transaction using the
// build package. The build package is designed to make it easier and more
// intuitive to configure and sign a transaction.
func ExampleLowLevelTransaction() {
	skp := keypair.MustParse("SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME")
	dkp := keypair.MustParse("GCICVEBF5JYDBCTR3TXFGN56WGYBAKKWVHUQYPM72F6ZEQ7BDQZT4NFZ")

	asset, err := xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
	if err != nil {
		panic(err)
	}

	var destination xdr.AccountId
	err = destination.SetAddress(dkp.Address())
	if err != nil {
		panic(err)
	}

	op := xdr.PaymentOp{
		Destination: destination,
		Asset:       asset,
		Amount:      50 * 10000000,
	}

	memo, err := xdr.NewMemo(xdr.MemoTypeMemoNone, nil)

	var source xdr.AccountId
	err = source.SetAddress(skp.Address())
	if err != nil {
		panic(err)
	}

	body, err := xdr.NewOperationBody(xdr.OperationTypePayment, op)
	if err != nil {
		panic(err)
	}

	tx := xdr.Transaction{
		SourceAccount: source,
		Fee:           10,
		SeqNum:        xdr.SequenceNumber(1),
		Memo:          memo,
		Operations: []xdr.Operation{
			{Body: body},
		},
	}

	var txBytes bytes.Buffer
	_, err = xdr.Marshal(&txBytes, tx)
	if err != nil {
		panic(err)
	}

	txHash := hash.Hash(txBytes.Bytes())
	signature, err := skp.Sign(txHash[:])
	if err != nil {
		panic(err)
	}

	ds := xdr.DecoratedSignature{
		Hint:      skp.Hint(),
		Signature: xdr.Signature(signature[:]),
	}

	txe := xdr.TransactionEnvelope{
		Tx:         tx,
		Signatures: []xdr.DecoratedSignature{ds},
	}

	var txeBytes bytes.Buffer
	_, err = xdr.Marshal(&txeBytes, txe)
	if err != nil {
		panic(err)
	}
	txeB64 := base64.StdEncoding.EncodeToString(txeBytes.Bytes())

	fmt.Printf("tx base64: %s", txeB64)
	// Output: tx base64: AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAkCqQJepwMIpx3O5TN76xsBApVqnpDD2f0X2SQ+EcMz4AAAAAAAAAAB3NZQAAAAAAAAAAAan/sYsAAABAIAWtiYQfI5rp4ZGE98rmyvfXvVX0+340nkEjYYDhFVFE7FCJERjlChY+pgR8THv7jgbtEgPZAgwJwXSrZh7mAw==
}

func ExampleUsingTransactionBuilder() {
	skp := keypair.MustParse("SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME")
	dkp := keypair.MustParse("GCICVEBF5JYDBCTR3TXFGN56WGYBAKKWVHUQYPM72F6ZEQ7BDQZT4NFZ")

	asset, err := xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
	if err != nil {
		panic(err)
	}

	var destination xdr.AccountId
	err = destination.SetAddress(dkp.Address())
	if err != nil {
		panic(err)
	}

	op := xdr.PaymentOp{
		Destination: destination,
		Asset:       asset,
		Amount:      50 * 10000000,
	}

	memo, err := xdr.NewMemo(xdr.MemoTypeMemoNone, nil)

	var source xdr.AccountId
	err = source.SetAddress(skp.Address())
	if err != nil {
		panic(err)
	}

	tx := xdr.Transaction{
		SourceAccount: source,
		Fee:           10,
		SeqNum:        xdr.SequenceNumber(1),
		Memo:          memo,
	}

	transactionBuilder := builder.GetInstance(&tx)
	transactionBuilder.MakeOperation(xdr.OperationTypePayment, op)
	err = transactionBuilder.Sign("SDKJ2BUKQ5TCMSLRQBAFSEVJ3LBXFGHEKKPTYNCDWSOJ4CFGFR5SKRME", "Test SDF Network ; September 2015")
	if err != nil {
		panic(err)
	}
	txeB64, err := transactionBuilder.ToBase64()
	if err != nil {
		panic(err)
	}
	fmt.Printf("tx base64: %s", txeB64)
	// Output: tx base64: AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAkCqQJepwMIpx3O5TN76xsBApVqnpDD2f0X2SQ+EcMz4AAAAAAAAAAB3NZQAAAAAAAAAAAan/sYsAAABA3LBOukhEwmXdSSi7zTUC5MfZM8opRWI0SZ8cscGO/av0aJXfZzRtGXiWE2cK7voF1GswJUiEpNo/EEWHYvzjCA==
}

func ExampleDecodeTransaction() {
	data := "AAAAABjCG5iSDJdtHOz38Hfkb0RYQP11Tu5cdDF+Teqp/7GLAAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAkCqQJepwMIpx3O5TN76xsBApVqnpDD2f0X2SQ+EcMz4AAAAAAAAAAB3NZQAAAAAAAAAAAan/sYsAAABA3LBOukhEwmXdSSi7zTUC5MfZM8opRWI0SZ8cscGO/av0aJXfZzRtGXiWE2cK7voF1GswJUiEpNo/EEWHYvzjCA=="
	rawr := strings.NewReader(data)
	b64r := base64.NewDecoder(base64.StdEncoding, rawr)

	var tx xdr.TransactionEnvelope
	bytesRead, err := xdr.Unmarshal(b64r, &tx)

	fmt.Println(tx)

	fmt.Printf("read %d bytes\n", bytesRead)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("This tx has %d operations\n", len(tx.Tx.Operations))
	// {{{PublicKeyTypePublicKeyTypeEd25519 0xc0000d01c0} 10 1 <nil> {MemoTypeMemoNone <nil> <nil> <nil> <nil>} [{<nil> {OperationTypePayment <nil> 0xc00008afc0 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>}}] {0}} [{[169 255 177 139] [220 176 78 186 72 68 194 101 221 73 40 187 205 53 2 228 199 217 51 202 41 69 98 52 73 159 28 177 193 142 253 171 244 104 149 223 103 52 109 25 120 150 19 103 10 238 250 5 212 107 48 37 72 132 164 218 63 16 69 135 98 252 227 8]}]}
	// Output: read 192 bytes
	// This tx has 1 operations
}
