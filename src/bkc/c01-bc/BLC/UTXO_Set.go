package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//UTXO持久化相关管理

//用于存入utxo的bucket
const  utxoTableName  = "utxoTable"

//utxoSet结构（保存指定区块链中所有的UTXO）
type UTXOSet struct{
	Blockchain *BlockChain
}

//序列化 txOutputs=>[]bytes
func (txOutputs *TXOutputs) Serilize()[]byte{
	var result bytes.Buffer

	encoder:=gob.NewEncoder(&result)
	if err:=encoder.Encode(txOutputs);nil!=err{
		log.Printf("serialize the utxo failed! %v\n",err)
	}
	return result.Bytes();
}

//反序列化 []bytes=>txOutputs
func DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs{
	var txOutputs TXOutputs
	decoder:=gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	if err:=decoder.Decode(&txOutputs);nil!=err{
		log.Panicf("deserialize the struct utxo failed! %v\n",err)
	}
	return &txOutputs
}


//更新


//查询余额
func (utxoSet *UTXOSet) GetBalance(address string)int{
	UTXOS:=utxoSet.FindUTXOWithAddress(address)
	var amount int
	for _, utxo:=range UTXOS{
		fmt.Printf("utxo-txhash:%x\n",utxo.TxHash)
		fmt.Printf("utxo-Index:%x\n",utxo.Index)
		fmt.Printf("utxo-Ripemd160Hash:%x\n",utxo.Output.Ripemd160Hash)
		fmt.Printf("utxo-Value:%x\n",utxo.Output.Value)
		amount+=utxo.Output.Value
	}
	return  amount
}


//查找
func(utxoSet *UTXOSet)FindUTXOWithAddress(address string)[]*UTXO{
	var utxos []*UTXO
	err:=utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		//1.获取utxotable表
		b:=tx.Bucket([]byte(utxoTableName))
		if nil!=b{
			c:=b.Cursor()//cursor--游标
			//通过游标遍历boltdb数据库中的数据
			for k,v:=c.First();k!=nil;k,v=c.Next(){
				txOutputs:=DeserializeTXOutputs(v)//从数据库得到utxo
				for _,utxo:=range txOutputs.TXOutputs{//遍历，看看是否匹配
					if utxo.UnLockScriptPubkeyWithAddress(address){
						utxo_signle:=UTXO{Output:utxo}
						utxos=append(utxos, &utxo_signle)
					}

				}
			}
		}
		return nil
	})
	if nil!=err{
		log.Printf("find the utxo of [%s] failed! %v\n",address,err)
	}
	return utxos
}




//重制,把UTXO存到DB里
func (utxoSet *UTXOSet)ResetUTXOSet()  {
	//在第一次创建的时候就更新utxo table
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		//查找utxo table
		b:=tx.Bucket([]byte(utxoTableName))
		if nil!=b{
			err:=tx.DeleteBucket([]byte(utxoTableName))//已经存在就删除
			if nil!=err{
				log.Panicf("delete the utxo table failed! %\n",err)
			}
		}
		//创建
		bucket,err:=tx.CreateBucket([]byte(utxoTableName))
		if nil!=err{
			log.Printf("create bucket failed! %v\n",err)
		}
		if nil!=bucket{
			//查找当前所有UTXO
			txOutputMap:=utxoSet.Blockchain.FindUTXOMap();
			for keyHash, outputs:=range txOutputMap{
				//将所有UTXO存入
				txHash, _:=hex.DecodeString(keyHash)//txHash序列化
				fmt.Printf("txHash: %x\n",txHash)

				//存入utxo talb,序列话之后，存入bucket
				err:=bucket.Put(txHash,outputs.Serilize())
				if nil!=err{
					log.Printf("put the utxo into table failed! %v\n",err)
				}
			}
		}
		return nil
	})
}

//更新
func (utxoSet *UTXOSet) update(){
	//获取最新区块
	lastest_block:=utxoSet.Blockchain.Iterator().Next()
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(utxoTableName))
		if nil!=b{
			//只需查找最新一个区块的交易列表，因为每上链一个区块
			//utxo table都更新一次， 所以只需要查找最近一个区块中的交易
			for _,tx:=range lastest_block.Txs{
				if !tx.IsCoinbaseTransaction(){
					//2.将已经被当前这笔交易的输入所引用的UTXO删除
					for _,vin:=range tx.Vins{
						//需要更新的输出
						updatedOutputs:=TXOutputs{}
						//获取指定输入所引用的交易哈希的输出
						outputBytes:=b.Get(vin.TxHash)
						//输出列表
						outs:=DeserializeTXOutputs(outputBytes)
						for outIdex,out:=range outs.TXOutputs{
							if vin.Vout!=outIdex{
								updatedOutputs.TXOutputs=append(updatedOutputs.TXOutputs,out)
							}
						}
						//如果交易中没有UTXO了，删除该交易
						if len(updatedOutputs.TXOutputs)==0{
							b.Delete(vin.TxHash)
						}else{
							//把更新之后的utxo数据存入数据库
							b.Put(vin.TxHash,updatedOutputs.Serilize())
						}
					}

				}
				//获取当前区块中新生成的交易输出
				//1.将最新区块中的UTXO插入
				newOutputs:=TXOutputs{}
				newOutputs.TXOutputs=append(newOutputs.TXOutputs,tx.Vouts...)
				b.Put(tx.TxHash,newOutputs.Serilize())
			}
		}
		return nil
	})
}