package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//区块基本结构与功能管理文件

//实现一个最基本的区块结构
type Block struct{
	TimeStamp 		int64	//区块时间戳，代表区块时间
	Hash			[]byte	//当前区块哈希
	PrevBlockHash 	[]byte	//前区块哈希
	Heigth			int64	//区块高度
	Txs             []*Transaction //交易数据
	//Data			[]byte	//交易数据
	Nonce			int64	//在运行pow时生成的哈希变化值，也代表pow运行的动态修改的数据
}

//新建区块
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block{
	var block Block

	block=Block{
		TimeStamp:time.Now().Unix(),
		Hash:nil,//这个是算出来的，现在先nil
		PrevBlockHash:prevBlockHash,
		Heigth:height,
		Txs:txs,
	}
	//通过POW 生成新的哈希值
	pow:= NewProofOfWork(&block)
	//执行工作量证明算法
	hash,nonce:=pow.Run()
	block.Hash=hash
	block.Nonce=int64(nonce)
	return &block
}

//生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block{
	return NewBlock(1,nil,txs)//区块高度从1开始
}

//boltDB存储的键值对的数据类型都是字节数组。所以在存储区块前需要对区块进行序列化，先转成字节数组[]byte
func (block *Block) Serialize() []byte  {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil{
		log.Panicf("serialize the block to []byte failed! %v \n",err)
	}
	return result.Bytes()
}

//读取区块的时候就需要做反序列化处理
func DeSerializeBlock(blockBytes []byte) *Block  {

	var block Block
	dencoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := dencoder.Decode(&block)
	if err != nil{
		log.Panicf("deserialize the []byte to block failed! %v\n",err)
	}

	return &block
}

//把指定区块所有交易结构都序列化(类Merkle的哈希计算方法)
func (block *Block) HashTransaction() []byte{
	var txHashes [][]byte
	//将指定区块中所有交易哈希进行拼接
	for _,tx:= range block.Txs{
		txHashes=append(txHashes,tx.TxHash)
	}
	txHash:=sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return txHash[:]
	//将交易数据存入Merkle树中，然后生成Merkle更节点
	mtree:=NewMerkleTree(txHashes)
	return mtree.RootNode.Data
}

