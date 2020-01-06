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
	fmt.Printf("\tcreateblockchain -address address -- 创建区块链\n")
	//添加区块
	fmt.Printf("\taddblock -data DATA -- 添加区块\n")
	//打印完整的区块信息
	fmt.Printf("\tprintchain -- 打印区块链\n")
	//通过命令转账
	fmt.Printf("\tsend- from FROM -to TO -amount AMOUNT -- 发起转账\n")
	fmt.Printf("\t转账参数说明:\n")
	fmt.Printf("\t\t-from FROM -- 转账源地址\n")
	fmt.Printf("\t\t-from TO -- 转账目标地址\n")
		fmt.Printf("\t\t-AMOUNT amount -- 转账金额\n")
}

//参数数量检测函数
func IsValidArgs(){
	if len(os.Args)<2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}

//发起交易
func (cli *CLI) send(){

}
//初始化区块链
func(cli *CLI) createBlockchain(address string){
	CreateBlockCHainWithGenesisBlock(address)
}

//添加区块
func (cli *CLI) addBlock(txs []*Transaction){
	//判断数据库是否存在
	if !dbExist(){
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain:=BlockchainObject()//获取到最新的blockchain的对象实例
	blockchain.AddBlock(txs)//新加区块
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
	//发起交易
	sendCmd:=flag.NewFlagSet("send",flag.ExitOnError)
	//数据参数处理
	//1.添加区块
	flagAddBlockArg:=addBlockCmd.String("data","sent 100 btc to user","添加区块数据")
	//2.创建区块链指定的矿工地址（矿工接收奖励）
	flagCreateBlockchainArg:=createBLCWithGenesisBlockCmd.String("address",
		"troytan","指定接收系统奖励的矿工地址")
	//发起交易参数
	flagSendFromArg:= sendCmd.String("from","","转账源地址")
	flagSendToArg:= sendCmd.String("to","","转账目标地址")
	flagSendAmountArg:= sendCmd.String("amount","","转账金额")

	//判断命令
	switch os.Args[1] {//判断第二个命令
	case "send":
		err:=sendCmd.Parse(os.Args[2:])
		if err!=nil {
			log.Printf("parse sendCmd failed! %v\n",err)
		}
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
	//发起转账
	if sendCmd.Parsed(){
		if *flagSendFromArg==""{//如果没有输入参数
			fmt.Printf("源地址不能为空")
			PrintUsage()
			os.Exit(1)//直接退出
		}
		if *flagSendToArg==""{//如果没有输入参数
			fmt.Printf("转账目标地址不能为空")
			PrintUsage()
			os.Exit(1)//直接退出
		}
		if *flagSendAmountArg==""{//如果没有输入参数
			fmt.Printf("转账金额不能为空")
			PrintUsage()
			os.Exit(1)//直接退出
		}
		//先打印测试一下看看
		fmt.Printf("\tFROM:[%s]\n",JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n",JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n",JSONToSlice(*flagSendAmountArg))

	}
	//添加区块命令
	if addBlockCmd.Parsed(){
		if *flagAddBlockArg==""{//如果没有输入参数
		PrintUsage()
		os.Exit(1)//直接退出
		}
		cli.addBlock([]*Transaction{})//如果有就把参数传进去 //暂时穿个空的
	}
	//输出区块链信息
	if printChainCmd.Parsed(){
		cli.printchain()
	}
	//创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed(){
		if *flagCreateBlockchainArg==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchainArg)//
	}
}