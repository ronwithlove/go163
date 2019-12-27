package BLC

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//区块基本结构与功能管理文件

//实现一个最基本的区块结构
type Block struct{
	TimeStamp 		int64	//区块时间戳，代表区块时间
	Hash			[]byte	//当前区块哈希
	PrevBlockHash 	[]byte	//前区块哈希
	Heigth			int64	//区块高度
	Data			[]byte	//交易数据
}

//新建区块
func NewBlock(height int64, prevBlockHash []byte, data []byte) *Block{
	var block Block

	block=Block{
		TimeStamp:time.Now().Unix(),
		Hash:nil,//这个是算出来的，现在先nil
		PrevBlockHash:prevBlockHash,
		Heigth:height,
		Data:data,
	}
	//生成哈希
	block.SetHash()
	return &block
}

//计算区块哈希，方法，因为只和Block有关
func(b *Block) SetHash(){//用指针不需要返回值，要不然太多余
	//调用sha256实现哈希生成
	timeStampBytes:=IntToHex(b.TimeStamp)//把int转成byte
	heightBytes:=IntToHex(b.Heigth)
	blockBytes:=bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		b.PrevBlockHash,
		b.Data,
	},[]byte{})//将一系列[]byte切片连接为一个[]byte切片，之间用sep来分隔，返回生成的新切片。
	hash:=sha256.Sum256(blockBytes)
	b.Hash=hash[:]//赋值切片所有内容
}

//生成创世区块
func CreateGenesisBlock(data []byte) *Block{
	return NewBlock(1,nil,data)//区块高度从1开始
}