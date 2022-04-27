package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

var (
	aesGCM   cipher.AEAD
	aesBlock cipher.Block
	nonce    []byte
)

func init() {
	var err error

	key, errKey := generateRandom(2 * aes.BlockSize)
	if errKey != nil {
		fmt.Println(errKey.Error())
	}

	aesBlock, err = aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	aesGCM, err = cipher.NewGCM(aesBlock)
	if err != nil {
		fmt.Println(err.Error())
	}

	nonce, err = generateRandom(aesGCM.NonceSize())
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Encode(userID string) (string, error) {
	src := []byte(userID)

	dst := aesGCM.Seal(nil, nonce, src, nil)

	sha := hex.EncodeToString(dst)

	return sha, nil
}

func Decode(sha string, userID *string) error {
	dst, err := hex.DecodeString(sha)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	src, errGCM := aesGCM.Open(nil, nonce, dst, nil)

	if errGCM != nil {
		return errGCM
	}

	*userID = string(src)

	return nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
