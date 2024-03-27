package encrypts

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	plain := "1234567"
	// AES 规定有3种长度的key: 16, 24, 32分别对应AES-128, AES-192, or AES-256
	key := "sdfgyrhgbxcdgryfhgywertd"
	// 加密
	cipherByte, err := Encrypt(plain, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %s\n", plain, cipherByte)
	// 解密
	plainText, err := Decrypt(cipherByte, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %s\n", cipherByte, plainText)
}

func TestMd5(t *testing.T) {
	plain := "1234567"
	md5 := Md5(plain)
	fmt.Println(md5)
}
