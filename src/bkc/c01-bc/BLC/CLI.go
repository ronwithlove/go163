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
	fmt.Println("Usage: \n")
	//初始化区块链
	fmt.Printf("\tcreateblockchain -address address -- 创建区块链\n")
	//添加区块
	//fmt.Printf("\taddblock -data DATA -- 添加区块\n")
	//打印完整的区块信息
	fmt.Printf("\tprintchain -- 打印区块链\n")
	//通过命令转账
	fmt.Printf("\tsend- from FROM -to TO -amount AMOUNT -- 发起转账\n")
	fmt.Printf("\t转账参数说明:\n")
	fmt.Printf("\t\t-from FROM -- 转账源地址\n")
	fmt.Printf("\t\t-from TO -- 转账目标地址\n")
	fmt.Printf("\t\t-AMOUNT amount -- 转账金额\n")
	//查询余额
	fmt.Printf("\tgetbalance -address FROM -- 查询指定地址的余额\n")
	fmt.Printf("\t查询余额参数说明：\n")
	fmt.Printf("\t-address -- 查询余额的地址\n")
	fmt.Println("\tcreatewallet -- 创建钱包\n")
	fmt.Println("\taccounts -- 获取钱包地址\n")
	fmt.Printf("\tutxo -method METHOD -- 测试UTXO Table 功能中指定的方法\n")
	fmt.Printf("\t\tMETHOD -- 方法名\n")
	fmt.Printf("\t\t\treset -- 重制UTXOtable\n")
	fmt.Printf("\t\t\tbalance - 查找所有UTXO")
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



//命令行运行函数
func (cli *CLI)Run(){
	//检测参数数量
	IsValidArgs()
	//新建相关命令
	//添加区块
	addBlockCmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	//创建钱包
	createWalletCmd:=flag.NewFlagSet("createwallet",flag.ExitOnError)
	//获取地址列表
	getAccountsCmd:=flag.NewFlagSet("accounts",flag.ExitOnError)
	//输出区块链完整信息
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	//创建区块链
	createBLCWithGenesisBlockCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	//发起交易
	sendCmd:=flag.NewFlagSet("send",flag.ExitOnError)
	//查询余额的命令
	getBalanceCmd:=flag.NewFlagSet("getbalance",flag.ExitOnError)
	//utxo测试命令
	UTXOTestCmd:=flag.NewFlagSet("utxo",flag.ExitOnError)
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
	//查询余额命令行参数
	flagGetBalanceArg:=getBalanceCmd.String("address","","要查询的地址")
	//UTXO测试命令行参数
	flagUTXOArg:=UTXOTestCmd.String("method","","UTXO Table相关操作\n")
	//判断命令
	switch os.Args[1] {//判断第二个命令
	case "utxo":
		err:=UTXOTestCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Printf("parse cmd operate utxo table fialed! %v\n",err)
		}
	case "accounts":
		err:=getAccountsCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Printf("parse cmd of get account fialed! %v\n",err)
		}
	case "createwallet":
		err:=createWalletCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Printf("parse cmd of create wallet fialed! %v\n",err)
		}
	case "getbalance":
		err:=getBalanceCmd.Parse(os.Args[2:])
		if err!=nil {
			log.Printf("parse cmd get Balance  failed! %v\n",err)
		}
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

	//utxo table 操作
	if UTXOTestCmd.Parsed(){
		switch *flagUTXOArg{
		case "balance":
			cli.TestFindUTXOMap()
		case"reset":
			cli.TestResetUTXO()
		default:


		}
	}
	//获取地址列表
	if getAccountsCmd.Parsed(){
		cli.GetAccounts()
	}
	//创建钱包
	if createWalletCmd.Parsed(){
		cli.CreateWallets()
	}
	//查询余额
	if getBalanceCmd.Parsed(){
		if *flagGetBalanceArg==""{
			fmt.Println("请输入查询地址")
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
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
		cli.send(JSONToSlice(*flagSendFromArg),JSONToSlice(*flagSendToArg),JSONToSlice(*flagSendAmountArg))
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