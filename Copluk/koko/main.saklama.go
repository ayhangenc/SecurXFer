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
	"os"
	"strings"
)

var (
	proto = flag.String("proto", "", "Transmission Protocol between peers (TCP -Default- or UDP)")
	enc   = flag.Int("enc", 0, "encryption enable (TRUE:1 or FALSE:0)")
)

func main() {

	cipherKey256 := []byte("_08_bit__16_bit__24_bit__32_bit_")
	encryptedMessage := "XX"

	//let's parse the arguments first...
	flag.Parse()
	fmt.Println("Transmission Type: ", *proto)
	fmt.Println("Encryption: ", *enc)

	//let's get the message from keyboard

	messageReader := bufio.NewReader(os.Stdin)
	fmt.Print("Please type your message here: ")
	inputMessage, _ := messageReader.ReadString('\n')

	// remove CR from end
	inputMessage = strings.TrimSuffix(inputMessage, "\n")
	//
	fmt.Println("Your message is: " + inputMessage)

	// Now we have the message, first check if encryption is enabled, if so encrypt the message

	switch *enc {
	case 0:
		// no encryption
		xferMessage := inputMessage
		fmt.Println("Your message to be transferred is: " + xferMessage)

	case 1:
		//encrypt the message
		fmt.Println("Your message is going to be encrypted. ")
		fmt.Println("mesaj :", inputMessage)
		fmt.Printf("mesajın tipi: %T\n", inputMessage)
		fmt.Println("sifer :", cipherKey256)

		bobi := ""

		encryptedMessage, _, bobi = encryptMessage(cipherKey256, inputMessage)
		// _, _, bobi := encryptMessage(cipherKey256, inputMessage)

		fmt.Println("bobi: ", bobi)

		//Print the key and cipher text:
		fmt.Println("Cipher Key for AES 256:", string(cipherKey256))
		fmt.Printf("Encrypted Message: %T\n", encryptedMessage)

	}

	// now we have the encrypted message

	// lets decrypt

	decryptedMessage, _ := decryptMessage(cipherKey256, "bobi")

	//Print re-decrypted text:
	fmt.Println("Decrypted Message: ", decryptedMessage)

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
func encryptMessage(key []byte, message string) (encoded string, err error, bibi string) {
	//Create byte array from the input string
	plainText := []byte(message)
	fmt.Println("message in byets: ", plainText)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	fmt.Println("block: ", block)
	//IF NewCipher failed, exit:
	// if err != nil {
	//	   return
	// }

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	fmt.Println("ciphertext: ", cipherText)
	fmt.Println("block size: ", aes.BlockSize)
	fmt.Println("iv: ", iv, len(iv))

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	// fmt.Println("stream: ", stream)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	fmt.Println("stream & plaintext: ", stream, plainText)
	//Return string encoded in base64
	fmt.Println("cipherText: ", cipherText)
	fmt.Println("cipherText encoded: ", base64.RawStdEncoding.EncodeToString(cipherText))
	bibi = base64.RawStdEncoding.EncodeToString(cipherText)

	return base64.RawStdEncoding.EncodeToString(cipherText), err, bibi
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
func decryptMessage(key []byte, secure string) (decoded string, err error) {
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
