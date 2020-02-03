package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
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
	tx.HashTransaction()//完整的交易生成
	//对交易进行签名
	bc.SignTransaction(&tx,wallet.PrivateKey)
	return &tx
}


//判断指定的交易是否是一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool{
	return  tx.Vins[0].Vout==-1 && len(tx.Vins[0].TxHash)==0//满足&&前一个判断也就是coinbase交易了
}

//交易签名
//prevTxs:代表当前交易的输入所引用的所有OUTPUT所属的交易
func (tx *Transaction)Sign(privateKey ecdsa.PrivateKey,prevTxs map[string]Transaction){
	//处理输入，保证交易的正确性
	//检查当前交易tx中每一个Input的哈希是否能在前面的preTxs中找到（目测这里并不能解决双花）
	//如果没有包含在里面，说明该交易被人修改了
	for _,vin:=range tx.Vins{
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash==nil{
			log.Panicf("ERROR:Prev transaction is not correct!\n")
		}
	}
	//提取需要签名的属性
	txCopy:=tx.TrimmedCopy()
	//处理交易副本的输入
	for vin_id, vin:=range txCopy.Vins{
		//获取关联交易
		prevTx:=prevTxs[hex.EncodeToString(vin.TxHash)]
		//找到发送者(当前输入引用的哈希--输出的哈希)
		//vin.Vout是交易Input里写的前交易内的index，这里使用原始的output做哈希
		txCopy.Vins[vin_id].PublicKey=prevTx.Vouts[vin.Vout].Ripemd160Hash
		//生成交易副本的哈希
		txCopy.TxHash=txCopy.Hash()
		//调用核心签名函数
		r,s,err:=ecdsa.Sign(rand.Reader,&privateKey,txCopy.TxHash)
		if nil!=err{
			log.Printf("sign to transaction [%x] failed! %v\n",err)
		}

		//组成交易签名
		signature:=append(r.Bytes(),s.Bytes()...)
		tx.Vins[vin_id].Signature=signature
	}

}

//交易拷贝，生成一个专门用于交易签名的副本
func(tx *Transaction) TrimmedCopy()Transaction{
	//重新组装生成一个新的交易
	var inputs []*TxInput
	var outputs []*TxOutput
	//组装input
	for _, vin:=range tx.Vins{
		inputs=append(inputs,&TxInput{vin.TxHash,vin.Vout,
			nil,nil})
	}
	//组装output
	for _,vout:=range tx.Vouts{
		outputs=append(outputs,&TxOutput{vout.Value,vout.Ripemd160Hash})
	}
	txCopy:=Transaction{tx.TxHash,inputs,outputs}
	return txCopy
}


//设置用于签名的交易的哈希
func (tx *Transaction)Hash()[]byte{
	txCopy :=tx
	txCopy.TxHash=[]byte{}
	//这里没有直接把tx.TxHash滞空，因为tx是引用指针的，会改变原有tx
	hash:=sha256.Sum256(txCopy.Serialize())
	return  hash[:]
}

//交易序列化
func (tx *Transaction) Serialize() []byte  {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil{
		log.Panicf("serialize the tx to []byte failed! %v \n",err)
	}
	return buffer.Bytes()
}