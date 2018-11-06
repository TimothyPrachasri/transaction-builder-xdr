package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	xdrHelper "github.com/kafakk/xdr-builder"
	h "github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

func main() {
	//CreateAccount
	destination, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(destination.Address())
	data, err := xdrHelper.CreateAccount(destination.Address(), -1)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	rawr := strings.NewReader(data)
	b64r := base64.NewDecoder(base64.StdEncoding, rawr)
	var op xdr.CreateAccountOp
	bytesRead, err := xdr.Unmarshal(b64r, &op)
	fmt.Printf("read %d bytes\n", bytesRead)
	body, err := xdr.NewOperationBody(xdr.OperationTypeCreateAccount, op)
	if err != nil {
		panic(err)
	}
	//SetOption
	// data, err := xdrHelper.SetOption("GAEBJVQJJO5ZPRJ2ZPNSDJLMNN64REZO7S5VUZAMNLI34B5XUQVD3URR", 0x1, 0x1, 3, 2,
	//  0, 1, "example.com", "", 0)
	// rawr := strings.NewReader(data)
	// b64r := base64.NewDecoder(base64.StdEncoding, rawr)
	// var op xdr.SetOptionsOp
	// bytesRead, err := xdr.Unmarshal(b64r, &op)
	// fmt.Println(bytesRead)
	// body, err := xdr.NewOperationBody(xdr.OperationTypeSetOptions, op)
	skp := keypair.MustParse("SBHUX66XAFVGQMBTOECROJR257JWRDTSGK7KOQZ7PUS3INLLOQ3TRIVG")
	var source xdr.AccountId
	err = source.SetAddress(skp.Address())
	if err != nil {
		panic(err)
	}
	tx := xdr.Transaction{
		SourceAccount: source,
		Fee:           100,
		SeqNum:        xdr.SequenceNumber(1008097543847999),
		// Memo:          memo,
		Operations: []xdr.Operation{
			// {SourceAccount: &source},
			{Body: body},
		},
	}
	var txBytes bytes.Buffer
	_, err = xdr.Marshal(&txBytes, tx)
	if err != nil {
		panic(err)
	}
	// txHash := hash.Hash(txBytes.Bytes())
	txHash, err := transactionHash(&tx, "Test SDF Network ; September 2015")
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
	resp, err := h.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		switch e := err.(type) {
		case *h.Error:
			fmt.Println("err type=" + e.Problem.Type)
			fmt.Println("err detailed=" + e.Problem.Detail)
			fmt.Println("err extras=" + string(e.Problem.Extras["result_codes"]))
		}
		log.Fatal(err)
	}
	log.Println("Transaction successfully: ")
	fmt.Println(resp)
}

// TransactionHash returns transaction hash for a given Transaction based on the network
func transactionHash(tx *xdr.Transaction, networkPassphrase string) ([32]byte, error) {
	var txBytes bytes.Buffer
	_, err := fmt.Fprintf(&txBytes, "%s", hash.Hash([]byte(networkPassphrase)))
	if err != nil {
		return [32]byte{}, err
	}
	_, err = xdr.Marshal(&txBytes, xdr.EnvelopeTypeEnvelopeTypeTx)
	if err != nil {
		return [32]byte{}, err
	}
	_, err = xdr.Marshal(&txBytes, tx)
	if err != nil {
		return [32]byte{}, err
	}
	return hash.Hash(txBytes.Bytes()), nil
}
