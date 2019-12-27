package main

import (
	"bkc/c01-bc/BLC"
	"fmt"
)

func main(){

	bc:=BLC.CreateBlockCHainWithGenesisBlock()
	fmt.Printf("blockCHain:%v\n",bc.Blocks[0])

	//上链
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Heigth+1,
		bc.Blocks[len(bc.Blocks)-1].Hash,
		[]byte("ron sent 10 tc to aaron"))//长度和索引差1
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Heigth+1,
		bc.Blocks[len(bc.Blocks)-1].Hash,
		[]byte("jacky sent 10 tc to aaron"))//长度和索引差1

		for _,block:=range bc.Blocks{
			//fmt.Printf("block : %v\n",block)
			fmt.Printf("上一个区块哈希：%x, 当前区块哈希: %x\n",block.PrevBlockHash,block.Hash)
	}
}