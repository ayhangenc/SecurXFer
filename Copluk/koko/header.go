package main
// header contents: 2 byte 0 F 0 0 + 4 byte lenght in hex + CRC32 (IEEE) 4 byte

import (
	"fmt"
	"encoding/binary"
	// "bytes"
	"encoding/hex"
	// "strconv"
)


func main() {

	headerSignature := []byte{0x0F, 0x0}
	fmt.Printf("headerSig :%x\n", headerSignature)
	fmt.Println(headerSignature)
	crcString :="3610a686"
	headerCRC,_ := hex.DecodeString(crcString)
	fmt.Printf("headerCRC :%x\n", headerCRC)
	fmt.Println(headerCRC)
	messageIn :="HM++" // WiNduHwjWPqeuwadBti+s1Jc1Q"
	message2Send := hex.EncodeToString([]byte(messageIn))

	fmt.Println("message2send", message2Send)

	messageLen := make([]byte, 2)
	fmt.Println("message2send Len", len(message2Send))

	binary.BigEndian.PutUint16(messageLen, uint16(len(message2Send)))

	// decode
	i := int16(binary.BigEndian.Uint16(messageLen))
	fmt.Println("i", i)
	// messageLenght := []byte(intconv.Itoa(len(message2Send)))
	fmt.Println("messageLenght", messageLen)


	headerI := append( headerSignature[:], headerCRC[:]...)
	header:= append( headerI[:], messageLen[:]...)
	fmt.Println("header", header, len(header))

	fullMessage2Send := append(header[:],message2Send[:]...)
	fmt.Println("full message to send", fullMessage2Send)
/*
	h1 := []byte{0xEF, 0xE0}
	h2 := []byte("def")
	h3 := []byte("1234")
	h4 := []byte("5678")

	tg := make([]byte, 3+3+4+(64)) // allocate enough space for the 3 IDs and a max 1MB of extra data

	hex.Decode(tg[:3], h1)

	// hh1,_ := hex.DecodeString(h1)

	// fmt.Printf("h1 : %x\n", hh1)

	// hh2, _ := hex.DecodeString(h2)

	// fmt.Printf("h2 : %x\n", hh2)

	hex.Decode(tg[3:6], h2)
	hex.Decode(tg[6:10], h3)
	l, _ := hex.Decode(tg[10:], h4)
	fmt.Printf("l : %v\n", l)
	fmt.Println("l : ", l)

	// tg = tg[:64]
	fmt.Printf("h1 : %x\n", h1)
	fmt.Printf("h2 : %s\n", h2)
	fmt.Printf("h3 : %s\n", h3)
	fmt.Printf("h4 : %s\n", h4)
	fmt.Printf("tg : %s\n", tg)
*/


}