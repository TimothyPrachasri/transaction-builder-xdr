package main

import (
	"fmt"

	"github.com/stellar/go/clients/horizon"
)

func main() {
	blob := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAAZAAI8fMAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAdfeOkY7szblb3J8lDy2i0o1ssnDcDOkFjxjwFx/sV+gAAAAAAAAAAAAAAGQAAAAAAAAAASraaa0AAABAQVkwp5giAPb+it3dZlGMDzC7YmbaC3ESwUF1ud440EYHSTCsa2lW547GSTZMFfRgb2gOAmuKooD7lH4qYC0dAg=="
	if _, err := horizon.DefaultTestNetClient.LoadAccount("GD3F7GSWWTP6MQVFO6ZT64TDS7XR2KUPKCFCQVBG4TF6DCYLF5SWJT7Z"); err != nil {
		panic(err)
	}
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(blob)
	if err != nil {
		panic(err)
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
}
