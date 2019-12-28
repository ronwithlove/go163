package main

import (
	"bkc/c01-bc/BLC"
	"fmt"
	"github.com/boltdb/bolt"
)

func main(){

	bc:=BLC.CreateBlockCHainWithGenesisBlock()
	defer bc.DB.Close()
	//上链
	bc.AddBlock([]byte("ron sent 10 tc to aaron"))
	bc.AddBlock([]byte("jacky sent 10 tc to aaron"))

	bc.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte("blocks"))
		if nil!=b{
			hash:=b.Get([]byte("1"))
			fmt.Printf("value: %x\n",hash)
		}
		return nil
	})


}