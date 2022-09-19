package main

import (
	"bufio"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	proto = flag.String("proto", "", "Transmission Protocol between peers (TCP -Default- or UDP)")
	enc   = flag.Int("enc", 0, "encryption enable (TRUE:1 or FALSE:0)")
)

func EncryptAES(key []byte, plaintext []byte) string {

	c, err := aes.NewCipher(key)
	CheckError(err)

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	CheckError(err)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	fmt.Println("DECRYPTED:", s)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	//let's parse the arguments first...
	flag.Parse()
	fmt.Println("Transmission Type: ", *proto)
	fmt.Println("Encryption: ", *enc)

	//let's get the message from keyboard

	messageReader := bufio.NewReader(os.Stdin)
	fmt.Print("Please type your message here: ")
	inputMessage, _ := messageReader.ReadString('\n')
	// fmt.Println("Your message is: " + inputMessage)

	// Now we have the message, first check if encryption is enabled, if so encrypt the message

	switch *enc {
	case 0:
		// no encription
		xferMessage := inputMessage
		fmt.Println("Your message to be transferred is: " + xferMessage)
	case 1:
		//encrypt the message
		fmt.Println("Your message is going to be encrypted. ")

		fmt.Printf("mesajın tipi: %T", inputMessage)

		// cipher key
		key := "thisis32bitlongpassphraseimusing"

		// plaintext
		// pt := inputMessage // "This is a secret"

		//Make the cipher text a byte array of size BlockSize + the length of the message
		pt := make([]byte, aes.BlockSize+len(inputMessage))

		iv := pt[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return
		}

		c := EncryptAES([]byte(key), pt)

		fmt.Println("Encypted message: ", c)

		// lets decrypt

	}

}
