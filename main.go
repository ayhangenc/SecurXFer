package main

import (
	"flag"
	"fmt"
	"log"
)

/*
*	FILE			: main.go
*	PROJECT			: Secure Message Transfer - Client/Server
*	PROGRAMMER		: Ayhan GENC, ref: https://github.com/ayhangenc
*	FIRST VERSION	: 19 Sept. 2022
*	DESCRIPTION		:
*		The project is a learning exercise for go. There would be different functions, including byte operations,
*		CRC generation and checking, encryption/decryption and TCP/UDP comm. using the Golang standard libraries
*		for AES(CFB), CRC etc...
*		DISCLAIMER: This is only for my personal learning. So NO WARRANTIES....
*		Credits: Daniel Pieczewski, ref: https://github.com/mickelsonm for AES encryption/decryption clues... .
* 				 Kevin FOO , ref: https://oofnivek.medium.com for CRC-32 clues
 */

const (
	IEEE = 0xedb88320 //CRC-32
	// Castagnoli's polynomial, used in iSCSI.
	// Has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/26.231911
	// Castagnoli = 0x82f63b78
	// Koopman's polynomial.
	// Also has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/DSN.2002.1028931
	// Koopman = 0xeb31d82e
)

var (
	cipherKey256 = []byte("_08_bit__16_bit__24_bit__32_bit_") //32-bit key for AES-256
	//cipherKey192 = []byte("_08_bit__16_bit__24_bit_") //24-bit key for AES-192
	//cipherKey128 = []byte("_08_bit__16_bit_") //16-bit key for AES-128

	mode      = flag.String("mode", "server", "server to listen, client to send")
	proto     = flag.String("proto", "TCP", "Transmission Protocol between peers (tcp -Default- or udp)")
	enc       = flag.Int("enc", 0, "Encryption enable (TRUE:1 or FALSE:0)")
	ipAddress = flag.String("ipAddress", "", "IP Address (A.B.C.D format in decimal")
	port      = flag.String("port", "5000", "IP port number - Default 5000")
	msg       = flag.String("msg", "", "Message to send to the other party (in quotes), if none, then user input via keyboard would be initialized")
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	flag.Parse()

	switch *mode {
	case "server":
		fmt.Println("mode: server")

		// Listen & Decode by calling server.go  - parameters: *proto, *port
		server(*proto, *port)
	case "client":

		// // Initiate & Encode by calling client.go in loop  - parameters: *proto, *enc, *msg, *ipAddress, *port
		for {
			client(*proto, *enc, *ipAddress, *port, *msg)
		}
	}
}
