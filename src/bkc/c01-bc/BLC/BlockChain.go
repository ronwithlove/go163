package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

//相关数据库属性
const dbName = "block.db"//数据库名
const blockTableName = "blocks"//表名
//区块链管理文件
type BlockChain struct{//直接用切片也可以，但是结构体比较正式一点
	//Blocks []*Block //区块的切片
	DB 		*bolt.DB	//数据库对象
	Tip 	[]byte		//最新区块的哈希值
}


//判断数据库文件是否存在
func dbExist() bool{
	if _,err:= os.Stat(dbName);os.IsNotExist(err){
		//数据库文件不存在
		return false
	}
	return  true
}

//初始化区块链
func CreateBlockCHainWithGenesisBlock(address string) *BlockChain{
	if dbExist(){//如果数据库已经存在
		fmt.Println("创世区块已存在")
		os.Exit(1)
	}

	//保持最新区块的哈希值
	var blockHash []byte
	//1.创建或者打开一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("create db [%s] failed %v\n",dbName,err)
	}
	//defer db.Close()

	//2.创建桶，把生成的创世区块存入数据库中
	db.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if b==nil{
			//没找到桶
			b, err := tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panicf("create bucket: [%s] failed %v\n", blockTableName,err)
			}
			//生成一个coinbase交易
			txCoinbase:=NewCoinbaseTransaction(address)
			//生成创世区块
			genesisBlock:=CreateGenesisBlock([]*Transaction{txCoinbase})
			//存储
			//1.key, value 分别以什么数据代表--hash
			//2.如何把block结构存入到数据库中--序列化

			err=b.Put(genesisBlock.Hash,genesisBlock.Serialize())//写入key 和value
			if nil !=err{
				log.Panicf("insert the genesis block failed %v\n",err)
			}
			blockHash=genesisBlock.Hash
			//存储最新的区块的哈希
			//1：lastet
			err=b.Put([]byte("1"),genesisBlock.Hash)
			if nil!=err{
				log.Panicf("save the hash of genesis block failed %v\n",err)
			}
		}
		return nil
	})
		return &BlockChain{DB:db,Tip:blockHash}
}


//添加区块到区块链中
func (bc*BlockChain) AddBlock(txs []*Transaction){
	//更新区块数据（insert)
	err:=bc.DB.Update(func(tx *bolt.Tx) error{
		//1.获取数据库桶
		b:=tx.Bucket([]byte(blockTableName))
		if nil !=b{
			//2.得到最后插入区块的序列化数据
			blockBytes:=b.Get(bc.Tip)
			//3.反序列化区块数据
			lastest_block:=DeSerializeBlock(blockBytes)
			//3.新建区块 （当前区块高度，上个区块的哈希，当前区块的数据）
			newBlock:=NewBlock(lastest_block.Heigth+1,lastest_block.Hash,txs)
			//4.存入数据库
			err:=b.Put(newBlock.Hash,newBlock.Serialize())
			if nil!=err{
				log.Panicf("insert the new block to db failed %v",err)
			}
			//更新最新区块的哈希（数据库）
			err=b.Put([]byte("1"),newBlock.Hash)
			if nil!=err{
				log.Panicf("update the latest block hash to db failed %v",err)
			}
			//更新区块链对象中的最新哈希
			bc.Tip=newBlock.Hash
		}
		return nil
	})
	if err!=nil{
		log.Panicf("insert block to db failed %v",err)
	}
}

//创建迭代器对象
func(blc *BlockChain)Iterator() *BlockChainIterator{
	return &BlockChainIterator{blc.DB,blc.Tip}
}

//遍历数据库，输出所有区块信息，区块链本身的方法
func (bc *BlockChain) PrintChain(){
	//读取数据库
	fmt.Println("区块链完整信息")
	var curBlock *Block
	bcit:=bc.Iterator()//获取迭代器对象
	//循环读取
	for{
		fmt.Println("--------------------------")
		curBlock=bcit.Next()
		//输出区块详情
		fmt.Printf("\tHash: %x\n",curBlock.Hash)
		fmt.Printf("\tPrevBlockHash: %x\n",curBlock.PrevBlockHash)
		fmt.Printf("\tTimeStamp: %s\n",time.Unix(curBlock.TimeStamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("\tHeigth: %d\n",curBlock.Heigth)
		fmt.Printf("\tNonce: %d\n",curBlock.Nonce)
		fmt.Printf("\tTxs: %v\n",curBlock.Txs)
		for _, tx:= range curBlock.Txs{
			fmt.Printf("\t\ttx-hash: %x\n",tx.TxHash)
			fmt.Printf("\t\t输入...\n")
			for _, vin:= range tx.Vins{
				fmt.Printf("\t\t\tvin-txHash : %x\n",vin.TxHash)
				fmt.Printf("\t\t\tprevious vout index: %x\n",vin.Vout)
				fmt.Printf("\t\t\tvin-scriptSig : %v\n",vin.ScriptSig)
			}
			fmt.Printf("\t\t输出...\n")
			for _, vout:=range tx.Vouts {
				fmt.Printf("\t\t\tout-value:%d\n",vout.Value)
				fmt.Printf("\t\t\tout-scriptPubkey:%v\n",vout.ScriptPubkey)

			}
		}

		//退出条件
		//转换为big.int
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt)==0{//0 表示比较的两者相等， 创世区块是nil,bigInt就是0
			break//遍历到创世区块,跳出循环
		}
	}
}

//获取一个blockchain对象
func BlockchainObject() *BlockChain {
	//获取DB
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("create db [%s] failed %v\n",dbName,err)
	}
	//defer db.Close()//不可以在这里关，这里把数据库实例给关里，接下来做任何操作都没用了。
	//获取Tip
	var tip []byte
	err=db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			tip=b.Get([]byte("1"))
		}
		return nil
	})
	if nil!=err{
		log.Panicf("get the blockchain object failed! %v\n",err)
	}
	return &BlockChain{db, tip}
}

//实现挖矿功能
//通过接收交易，生成区块
func(blockchain *BlockChain) MineNewBlock(from, to , amount []string){
	var block *Block
	//搁置交易生成步骤
	var txs []*Transaction
	value,_:=strconv.Atoi(amount[0])//转成int
	//生成新的交易
	tx:=NewSimpleTransaction(from[0],to[0],value)
	//最加到txs的交易列表中去
	txs=append(txs,tx)
	//从数据库中获取最新的一个区块
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if b!=nil{
			//获取最新的区块哈希值放在key为1的上面
			hash:=b.Get([]byte("1"))
			//再通过哈希获取最新区块
			blockBytes:=b.Get(hash)
			//反序列化
			block=DeSerializeBlock(blockBytes)
		}
		return nil
	})
	//通过数据库中最新的区块去生成更新的区块
	block=NewBlock(block.Heigth+1,block.Hash,txs)
	//持计划新生成的区块到数据库中
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if nil!=b{
			err:=b.Put(block.Hash, block.Serialize())
			if err!=nil{
				log.Printf("update the new block to db failed %v\n",err)
			}
			//更新区块的哈希值
			err=b.Put([]byte("1"),block.Hash)
			if err!=nil{
				log.Printf("update the lastest block hash to db failed %v\n",err)
			}
			blockchain.Tip=block.Hash
		}
		return nil
	})
}

//查找指定地址的UTXO
func(blockchain *BlockChain) UnUTXOS(address string) []*TxOutput{//整条链可能会有多个，所以要数组
	fmt.Printf("exec the UnUTXOS function\n")
	return nil
}