package utils

import (
	"errors"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// SetAccountId receive address of AccountId, point it and then set address itself.
func SetAccountId(addressOrSeed string, aid *xdr.AccountId) error {
	kp, err := keypair.Parse(addressOrSeed)
	if err != nil {
		return err
	}

	if aid == nil {
		return errors.New("aid is nil in setAccountId")
	}

	return aid.SetAddress(kp.Address())
}
