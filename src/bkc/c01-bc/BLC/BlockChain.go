package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
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
func CreateBlockCHainWithGenesisBlock(txs []*Transaction) *BlockChain{
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
			//生成创世区块
			genesisBlock:=CreateGenesisBlock(txs)
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
		fmt.Printf("\tTxs: %v\n",curBlock.Txs)
		fmt.Printf("\tHeigth: %d\n",curBlock.Heigth)
		fmt.Printf("\tNonce: %d\n",curBlock.Nonce)
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