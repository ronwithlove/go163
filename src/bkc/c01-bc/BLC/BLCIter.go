package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代管理文件

//迭代器基本结构
type BlockChainIterator struct {
	DB				*bolt.DB	//迭代目标
	CurrentHash		[]byte		//当前迭代目标的哈希
}



//实现迭代函数next, 获取每一个区块
func (bcit *BlockChainIterator) Next() *Block{
	var block *Block
	err:=bcit.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			currentBlockBytes:=b.Get(bcit.CurrentHash)
			block=DeSerializeBlock(currentBlockBytes)
			//更新迭代器中区块的哈希值
			bcit.CurrentHash=block.PrevBlockHash
		}
		return nil
	})
	if nil!=err{
		log.Printf("iterator the DB failed! %v\n",err)
	}
	return block
}