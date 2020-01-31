package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main(){
	//生成钥匙
	privateKey,err:=ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
	if nil!=err{
		panic(err)
	}
	//私钥+哈希（原文）=签名
	msg:="hello, world"//原文
	hash:=sha256.Sum256([]byte(msg))
	r,s,err:=ecdsa.Sign(rand.Reader,privateKey,hash[:])
	if nil!=err{
		panic(err)
	}
	fmt.Printf("signature:(0x%x,0x%x)\n",r,s)

	//验证签名：公钥+哈希（原文）+签名
	valid:=ecdsa.Verify(&privateKey.PublicKey,hash[:],r,s)
	fmt.Printf("the result of signature verify :%v\n",valid)
}
