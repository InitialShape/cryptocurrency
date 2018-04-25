package utils

import (
	"golang.org/x/crypto/ed25519"
	"github.com/mr-tron/base58/base58"
	"crypto/rand"
	"os"
	"fmt"
	"log"
)

func GenerateWallet() error {
		path := "wallet.txt"
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			var file, err = os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()

			file, err = os.OpenFile(path, os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			defer file.Close()

			publicKeyBase58 := base58.Encode(publicKey)
			privateKeyBase58 := base58.Encode(privateKey)
			fileEntry := fmt.Sprintf("%s\n%s", publicKeyBase58, privateKeyBase58)
			_, err = file.WriteString(fileEntry)
			if err != nil {
				return err
			}

			err = file.Sync()
			if err != nil {
				return err
			}
			log.Println("New wallet.txt file created")
			return err
		} else {
			log.Println("Didn't create new key file as wallet.txt already exists")
			return nil
		}
}
