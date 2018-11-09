package xdrsigner

import (
	"bytes"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/xdr"
)

// TransactionEnvelope represents a struct for collectioning crucials fields for encryping to envelope
type TransactionEnvelope struct {
	TransactionXDR *xdr.Transaction
	Signatures     []xdr.DecoratedSignature `xdrmaxsize:"20"`
}

// Sign method will be called after make or makeAll method for transaction signing protocol.
// It will receive signer and networkPassPhrase as a string.
func (transactionEnvelope *TransactionEnvelope) Sign(signer, networkPassPhrase string) (err error) {
	kp, err := keypair.Parse(signer)
	if err != nil {
		return errors.Wrap(err, "parse failed")
	}

	hash, err := network.HashTransaction(transactionEnvelope.TransactionXDR, networkPassPhrase)
	if err != nil {
		return errors.Wrap(err, "hash tx failed")
	}

	sig, err := kp.SignDecorated(hash[:])
	if err != nil {
		return errors.Wrap(err, "sign tx failed")
	}

	transactionEnvelope.Signatures = append(transactionEnvelope.Signatures, sig)
	return nil
}

// SignAll is similar to Sign method but it will receive an array of SignerForm which contains signer and networkPassPhrase.
func (transactionEnvelope *TransactionEnvelope) SignAll(signers []string, networkPassPhrase string) (err error) {
	for _, ech := range signers {
		err := transactionEnvelope.Sign(ech, networkPassPhrase)
		if err != nil {
			break
		}
	}
	return
}

// ToBytesEncoded encodes the builder's underlying envelope to XDR using xdr based Marshal
func (transactionEnvelope *TransactionEnvelope) ToBytesEncoded() ([]byte, error) {
	var txBytes bytes.Buffer
	_, err := xdr.Marshal(&txBytes, struct {
		TransactionXDR xdr.Transaction
		Signatures     []xdr.DecoratedSignature `xdrmaxsize:"20"`
	}{
		*transactionEnvelope.TransactionXDR,
		transactionEnvelope.Signatures,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshal xdr failed")
	}

	return txBytes.Bytes(), nil
}

// ToBase64 change format of XDR encoded into base64 string
func (transactionEnvelope *TransactionEnvelope) ToBase64() (string, error) {
	bs, err := transactionEnvelope.ToBytesEncoded()
	if err != nil {
		return "", errors.Wrap(err, "get raw bytes failed")
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

// GetEnveloperInstance is like a new object method in OOP
func GetEnveloperInstance(transactionXdrBase64 string) (result TransactionEnvelope) {
	var transaction xdr.Transaction
	xdr.SafeUnmarshalBase64(transactionXdrBase64, &transaction)
	result = TransactionEnvelope{TransactionXDR: &transaction}
	return
}
