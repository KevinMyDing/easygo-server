/**
作者:guangbo
模块：des加解密接口
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type Des struct {
	block cipher.Block
}

func NewDes(key []byte) *Des {
	des := new(Des)
	cblock, err := aes.NewCipher(key)
	if err != nil {
		panic("aes.NewCipher: " + err.Error())
	}

	des.block = cblock
	return des
}

// AES加密
func (d *Des) Encrypt(src []byte) ([]byte, error) {
	// 必须为aes.Blocksize的倍数
	if len(src)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	encryptText := make([]byte, aes.BlockSize+len(src))

	iv := encryptText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(d.block, iv)

	mode.CryptBlocks(encryptText[aes.BlockSize:], src)

	return encryptText, nil
}

// AES解密
func (d *Des) Decrypt(src []byte) ([]byte, error) {
	// hex
	decryptText, err := hex.DecodeString(fmt.Sprintf("%x", string(src)))
	if err != nil {
		return nil, err
	}

	// 长度不能小于aes.Blocksize
	if len(decryptText) < aes.BlockSize {
		return nil, errors.New("crypto/cipher: ciphertext too short")
	}

	iv := decryptText[:aes.BlockSize]
	decryptText = decryptText[aes.BlockSize:]

	// 必须为aes.Blocksize的倍数
	if len(decryptText)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(d.block, iv)

	mode.CryptBlocks(decryptText, decryptText)

	return decryptText, nil
}
