package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

/*
 *	FILE			: client.go
 *	PROJECT			: Secure Message Transfer - Client/Server
 *	PROGRAMMER		: Ayhan GENC, ref: https://github.com/ayhangenc
 *	FIRST VERSION	: 19 Sept. 2022
 *	DESCRIPTION		:
 *		Client code for SecureXFER
 *		The project is a learning exercise for go. There would be different functions, including byte operations,
 *		CRC generation and checking, encryption/decryption and TCP/UDP comm. using the Golang standard libraries
 *		for AES(CFB), CRC etc...
 *		DISCLAIMER: This is only for my personal learning. So NO WARRANTIES....
 *		Credits: Daniel Pieczewski, ref: https://github.com/mickelsonm for AES encryption/decryption clues... .
 * 				 Kevin FOO , ref: https://oofnivek.medium.com for CRC-32 clues
 */

func client(proto string, enc int, ipAddress string, port string, msg string) {

	messageFromCLI := msg
	inputMessage := ""
	if messageFromCLI == "" {
		messageReader := bufio.NewReader(os.Stdin)
		fmt.Print("Please type your message here: ")
		inputMessage, _ = messageReader.ReadString('\n')
		inputMessage = strings.TrimSuffix(inputMessage, "\n") // remove CR from end
	} else {
		inputMessage = messageFromCLI
	}

	var messagetoCRC []byte
	switch enc {
	case 0: // no encryption
		messagetoCRC = []byte(inputMessage)

	case 1: //encrypt the message
		encrypted, err := encryptMessage(cipherKey256, inputMessage)
		if err != nil { //IF the encryption failed:
			log.Println(err) //Print error message:
			os.Exit(-3)      // -3: Encryption error
		}
		messagetoCRC = []byte(encrypted)
	}
	messageCRC := crcGenerate(messagetoCRC)                         // CRC Generation
	messageHeader := headerGenerate(messagetoCRC, messageCRC, &enc) // Header Generation
	fullMessage2Send := append(messageHeader[:], messagetoCRC[:]...)
	addresstoSend := ipAddress + ":" + port
	tcpAddr, err := net.ResolveTCPAddr(proto, addresstoSend)
	checkError(err)
	conn, err := net.DialTCP(proto, nil, tcpAddr)
	checkError(err)
	_, err = conn.Write(fullMessage2Send)
	checkError(err)
	fmt.Println("Message sent!...")
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

func encryptMessage(cipherKey []byte, messagetoEncrypt string) (encodedMessage string, err error) {

	messageText := []byte(messagetoEncrypt)          //Create byte array from the input string
	encryptionBlock, err := aes.NewCipher(cipherKey) //Create a new AES cipher using the key
	if err != nil {                                  //if failed, exit:
		return
	}
	cipherText := make([]byte, aes.BlockSize+len(messageText))           //Make the cipher text a byte array of size BlockSize + the length of the message
	intermediateText := cipherText[:aes.BlockSize]                       //intermediateText is the ciphertext up to the blocksize (16)
	if _, err = io.ReadFull(rand.Reader, intermediateText); err != nil { //if failed, exit:
		return
	}
	encryptedStream := cipher.NewCFBEncrypter(encryptionBlock, intermediateText) //Encrypt the message
	encryptedStream.XORKeyStream(cipherText[aes.BlockSize:], messageText)

	return base64.RawStdEncoding.EncodeToString(cipherText), err //Return string encoded in base64
}

/*
 *	FUNCTION		: CRC Generation
 *	DESCRIPTION		:
 *		This function takes a string generate CRC-32
 *
 *	PARAMETERS		:
 *		TBD
 *		XX byte[] key	: Byte array containing the cipher key
 *		XX string secure	: String containing an encrypted message
 *
 *	RETURNS			:
 *		TBD
 *		XX string decoded	: String containing the decrypted equivalent of secure
 *		XX error err	: Error message
 */

func crcGenerate(message2CRC []byte) (crcFromMessage []byte) {
	crc32Table := crc32.MakeTable(IEEE)
	crcIntermediate := crc32.Checksum(message2CRC, crc32Table)
	crcFromMessage = make([]byte, 4)
	binary.BigEndian.PutUint16(crcFromMessage, uint16(crcIntermediate))
	return crcFromMessage
}

/*
 *	FUNCTION		: Header Generation
 *	DESCRIPTION		:
 *		This function takes the message and CRC values to generate the header for the message to wire..
 *
 *	PARAMETERS		:
 *		TBD
 *		XX byte[] key	: Byte array containing the cipher key
 *		XX string secure	: String containing an encrypted message
 *
 *	RETURNS			:
 *		TBD
 *		XX string decoded	: String containing the decrypted equivalent of secure
 *		XX error err	: Error message
 */

func headerGenerate(message2Send []byte, crcfromMessage []byte, enc *int) (headertoMessage []byte) {

	var headerSignature []byte
	switch *enc {
	case 0:
		headerSignature = []byte{0xF0, 0x00}
	case 1:
		headerSignature = []byte{0xF0, 0x01}
	}
	crcString := fmt.Sprintf("%x", crcfromMessage)
	headerCRC, _ := hex.DecodeString(crcString)
	messageLen := make([]byte, 2)
	binary.BigEndian.PutUint16(messageLen, uint16(len(message2Send)))
	headerI := append(headerSignature[:], headerCRC[:]...)
	header := append(headerI[:], messageLen[:]...)

	return header
}
