package transactionbuilder

import (
	"bytes"
	"encoding/base64"
	"transaction-builder-xdr/utils"

	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/xdr"
)

var DefaultBaseFee uint64 = 100

// TransactionBuilder is a decorator struct contain xdr.Transaction, signatures, and baseFee
type TransactionBuilder struct {
	TransactionXDR *xdr.Transaction
	Signatures     []xdr.DecoratedSignature `xdrmaxsize:"20"`
	BaseFee        uint64
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

// MakeDefaultFee create a default fee for transaction in case fee is not defined yet
func (transactionBuilder *TransactionBuilder) MakeDefaultFee() error {
	if transactionBuilder.BaseFee == 0 {
		transactionBuilder.BaseFee = DefaultBaseFee
	}
	if transactionBuilder.TransactionXDR.Fee == 0 {
		transactionBuilder.TransactionXDR.Fee = xdr.Uint32(int(transactionBuilder.BaseFee) * len(transactionBuilder.TransactionXDR.Operations))
	}
	return nil
}

// MakeOperation method is a method in which TransactionBuilder make a TransactionXDR based on operationType and paymentOperationForm.
func (transactionBuilder *TransactionBuilder) MakeOperation(operationType xdr.OperationType, paymentOperationForm interface{}) (err error) {
	operation := xdr.Operation{SourceAccount: &xdr.AccountId{}}
	operation.Body, err = xdr.NewOperationBody(operationType, paymentOperationForm)
	transactionBuilder.TransactionXDR.Operations = append(transactionBuilder.TransactionXDR.Operations, operation)
	return
}

// MakeAllOperations method is like a Make method but makeAll receives an array of OperationForm which contain operationType and paymentOperationForm.
func (transactionBuilder *TransactionBuilder) MakeAllOperations(forms []OperationForm) (err error) {

	for _, ech := range forms {
		err = transactionBuilder.MakeOperation(ech.OperationType, ech.PaymentOperationForm)
		if err != nil {
			break
		}
	}
	return
}

// MakeSourceAccount create source account for transaction
func (transactionBuilder *TransactionBuilder) MakeSourceAccount(addressOrSeed string) error {
	return utils.SetAccountId(addressOrSeed, &transactionBuilder.TransactionXDR.SourceAccount)
}

// MakeFee create fee for transaction
func (transactionBuilder *TransactionBuilder) MakeFee(fee xdr.Uint32) error {
	transactionBuilder.TransactionXDR.Fee = fee
	return nil
}

// MakeBaseFee create basefee for transaction
func (transactionBuilder *TransactionBuilder) MakeBaseFee(baseFee uint64) error {
	transactionBuilder.BaseFee = baseFee
	return nil
}

// MakeSequenceNumber create sequence number for transaction
func (transactionBuilder *TransactionBuilder) MakeSequenceNumber(sequenceNumber uint64) (err error) {
	transactionBuilder.TransactionXDR.SeqNum = xdr.SequenceNumber(sequenceNumber)
	return nil
}

// Sign method will be called after make or makeAll method for transaction signing protocol.
// It will receive signer and networkPassPhrase as a string.
func (transactionBuilder *TransactionBuilder) Sign(signer, networkPassPhrase string) (err error) {
	transactionBuilder.MakeDefaultFee()

	kp, err := keypair.Parse(signer)
	if err != nil {
		return errors.Wrap(err, "parse failed")
	}

	hash, err := network.HashTransaction(transactionBuilder.TransactionXDR, networkPassPhrase)
	if err != nil {
		return errors.Wrap(err, "hash tx failed")
	}

	sig, err := kp.SignDecorated(hash[:])
	if err != nil {
		return errors.Wrap(err, "sign tx failed")
	}

	transactionBuilder.Signatures = append(transactionBuilder.Signatures, sig)
	return nil
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

// ToBytesEncoded encodes the builder's underlying envelope to XDR using xdr based Marshal
func (transactionBuilder *TransactionBuilder) ToBytesEncoded() ([]byte, error) {
	var txBytes bytes.Buffer
	_, err := xdr.Marshal(&txBytes, struct {
		TransactionXDR xdr.Transaction
		Signatures     []xdr.DecoratedSignature `xdrmaxsize:"20"`
	}{
		*transactionBuilder.TransactionXDR,
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

// GetInstance is like a new object method in OOP
func GetInstance(transaction ...*xdr.Transaction) (result TransactionBuilder) {
	if len(transaction) == 1 {
		result = TransactionBuilder{TransactionXDR: transaction[0]}
		return
	}
	result = TransactionBuilder{TransactionXDR: &xdr.Transaction{}}
	return
}
