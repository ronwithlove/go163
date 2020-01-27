package BLC

import (
	"fmt"
	"os"
)

//打印完整区块链信息
func(cli *CLI) printchain(){
	//判断数据库是否存在
	if !dbExist(){
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain:=BlockchainObject()//获取到最新的blockchain的对象实例
	blockchain.PrintChain()
}
