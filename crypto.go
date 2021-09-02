package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"bufio"
	"os"
	"io"
)

func AesEncryptStr(origData, key []byte) (string, error) {

	result, err := AesEncrypt(origData, key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), nil

}
func AesDecryptStr(crypted string, key []byte) ([]byte, error) {

	destr, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return nil, err
	}

	origData, err := AesDecrypt(destr, key)
	if err != nil {
		return nil, err
	}

	return origData, nil

}
func AesEncryptBase62(origData, key []byte) (string, error) {

	result, err := AesEncrypt(origData, key)
	if err != nil {
		return "", err
	}

	uEnc := base64.URLEncoding.EncodeToString([]byte(result))
	return uEnc, nil //base64.StdEncoding.EncodeToString(result), nil

}
func AesDecryptBase62(crypted string, key []byte) ([]byte, error) {

	destr, err := base64.URLEncoding.DecodeString(crypted) //base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return nil, err
	}

	origData, err := AesDecrypt(destr, key)
	if err != nil {
		return nil, err
	}

	return origData, nil

}
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Md5(buf []byte) string {
	//hash := md5.New()
	//hash.Write(buf)
	//return fmt.Sprintf("%x", hash.Sum(nil))
	ret := md5.Sum(buf)
	return hex.EncodeToString(ret[:])
}
func Sha1s(s string) string {
	r := sha1.Sum([]byte(s))
	sh := hex.EncodeToString(r[:])
	return strings.ToUpper(sh)
}
func FileHash(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {

		return "", err
	}
	defer f.Close()

	br := bufio.NewReader(f)

	h := sha1.New()
	_, err = io.Copy(h, br)

	if err != nil {
	
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
