package main

import (
	"bkc/c01-bc/BLC"
)

func main(){

	bc:=BLC.CreateBlockCHainWithGenesisBlock()
	defer bc.DB.Close()
	//上链
	bc.AddBlock([]byte("ron sent 10 tc to aaron"))
	bc.AddBlock([]byte("jacky sent 10 tc to aaron"))

	bc.PrintChain()

}