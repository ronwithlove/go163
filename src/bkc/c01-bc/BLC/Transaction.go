package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//交易管理文件

//定义一个交易基本结构
type Transaction struct{
	//交易哈希（标识）
	TxHash     []byte
	//输入列表
	Vins 		[]*TxInput
	//输出列表
	Vous 		[]*TxOutput
}

//实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction{
	//输入,txHash为nil,这里的索引用-1，挖坑没有人就用system reward
	txInput:=&TxInput{[]byte{},-1,"system reward"}
	//输出,value暂定给10,address
	txOutput:=&TxOutput{10,address}

	txCoinbase:=&Transaction{
		nil,
		[]*TxInput{txInput},
		[]*TxOutput{txOutput},
	}
	//交易哈希生成
	txCoinbase.HashTransaction()
	return txCoinbase
}

//生产交易哈希（交易序列化）
func (tx *Transaction) HashTransaction(){
	var result bytes.Buffer
	//设置编码对象
	encoder:=gob.NewEncoder(&result)
	if err:=encoder.Encode(tx);err!=nil{
		log.Panicf("tx Hash encoded failed %v\n",err)
	}
	//生成哈希值
	hash:=sha256.Sum256(result.Bytes())
	tx.TxHash=hash[:]
}