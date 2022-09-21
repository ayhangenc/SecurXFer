package main

/*
*	FILE			: server.go
*	PROJECT			: Secure Message Transfer - Client/Server
*	PROGRAMMER		: Ayhan GENC, ref: https://github.com/ayhangenc
*	FIRST VERSION	: 19 Sept. 2022
*	DESCRIPTION		:
*		Server code for SecureXFER
*		The project is a learning exercise for go. There would be different functions, including byte operations,
*		CRC generation and checking, encryption/decryption and TCP/UDP comm. using the Golang standard libraries
*		for AES(CFB), CRC etc...
*		DISCLAIMER: This is only for my personal learning. So NO WARRANTIES....
*		Credits: Daniel Pieczewski, ref: https://github.com/mickelsonm for AES encryption/decryption clues... .
* 				 Kevin FOO , ref: https://oofnivek.medium.com for CRC-32 clues
 */

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	input := make([]byte, 1024)
	nRead, err := conn.Read(input[0:])
	checkError(err)
	fullMessageReceived := input[:nRead]

	headSign := fullMessageReceived[:2] // check if message is authentic (header first 2-digits are FO:O0/01)
	if headSign[0] == 0xf0 {
		if headSign[1] == 0x00 || headSign[1] == 0x01 {
			fmt.Println("Message Is Authentic.")
		} else {
			fmt.Println("Message Is NOT Authentic!..")
			return
		}
	} else {
		fmt.Println("Message Is NOT Authentic!..")
		return
	}

	lenFromHeader := fullMessageReceived[6:8] // message lenght check (header digits 6 & 7 are lenght of message in hex)
	lenFromMessage := len(fullMessageReceived) - 8
	lenXCheck := int(lenFromHeader[0])*256 + int(lenFromHeader[1]) // (hex to int)
	if lenXCheck != lenFromMessage {
		fmt.Println("Message size is DIFFERENT from header!...")
		os.Exit(5) //message altered or corrupt during transmission
	}

	crcFromHeader := fullMessageReceived[2:6] // message crc check (header digits 2,3,4,5 are CRC digits in hex)
	crcFromMessage := crcGenerate(fullMessageReceived[8:])
	crcXCheck := bytes.Compare(crcFromHeader, crcFromMessage)
	if crcXCheck != 0 {
		fmt.Println("Message CRC is DIFFERENT from header!...")
		os.Exit(6) //message altered or corrupt during transmission
	}

	var messageRXBody []byte
	switch headSign[1] { // check if message is encrypted
	case 0x01: // encrypted
		fmt.Println("Message Is Encypted!...")
		messageRXBodySTR, err := decryptMessage(cipherKey256, string(fullMessageReceived[8:]))

		messageRXBody = []byte(messageRXBodySTR)

		if err != nil { //if message decrypt fails...
			log.Println(err)
			os.Exit(-3)
		}
	case 0x00: // not encrypted
		fmt.Println("Message Is Not Encrypted!..")
		messageRXBody = fullMessageReceived[8:]
	}
	fmt.Printf("Message RECEIVED from other party (be it encrypted or not) : %s\n", messageRXBody)
}

func server(proto string, port string) { // (proto string, port string)

	tcpAddr, err := net.ResolveTCPAddr(proto, fmt.Sprintf("127.0.0.1:%s", port))
	checkError(err)
	fmt.Println("IP add: ", tcpAddr)
	listener, err := net.ListenTCP(proto, tcpAddr)
	checkError(err)

	for {
		connectTo, err := listener.Accept()
		checkError(err)

		go handleConnection(connectTo)
	}
}

/*
 *	FUNCTION		: decrypt
 *	DESCRIPTION		:
 *		This function takes a string and a key and uses AES to decrypt the string into plain text
 *
 *	PARAMETERS		:
 *		byte[] key	: Byte array containing the cipher key
 *		string secure	: String containing an encrypted message
 *
 *	RETURNS			:
 *		string decoded	: String containing the decrypted equivalent of secure
 *		error err	: Error message
 */

func decryptMessage(cipherKey []byte, secureMessage string) (decodedMessage string, err error) {

	cipherText, err := base64.RawStdEncoding.DecodeString(secureMessage) // decode base64
	if err != nil {                                                      //IF DecodeString failed, exit:
		return
	}
	decryptionBlock, err := aes.NewCipher(cipherKey) //Create a new AES cipher with the key and encrypted message
	if err != nil {                                  //IF NewCipher failed, exit:
		return
	}
	if len(cipherText) < aes.BlockSize { //IF the length of the cipherText is less than 16 Bytes:
		log.Println("Ciphertext block size is too short!", err)
		os.Exit(-3)
	}
	intermediateText := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(decryptionBlock, intermediateText) //Decrypt the message
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}
