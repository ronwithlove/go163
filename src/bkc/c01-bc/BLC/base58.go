package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

//base64:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/
//base58:去掉0(零)，O(大写的 o)，I(大写的i)，l(小写的 L)，+，/

//base58编码
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte{
	var result []byte

	x:= big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)

	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x,base,mod)  // 对x取余数
		result =  append(result, b58Alphabet[mod.Int64()])
	}



	ReverseBytes(result)

	for _,b:=range input{

		if b ==0x00{
			result =  append([]byte{b58Alphabet[0]},result...)
		}else{
			break
		}
	}


	return result

}


//字节数组的反转
func ReverseBytes(data []byte){
	for i,j :=0,len(data) - 1;i<j;i,j = i+1,j - 1{
		data[i],data[j] = data[j],data[i]
	}
}
func generatePrivateKey(hexprivatekey string,compressed bool) []byte{
	versionstr :=""
	//判断是否对应的是压缩的公钥，如果是，需要在后面加上0x01这个字节。同时任何的私钥，我们需要在前方0x80的字节
	if compressed{
		versionstr  = "80" + hexprivatekey + "01"
	}else{
		versionstr  = "80" + hexprivatekey
	}
	//字符串转化为16进制的字节
	privatekey,_:=hex.DecodeString(versionstr)
	//通过 double hash 计算checksum.checksum他是两次hash256以后的前4个字节。
	firsthash:=sha256.Sum256(privatekey)

	secondhash:= sha256.Sum256(firsthash[:])

	checksum := secondhash[:4]

	//拼接
	result := append(privatekey,checksum...)

	//最后进行base58的编码
	base58result :=Base58Encode(result)
	return base58result
}



func Base58Decode(input []byte) []byte{
	result :=  big.NewInt(0)
	zeroBytes :=0
	for _,b :=range input{
		if b=='1'{
			zeroBytes++
		}else{
			break
		}
	}

	payload:= input[zeroBytes:]

	for _,b := range payload{
		charIndex := bytes.IndexByte(b58Alphabet,b)  //反推出余数

		result.Mul(result,big.NewInt(58))   //之前的结果乘以58

		result.Add(result,big.NewInt(int64(charIndex)))  //加上这个余数

	}

	decoded :=result.Bytes()


	decoded =  append(bytes.Repeat([]byte{0x00},zeroBytes),decoded...)
	return decoded
}
