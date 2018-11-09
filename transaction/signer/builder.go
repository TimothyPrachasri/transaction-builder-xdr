package xdrsigner

import (
	"bytes"
	"encoding/base64"

	"transaction-builder-xdr/utils"

	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

var defaultBaseFee uint64 = 100

// TransactionBuilder is a decorator struct contain xdr.Transaction, and baseFee
type TransactionBuilder struct {
	TransactionXDR *xdr.Transaction
	BaseFee        uint64
}

// MakeMemo create memo for transaction.
func (transactionBuilder *TransactionBuilder) MakeMemo(memoType xdr.MemoType, value interface{}) (err error) {
	transactionBuilder.TransactionXDR.Memo, err = xdr.NewMemo(memoType, value)
	return
}

// MakeDefaultFee create a default fee for transaction in case fee is not defined yet
func (transactionBuilder *TransactionBuilder) MakeDefaultFee() error {
	if transactionBuilder.BaseFee == 0 {
		transactionBuilder.BaseFee = defaultBaseFee
	}
	if transactionBuilder.TransactionXDR.Fee == 0 {
		transactionBuilder.TransactionXDR.Fee = xdr.Uint32(int(transactionBuilder.BaseFee) * len(transactionBuilder.TransactionXDR.Operations))
	}
	return nil
}

// MakeOperation method is a method in which TransactionBuilder make a TransactionXDR based on operationType and paymentOperationForm.
func (transactionBuilder *TransactionBuilder) MakeOperation(operationB64 string) (err error) {
	var operation xdr.Operation
	xdr.SafeUnmarshalBase64(operationB64, &operation)
	transactionBuilder.TransactionXDR.Operations = append(transactionBuilder.TransactionXDR.Operations, operation)
	return
}

// MakeAllOperations method is like a Make method but makeAll receives an array of OperationForm which contain operationType and paymentOperationForm.
func (transactionBuilder *TransactionBuilder) MakeAllOperations(forms []string) (err error) {

	for _, ech := range forms {
		err = transactionBuilder.MakeOperation(ech)
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

// ToBytesEncoded encodes the builder's underlying envelope to XDR using xdr based Marshal
func (transactionBuilder *TransactionBuilder) ToBytesEncoded() ([]byte, error) {
	transactionBuilder.MakeDefaultFee()
	var txBytes bytes.Buffer
	_, err := xdr.Marshal(&txBytes, transactionBuilder.TransactionXDR)
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

// GetBuilderInstance is like a new object method in OOP
func GetBuilderInstance(transaction ...*xdr.Transaction) (result TransactionBuilder) {
	if len(transaction) == 1 {
		tx := *transaction[0]
		result = TransactionBuilder{TransactionXDR: &tx}
		return
	}
	result = TransactionBuilder{TransactionXDR: &xdr.Transaction{}}
	return
}
