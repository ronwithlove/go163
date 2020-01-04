package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//共识算法管理文件

//实现POW实例以及相关功能

//目标难度值
const targetBit=16//前两位为0，一个bit8位
//工作量证明的结构
type ProofOfWork struct{
	//需要共识验证的区块
	Block *Block
	//目标难度的哈希
	target	*big.Int
}

//创建一个POW对象
func NewProofOfWork(block *Block) *ProofOfWork{
	target:=big.NewInt(1)
	//数据长度为8位
	//需求：需求满足前两位为0，才能解决问题
	//1*2<<(8-2)=64     64 0100 0000
	//00xx xxxx  就算全是0011 1111 =63
	//sha256.Sum256返回的事32位的字节数bytes,=32*8bit= 256bit,所以下面事256
	target=target.Lsh(target,256-targetBit)//左移计算
	return &ProofOfWork{Block:block,target:target}
}

//执行pow,比较哈希
//返回哈希值 碰撞次数
func(ProofOfWork *ProofOfWork) Run()([]byte, int){
	//碰撞次数
	var nonce=0
	var hashInt big.Int
	var hash [32]byte //生成的哈希值
	//无限循环，生成符合条件的哈希值
	for{
		//生成准备数据
		dataBytes:=ProofOfWork.prepareData(int64((nonce)))
		hash =sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//检测生成的哈希值是否符合条件
		if ProofOfWork.target.Cmp(&hashInt)==1{// ==1：前面的big.Int 实例大于cmp方法里的hashInt 参数
			//找到了符合条件的哈希值，中断循环
			break
		}
		nonce++
	}
	fmt.Printf("\n碰撞次数：%d\n",nonce)
	return hash[:],nonce
}

//生成准备数据
func (pow *ProofOfWork)prepareData(nonce int64) []byte{
	var data []byte
	//拼接区块属性，进行哈希计算
	timeStampBytes:=IntToHex(pow.Block.TimeStamp)//把int转成byte
	heightBytes:=IntToHex(pow.Block.Heigth)
	data = bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		IntToHex(nonce),
		IntToHex(targetBit),
	},[]byte{})
	return data
}


