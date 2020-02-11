package BLC

//初始化区块链
func(cli *CLI) createBlockchain(address string,nodeID string){
	bc:=CreateBlockCHainWithGenesisBlock(address,nodeID)
	defer bc.DB.Close()

	//设置utxo重制操作
	utxoSet:=&UTXOSet{bc}
	utxoSet.ResetUTXOSet()


}
