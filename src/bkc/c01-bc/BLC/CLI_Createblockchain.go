package BLC

//初始化区块链
func(cli *CLI) createBlockchain(address string){
	bc:=CreateBlockCHainWithGenesisBlock(address)
	defer bc.DB.Close()

	//设置utxo重制操作
	utxoSet:=&UTXOSet{bc}
	utxoSet.ResetUTXOSet()


}
