package transactionbuilder

import (
	"github.com/stellar/go/xdr"
)

// TransactionBuilder is a struct contain xdrTransaction
type TransactionBuilder struct {
	TransactionXDR *xdr.Transaction
}

// OperationForm is a form for building operationForm and send it to makeAll method
type OperationForm struct {
	operationType        xdr.OperationType
	paymentOperationForm interface{}
}

func (transactionBuilder *TransactionBuilder) make(operationType xdr.OperationType, paymentOperationForm interface{}) (err error) {
	operation := xdr.Operation{SourceAccount: &xdr.AccountId{}}
	operation.Body, err = xdr.NewOperationBody(operationType, paymentOperationForm)
	transactionBuilder.TransactionXDR.Operations = append(transactionBuilder.TransactionXDR.Operations, operation)
	return
}

func (transactionBuilder *TransactionBuilder) makeAll(forms []OperationForm) (err error) {
	for _, ech := range forms {
		err = transactionBuilder.make(ech.operationType, ech.paymentOperationForm)
		if err != nil {
			break
		}
	}
	return
}

func getInstance() TransactionBuilder {
	result := TransactionBuilder{&xdr.Transaction{}}
	return result
}
