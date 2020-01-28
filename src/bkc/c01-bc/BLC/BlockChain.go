package BLC

import (
	"encoding/hex"
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
	tx:=NewSimpleTransaction(from[0],to[0],value,blockchain)
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

//获取指定地址已花费输出,虽然叫已花费输出，就是区块交易中的Input,不是Output不要搞混
func (blockchain *BlockChain) SpentOutpts(address string) map[string][]int {
	//已花费输出缓存
	spentTXOutputs:=make(map[string][]int)
	bcit:=blockchain.Iterator()
	for{
		block:=bcit.Next()
		for _, tx:= range block.Txs{//一个block里会有多个交易
			//排除coinbase交易
			if!tx.IsCoinbaseTransaction(){
				for _, in:=range tx.Vins{//一个交易里会有多个input
					if in.CheckPubkeyWithAddress(address){}
					key:=hex.EncodeToString(in.TxHash)//交易哈希转成string保存，作为key
					//添加到已花费输出的缓存中
					spentTXOutputs[key]= append(spentTXOutputs[key], in.Vout)
					//在一个Input中某个人可能有多条记录in.Vout(index,引用上一笔交易的输出索引号)
					//保存在以哈希为key,value是int的数组中
				}
			}
		}
		//退出循环条件,直到创世区块
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}
	return spentTXOutputs
}

//查找指定地址的UTXO
/*
	遍历查找区块链数据库中的每一个区块中的每一个交易
	查找每一个交易中的每一个输出
	判断每个输出是否满足下列条件
	1.属于传入的地址(就是谁，谁花了钱)
	2.是否未被花费
		1.先遍历一次区块链数据库，将所有自己花费的OUTPUT存入一个缓存
		2.再次遍历区块链数据库，检查每一个VOUT是否包含在前面的已花费的缓存中
 */
func(blockchain *BlockChain) UnUTXOS(address string,txs []*Transaction) []*UTXO{//整条链可能会有多个，所以要数组
	//1.遍历数据库，查找所有与address相关的交易
	//获取迭代器
	bcit:=blockchain.Iterator()
	//当前地址的未花费输出列表
	var unUTXOS []*UTXO
	//获取指定地址所有已花费输出，得到改地址的所有的input
	spentTXOutputs:=blockchain.SpentOutpts(address)
	//缓存迭代
	//查找缓存中的已花费输出
	for _,tx:=range txs{
		//判断coinbaseTransaction
		if!tx.IsCoinbaseTransaction(){
			for _,in:=range tx.Vins{
				//判断用户
				if in.CheckPubkeyWithAddress(address){
					//添加到已花费输出的map中
					key:=hex.EncodeToString(in.TxHash)
					spentTXOutputs[key]=append(spentTXOutputs[key],in.Vout)
				}
			}
		}
	}
	//遍历缓存中的UTXO
	for _, tx:=range txs{
		//添加一个缓存输出的跳转
		WorkCacheTx:
		for index,vout:=range tx.Vouts{
			if vout.CheckPubkeyWithAddress(address){
				if len(spentTXOutputs)!=0{
					var isUtxoTx bool //判断交易是否被其他交易引用
					for txHash, indexArray:=range spentTXOutputs{
						txHashStr:= hex.EncodeToString(tx.TxHash)
						if txHash ==txHashStr{
							//当前遍历到的交易已经有输出被其他交易的输入所引用
							isUtxoTx=true
							//添加状态变量，判断指定的output是否被引用
							var isSpentUTXO bool
							for _,voutIndex:=range indexArray{
								if index==voutIndex{
									//该输出被引用
									isSpentUTXO=true
									//跳出当前vout判断逻辑，进行下一个输出判断
									continue WorkCacheTx
								}
							}
							if isSpentUTXO==false{
								utxo:=&UTXO{tx.TxHash,index,vout}
								unUTXOS=append(unUTXOS,utxo)
							}
						}
					}
					if isUtxoTx==false{
						//说明当前交易中所有与address 相关的outputs 都是UTXO
						utxo:=&UTXO{tx.TxHash,index,vout}
						unUTXOS=append(unUTXOS,utxo)
					}
				}else{
					utxo:=&UTXO{tx.TxHash,index,vout}
					unUTXOS=append(unUTXOS,utxo)
				}
			}
		}
	}

	//优先遍历缓存中的UTXO(因为多笔交易可能刚好就在一个区块中)，如果余额足够，直接返回，如果不足，再遍历db文件中的UTXO
	//数据库迭代，不断获取下一个区块
	//迭代，不断获取下一个区块
	for{
		block:=bcit.Next()
		//遍历区块中的每笔交易
		for _, tx:= range block.Txs{//每个区块有多个交易
			//跳转
			work:
			for index,vout:=range tx.Vouts{//每个交易有多个output(tx中output是数组)
				//index：当前输出在当前交易的中索引位置
				//vout:当前输出
				if vout.CheckPubkeyWithAddress(address){
					//当前vout属于传入地址
					if len(spentTXOutputs)!=0{
						var isSpentOutput bool//默认就是false
						for txHash, indexArray :=range spentTXOutputs{//遍历key=txHash,value=indexArray
							for _, i:=range indexArray{//遍历array内保存的多条in.Vout记录
								if txHash==hex.EncodeToString(tx.TxHash)&&index==i{
									//txHash== hex.EncodeToString(tx.TxHash),
									//说明当前的交易tx至少已经有输出被其他交易的输入引用
									//index==i 说明正好是当前的输出被其他交易引用
									//跳转到最外层循环，判断下一个VOUT
									isSpentOutput=true
									continue work
								}
							}
						}
						if isSpentOutput==false{
							utxo:=&UTXO{tx.TxHash,index,vout}
							unUTXOS=append(unUTXOS,utxo)
						}
					}else{//如果长度为0，表示没有找到output信息，就表示他没花过钱
						//将当前地址所有输出都添加到未花费输出中
						utxo:=&UTXO{tx.TxHash,index,vout}
						unUTXOS=append(unUTXOS,utxo)
					}
				}
			}
		}

		//退出循环条件,直到创世区块
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}

	return unUTXOS
}


//查询余额
func (blockchain *BlockChain) getBalance (address string) int{
	var amount int
	utxos:=blockchain.UnUTXOS(address)
	for _, utxo := range utxos{
		amount+=utxo.Output.Value
	}
	return  amount
}

//查找指定地址的可用UTXO,超过amount就中断查找
//更新当前数据库中指定地址的UTXO数量
//txs:缓存中的交易列表（用于多笔交易处理）
func(blockchain * BlockChain) FindSpendableUTXO(from string, amount int,txs[]*Transaction)(int, map[string][]int){
	spendableUTXO:= make(map[string][]int)

	var value int
	utxos:=blockchain.UnUTXOS(from,txs)
	//遍历UTXO
	for _, utxo := range utxos{
		value += utxo.Output.Value
		//计算交易哈希
		hash:=hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash]=append(spendableUTXO[hash],utxo.Index)
		if value>=amount{
			break
		}
	}

	//遍历完所有交易后，依然小于amount
	//资金不足
	if value<amount{
		fmt.Printf("地址[%s]余额不足，当前余额[%d],转账金额[%d]\n",from,value,amount)
		os.Exit(1)
	}
	return value, spendableUTXO
}