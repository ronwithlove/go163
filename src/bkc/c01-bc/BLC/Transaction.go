package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
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
	Vouts []*TxOutput
}

//实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction{
	//输入,txHash为nil,这里的索引用-1，挖坑没有人就用system reward
	txInput:=&TxInput{[]byte{},-1,nil,nil}
	//输出,value暂定给10,address
	//txOutput:=&TxOutput{10,StringToHash160(address)}
	txOutput:=NewTxOutput(10,address)
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

//生成普通转账交易
func NewSimpleTransaction(from string, to string, amount int,bc *BlockChain,txs []*Transaction) *Transaction{
	var txInputs []*TxInput
	var txOutupts []*TxOutput

	//调用可花费UTXO函数
	money,spendableUTXODic:=bc.FindSpendableUTXO(from,amount,txs)
	fmt.Printf("money:%v\n",money)
	//获取钱包集合对象
	wallets:=NewWallets()
	//查找对应的钱包结构
	wallet:=wallets.Wallets[from]

	//输入
	for txHash,indexArray:=range spendableUTXODic{
		txHashesBytes,err:=hex.DecodeString(txHash)
		if nil!=err{
			log.Panicf("decode string to []byte failed! %v\n",err)
		}
		//遍历索引列表
		for _, index:=range indexArray{
			txInput:=&TxInput{txHashesBytes,index,nil,wallet.PublicKey}
			txInputs=append(txInputs,txInput)
		}
	}

	//输出
	//txOutput :=&TxOutput{
	//	Value:        amount,
	//	ScriptPubkey: to,
	//}
	txOutput:=NewTxOutput(amount,to)
	txOutupts=append(txOutupts, txOutput) //追加到输出交易中
	//输出（找零）
	if money>amount{
		//txOutput =&TxOutput{money-amount,from}//找零，找回给自己
		txOutput=NewTxOutput(money-amount,from)//找零，找回给自己
		txOutupts=append(txOutupts, txOutput) //再把这笔交易追加到输出交易中
	}else{
		log.Panicf("余额不足...\n")
	}
	tx:=Transaction{nil,txInputs,txOutupts}
	tx.HashTransaction()
	return &tx
}


//判断指定的交易是否是一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool{
	return  tx.Vins[0].Vout==-1 && len(tx.Vins[0].TxHash)==0//满足&&前一个判断也就是coinbase交易了
}
