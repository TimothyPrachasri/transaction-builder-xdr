package main

import (
	"encoding/json"
	"fmt"

	"github.com/stellar/go/clients/horizon"
)

func main() {
	blob := "AAAAAFsAPNHwcy2ZPYftEEoI+dAPr0ZBN+vuXUKPEKcq2mmtAAAAZAAI8fMAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAdfeOkY7szblb3J8lDy2i0o1ssnDcDOkFjxjwFx/sV+gAAAAAAAAAADuaygAAAAAAAAAAASraaa0AAABAgoMYas+6bNnZHvW4tVFg+ZNYNOmzI+WYwqf/ZRJ22m3KOQ4fOHQbA8bY1IePnsGquwh1lfH3y8HJ0fBaK/jMCg=="
	acct, err := horizon.DefaultTestNetClient.LoadAccount("GD3F7GSWWTP6MQVFO6ZT64TDS7XR2KUPKCFCQVBG4TF6DCYLF5SWJT7Z")
	res, _ := json.Marshal(acct)
	fmt.Println(string(res))
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(blob)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
}
