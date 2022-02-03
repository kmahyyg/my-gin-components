package main

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var Operation = flag.String("a", "", "Operation: e for encrypt, d for decrypt")
var InputFile = flag.String("i", "", "Input file path")
var OutputFile = flag.String("o", "", "Output file path")
var EncryptPass = flag.String("p", "", "Password String")
var encryptPassStr = ""

func init() {
	flag.Parse()
}

func main() {
	if len(*EncryptPass) != 32 {
		pwdHash := sha256.Sum256([]byte(*EncryptPass))
		encryptPassStr = string(pwdHash[:])
	}
	switch *Operation {
	case "e":
		DoCrypt("e", *InputFile, *OutputFile, encryptPassStr)
	case "d":
		DoCrypt("d", *InputFile, *OutputFile, encryptPassStr)
	default:
		log.Println("Illegal Operation.")
	}
}

func IsFileExists(fileLoc string) (string, bool) {
	fullP, err := filepath.Abs(fileLoc)
	if err != nil {
		return "", false
	}
	_, err = os.Stat(fullP)
	if err != nil {
		return "", false
	}
	return fullP, true
}

func DoCrypt(op string, iptfd string, optfd string, passwd string) {
	// check if input exists
	fullP, ok := IsFileExists(iptfd)
	if !ok {
		log.Fatalln("Input File not exists or cannot be accessed.")
	}
	// some prepartion
	oriFile, err := os.Open(fullP)
	if err != nil {
		log.Fatalln("File Open error.")
	}

	fullOptP, err := filepath.Abs(optfd)
	if err != nil {
		log.Fatalln("Output File Is incorrect.")
	}
	optFile, err := os.OpenFile(fullOptP, os.O_WRONLY|os.O_CREATE, 0644)
	defer optFile.Close()
	if err != nil {
		log.Fatalln("Output file is not writable.")
	}
	plaintext, err := ioutil.ReadAll(oriFile)
	if err != nil {
		log.Fatalln("File Read error.")
	}

	// crypt
	key := []byte(passwd)

	// cryptor initialize
	aescipher, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("AES Cannot initialize.")
	}

	// ctr
	var ciphertext []byte
	var iv []byte
	if op == "e" {
		ciphertext = make([]byte, aes.BlockSize+len(plaintext))
		iv = ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(crand.Reader, iv); err != nil {
			log.Fatalln(err)
		}
	} else if op == "d" {
		ciphertext = make([]byte, len(plaintext)-aes.BlockSize)
		iv = plaintext[:aes.BlockSize] // something like CTS
	} else {
		log.Fatalln("Unknown operation.")
	}
	ctrcipher := cipher.NewCTR(aescipher, iv) // not AEAD, no fail-safe

	if op == "e" {
		ctrcipher.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	} else if op == "d" {
		ctrcipher.XORKeyStream(ciphertext, plaintext[aes.BlockSize:])
	} else {
		log.Fatalln("Unknown operation.")
	}

	// write to file
	_, err = optFile.Write(ciphertext)
	if err != nil {
		log.Fatalln(err)
	}
	_ = optFile.Sync()
	// log done
	log.Println("Done.")
}
