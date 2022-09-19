package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
 */

var (
	proto = flag.String("proto", "", "Transmission Protocol between peers (TCP -Default- or UDP)")
	enc   = flag.Int("enc", 0, "Encryption enable (TRUE:1 or FALSE:0)")
	msg   = flag.String("msg", "", "Message to send to the other party (in quotes), if none, then user input via keyboard would be initialized")
)

func main() {
	cipherKey256 := []byte("_08_bit__16_bit__24_bit__32_bit_") //32 bit key for AES-256
	//cipherKey192:= []byte("_08_bit__16_bit__24_bit_") //24 bit key for AES-192
	//cipherKey128:= []byte("_08_bit__16_bit_") //16 bit key for AES-128

	var inputMessage string

	flag.Parse()
	fmt.Println("Transmission Type: ", *proto)
	fmt.Println("Encryption: ", *enc)
	fmt.Println("Message to transfer: ", *msg)

	messageFromCLI := *msg
	if messageFromCLI == "" {
		messageReader := bufio.NewReader(os.Stdin)
		fmt.Print("Please type your message here: ")
		inputMessage, _ = messageReader.ReadString('\n')
		// remove CR from end
		inputMessage = strings.TrimSuffix(inputMessage, "\n")
	} else {
		inputMessage = messageFromCLI
	}

	fmt.Println("mesaj / boyu: ", inputMessage, len(inputMessage))

	//Encrypt the message:
	encrypted, err := encrypt(cipherKey256, inputMessage)

	//IF the encryption failed:
	if err != nil {
		//Print error message:
		log.Println(err)
		os.Exit(-2)
	}

	//Print the key and cipher text:
	fmt.Printf("\n\tCIPHER KEY: %s\n", string(cipherKey256))
	fmt.Printf("\tENCRYPTED: %s\n", encrypted)

	//Decrypt the text:
	decrypted, err := decrypt(cipherKey256, encrypted)

	//IF the decryption failed:
	if err != nil {
		log.Println(err)
		os.Exit(-3)
	}

	//Print re-decrypted text:
	fmt.Printf("\tDECRYPTED: %s\n\n", decrypted)
}

/*
 *	FUNCTION		: encrypt
 *	DESCRIPTION		:
 *		This function takes a string and a cipher key and uses AES to encrypt the message
 *
 *	PARAMETERS		:
 *		byte[] key	: Byte array containing the cipher key
 *		string message	: String containing the message to encrypt
 *
 *	RETURNS			:
 *		string encoded	: String containing the encoded user input
 *		error err	: Error message
 */

func encrypt(key []byte, message string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)
	fmt.Println("message in byets: ", plainText)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	fmt.Println("block: ", block)
	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	fmt.Println("ciphertext: ", cipherText)

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	fmt.Println("block size: ", aes.BlockSize)
	fmt.Println("iv: ", iv, len(iv))

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	fmt.Println("stream: ", stream)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
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

func decrypt(key []byte, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}
