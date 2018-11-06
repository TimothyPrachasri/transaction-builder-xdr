package transactionbuilder

import (
	"bytes"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/xdr"
)

// TransactionBuilder is a struct contain xdrTransaction
type TransactionBuilder struct {
	TransactionXDR *xdr.Transaction
	Signatures     []xdr.DecoratedSignature `xdrmaxsize:"20"`
}

// OperationForm is a form for building operationForm and send it to makeAll method.
type OperationForm struct {
	OperationType        xdr.OperationType
	PaymentOperationForm interface{}
}

// SignerForm is a form for transaction signing protocol.
type SignerForm struct {
	Signer            string
	NetworkPassPhrase string
}

// Make method is a method in which TransactionBuilder make a TransactionXDR based on operationType and paymentOperationForm.
func (transactionBuilder *TransactionBuilder) Make(operationType xdr.OperationType, paymentOperationForm interface{}) (err error) {
	operation := xdr.Operation{SourceAccount: &xdr.AccountId{}}
	operation.Body, err = xdr.NewOperationBody(operationType, paymentOperationForm)
	transactionBuilder.TransactionXDR.Operations = append(transactionBuilder.TransactionXDR.Operations, operation)
	return
}

// MakeAll method is like a Make method but makeAll receives an array of OperationForm which contain operationType and paymentOperationForm.
func (transactionBuilder *TransactionBuilder) MakeAll(forms []OperationForm) (err error) {
	for _, ech := range forms {
		err = transactionBuilder.Make(ech.OperationType, ech.PaymentOperationForm)
		if err != nil {
			break
		}
	}
	return
}

// Sign method will be called after make or makeAll method for transaction signing protocol.
// It will receive signer and networkPassPhrase as a string.
func (transactionBuilder *TransactionBuilder) Sign(signer, networkPassPhrase string) (err error) {
	hash, err := network.HashTransaction(transactionBuilder.TransactionXDR, networkPassPhrase)
	if err != nil {
		return errors.Wrap(err, "hash tx failed")
	}

	kp, err := keypair.Parse(signer)
	if err != nil {
		return errors.Wrap(err, "parse failed")
	}

	sig, err := kp.SignDecorated(hash[:])
	if err != nil {
		return errors.Wrap(err, "sign tx failed")
	}

	transactionBuilder.Signatures = append(transactionBuilder.Signatures, sig)
	return nil
}

// ToBytesEncoded encodes the builder's underlying envelope to XDR using xdr based Marshal
func (transactionBuilder *TransactionBuilder) ToBytesEncoded() ([]byte, error) {
	var txBytes bytes.Buffer
	_, err := xdr.Marshal(&txBytes, struct {
		TransactionXDR *xdr.Transaction
		Signatures     []xdr.DecoratedSignature `xdrmaxsize:"20"`
	}{
		transactionBuilder.TransactionXDR,
		transactionBuilder.Signatures,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshal xdr failed")
	}

	return txBytes.Bytes(), nil
}

// ToBase64 change format of XDR encoded into base64 string
func (transactionBuilder *TransactionBuilder) ToBase64() (string, error) {
	bs, err := transactionBuilder.ToBytesEncoded()
	if err != nil {
		return "", errors.Wrap(err, "get raw bytes failed")
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

// SignAll is similar to Sign method but it will receive an array of SignerForm which contains signer and networkPassPhrase.
func (transactionBuilder *TransactionBuilder) SignAll(forms []SignerForm) (err error) {
	for _, ech := range forms {
		err := transactionBuilder.Sign(ech.Signer, ech.NetworkPassPhrase)
		if err != nil {
			break
		}
	}
	return
}

// GetInstance is like a new object method in OOP
func GetInstance() TransactionBuilder {
	result := TransactionBuilder{TransactionXDR: &xdr.Transaction{}}
	return result
}
