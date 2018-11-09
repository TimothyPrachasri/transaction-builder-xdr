# transaction-builder-xdr
This is a library for parsing transaction and encoded payment xdr to encoded transaction xdr
# There are 2 parts

## Builder part

> builder path is to create encoded xdr transaction base 64 string and send to enveloper to envelope it with signature

### There are 2 ways to fill a transaction

- First one is creating xdr.Transaction struct and send it to parameter (_think like new object()_)
   ```go
    tx = xdr.Transaction{
		SourceAccount: Source,
		Fee:           xdr.Uint32(int(100) * 1),
		SeqNum:        xdr.SequenceNumber(1),
		Memo:          Memo,
	}
	transactionBuilder = xdrSigner.GetBuilderInstance(&tx)
	tB64, err = transactionBuilder.ToBase64()
   ```
- Second is creating empty builder instance and then make it through provided function
  ```go 
	transactionBuilder = xdrSigner.GetBuilderInstance()
	transactionBuilder.MakeMemo(xdr.MemoTypeMemoNone, nil)
	transactionBuilder.MakeSequenceNumber(uint64(1))
	transactionBuilder.MakeFee(100)
	transactionBuilder.MakeSourceAccount(SourceSeed)
	transactionBuilder.MakeOperation(opB64)
	tB64, err = transactionBuilder.ToBase64()
  ```
	- Here is a list of provided functions
		- `MakeMemo(memoType xdr.MemoType, value interface{}) (err error)`
		- `MakeOperation(operationB64 string) (err error)`
		- `MakeAllOperations(forms []string) (err error)`
		- `MakeSourceAccount(addressOrSeed string) (err error)`
		- `MakeFee(fee xdr.Uint32) (err error)`
		- `MakeBaseFee(baseFee uint64) (err error)`
		- `MakeSequenceNumber(sequenceNumber uint64) (err error)`
		- `ToBytesEncoded() ([]byte, error)`
		- `ToBase64() (string, error) (err error)`
___

## Enveloper part

> Enveloper will receive xdr transaction as a string and provide function sign in order to let user sign, envelope all of this stuff, and turn transaction with signature(s) into transaction envelope base 64 string

```go
	transactionEnvelope = xdrSigner.GetEnveloperInstance(tB64)
	err = transactionEnvelope.Sign(SourceSeed, PassPhrase)
	txeB64, err = transactionEnvelope.ToBase64()
	_, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
```

- Here is a list of provided functions
	- `Sign(signer, networkPassPhrase string) (err error)`
	- `SignAll(signers []string, networkPassPhrase string) (err error)`
	- `ToBytesEncoded() ([]byte, error)`
	- `ToBase64() (string, error)`

___

