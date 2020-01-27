package BLC

import (
	"fmt"
	"os"
)

//发起交易
func (cli *CLI) send(from, to , amount []string){
	if !dbExist(){
		fmt.Printf("数据库不存在")
		os.Exit(1)
	}
	//获取区块链对象
	blockchain:=BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to , amount)
}