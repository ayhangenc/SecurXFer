# SecurXFer (copyleft-2022)

secure message trasfer client/server mode // very beginning of a personal learning effort...

baby steps:


1. Parameter parsing -from CLI  -> OK

2. Message load (from file/command line parameter/typing etc) - from CLI and KB typing OK, load from file is possible but not implemented -> OK

3. Encryption - optional -aes 256 trials... -> OK

4. Header (Lenght of msg, encryption status, checksum) -> OK

5. Transmit (TCP/UDP) TCP for start, would add UDP in future, should work on fragmented transfer for large chunks... -> OK

6. Receive -> OK

7. Header processing -> OK

8. Decryption -> OK

9. Delivery (to screen/file etc) to screen for the moment, file is possible but not implemented -> OK

Additional notes:

* Consider sending large chunks in fregments... good homework :)
* Add an HTML interface as UI for the next release
* Public/Private Key approach for encryption, homework for encryption :)
(Sept. 21, 2022 first release)

