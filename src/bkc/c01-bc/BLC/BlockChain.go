package BLC
//区块链管理文件
type BlockChain struct{//直接用切片也可以，但是结构体比较正式一点
	Blocks []*Block //区块的切片
}


//初始化区块链
func CreateBlockCHainWithGenesisBlock() *BlockChain{
	block:=CreateGenesisBlock([]byte("init blockchain"))
	return &BlockChain{[]*Block{block}}//把第一个区块添加到区块链中去了
}


//添加区块到区块链中
func (bc*BlockChain) AddBlock(height int64,  prevBlockHash []byte, data []byte){
	newBlock:=NewBlock(height,prevBlockHash,data)
	bc.Blocks=append(bc.Blocks,newBlock)//往BlockChain结构体中的Blocks切片加区块
}
