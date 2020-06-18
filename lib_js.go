package main

import (
	"fmt"
	"syscall/js"
)

func access(satAddr string, apiKey string, passphrase string) string {
	return fmt.Sprintf("%s - %s - %s", satAddr, apiKey, passphrase)
}

func accessThis(this js.Value, satAddr string, apiKey string, passphrase string) string {
	return fmt.Sprintf("%s: %s - %s - %s", this.Get("name").String(), satAddr, apiKey, passphrase)
}

func different(s string, b bool, f float32, i int8, u uint64) float32 {
	println(s, b, f, i, u)
	return 15.524
}
