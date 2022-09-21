package main

// Kevin FOO , https://oofnivek.medium.com, 2020

import (
"os"
"fmt"
"hash/crc32"
)
const (
	// IEEE is by far and away the most common CRC-32 polynomial.
	// Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
	IEEE = 0xedb88320
	// Castagnoli's polynomial, used in iSCSI.
	// Has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/26.231911
	Castagnoli = 0x82f63b78
	// Koopman's polynomial.
	// Also has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/DSN.2002.1028931
	Koopman = 0xeb31d82e
)
func main() {
	crc32q := crc32.MakeTable(IEEE)
	fmt.Println("crc32q", crc32q)
	str := os.Args[1]
	fmt.Println("str", str)
	fmt.Printf("%x\n", crc32.Checksum([]byte(str), crc32q))
	fmt.Println("crc", crc32.Checksum([]byte(str), crc32q))
}
