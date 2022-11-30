package crypto

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/tjfoc/gmsm/sm4"
)

// Sm4StrEncrypt 返回base64后的加密串
func Sm4StrEncrypt(key, data string) (string, error) {
	keyBytes := []byte(key)
	dataBytes := []byte(data)
	fmt.Println("明文串:", data)
	iv := make([]byte, sm4.BlockSize)
	ciphertxt, err := sm4Encrypt(keyBytes, iv, dataBytes)
	if err != nil {
		fmt.Println("sm4加密异常:", err)
		return "", err
	}
	base64Result := base64.StdEncoding.EncodeToString(ciphertxt)
	fmt.Println("加密串:", base64Result)
	return base64Result, nil
}

// Sm4StrDecrypt 传入base64串
func Sm4StrDecrypt(key, base64Str string) (string, error) {
	keyBytes := []byte(key)
	fmt.Println("加密串:", base64Str)
	base64Data, _ := base64.StdEncoding.DecodeString(base64Str)
	iv := make([]byte, sm4.BlockSize)
	decbytes, err := sm4Decrypt(keyBytes, iv, []byte(base64Data))
	if err != nil {
		fmt.Println("sm4解密异常:", err)
		return "", err
	}
	oriData := string(decbytes)
	fmt.Println("明文串:", oriData)
	return oriData, nil
}

func sm4Encrypt(key, iv, plainText []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData := pkcs5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

func sm4Decrypt(key, iv, cipherText []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

// pkcs5填充
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return nil
	}
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

/*
func mytest(){
	// 128比特密钥
	key := []byte("1234567890abcdef")
	// 128比特iv
	pwdIV := make([]byte, sm4.BlockSize)
	data := []byte("golden")
	fmt.Printf("原文: %s\n", string(data))//原文: golden
	ciphertxt,_ := sm4Encrypt(key,pwdIV, data)
	fmt.Printf("密文: %x\n", ciphertxt)//密文: 816517345af210927d523c0c8cea0c29
	str := base64.StdEncoding.EncodeToString(ciphertxt)
	fmt.Printf("base64加密结果: %s\n", str)//base64加密结果: gWUXNFryEJJ9UjwMjOoMKQ==

	//解密
	str2,_ := base64.StdEncoding.DecodeString(str)
	fmt.Printf("解base64后: %x\n", str2)//解base64后: 816517345af210927d523c0c8cea0c29
	decbytes,_ := sm4Decrypt(key,pwdIV,[]byte(str2))
	fmt.Println("解密结果:", string(decbytes)) //解密结果: golden
}
*/
