package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

//对blockchain的命令操作进行管理

//client对象
type CLI struct {
	//不需要
	//BC *BlockChain		//blockchain对象
}


//用法展示
func PrintUsage(){
	fmt.Println("Usage:")
	//初始化区块链
	fmt.Printf("\tcreateblockchain -- 创建区块链\n")
	//添加区块
	fmt.Printf("\taddblock -data DATA -- 添加区块\n")
	//打印完整的区块信息
	fmt.Printf("\tprintchain -- 打印区块链\n")
}

//初始化区块链
func(cli *CLI) createBlockchain(){
	CreateBlockCHainWithGenesisBlock()
}

//添加区块
func (cli *CLI) addBlock(data string){
	//判断数据库是否存在
	if !dbExist(){
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain:=BlockchainObject()//获取到最新的blockchain的对象实例
	blockchain.AddBlock([]byte(data))//新加区块
	//cli.BC.AddBlock([]byte(data))//删除，没有必要了
}

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

//参数数量检测函数
func IsValidArgs(){
	if len(os.Args)<2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}

//命令行运行函数
func (cli *CLI)Run(){
	//检测参数数量
	IsValidArgs()
	//新建相关命令
	//添加区块
	addBlockCmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	//输出区块链完整信息
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	//创建区块链
	createBLCWithGenesisBlockCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)

	//数据参数处理
	flagAddBlockArg:=addBlockCmd.String("data","sent 100 btc to user","添加区块数据")

	//判断命令
	switch os.Args[1] {//判断第二个命令
	case "addblock":
		err:=addBlockCmd.Parse(os.Args[2:])
		if err!=nil {
			log.Printf("parse addBlockCmd failed! %v\n",err)
		}
	case "printchain":
		err:=printChainCmd.Parse(os.Args[2:])
		if err!=nil {
			log.Printf("parse printChainCmd failed! %v\n",err)
		}
	case "createblockchain":
		err:=createBLCWithGenesisBlockCmd.Parse(os.Args[2:])
		if err!=nil {
			log.Printf("parse createBLCWithGenesisBlockCmd failed! %v\n",err)
		}
	default:
		//没有传递以上命令
		PrintUsage()
		os.Exit(1)

	}

	//添加区块命令
	if addBlockCmd.Parsed(){
		if *flagAddBlockArg==""{//如果没有输入参数
		PrintUsage()
		os.Exit(1)//直接退出
		}
		cli.addBlock(*flagAddBlockArg)//如果有就把参数传进去
	}
	//输出区块链信息
	if printChainCmd.Parsed(){
		cli.printchain()
	}
	//创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed(){
		cli.createBlockchain()
	}
}