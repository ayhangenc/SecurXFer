package main

import (
	"encoding/base64"
	"fmt"
)

func main (){

	pot := "OeRl4/dLt/finndLCdKtgz4"

	nano := deco(pot)

	fmt.Println("nano: ",nano)
}

func deco(sec string) (chipherText string) {
	fmt.Println("secure: ", sec)

	cipherText, err := base64.RawStdEncoding.DecodeString(sec) //"OeRl4/dLt/finndLCdKtgz4") // 6z6bbY0ZihAUZYUv7h8sP6o"
	fmt.Println("cipherText: ", cipherText)
	//IF DecodeString failed, exit:
	if err != nil {
		return
	}
return
}
