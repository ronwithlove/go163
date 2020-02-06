package BLC

//重制 utxo table
func (cli *CLI) TestResetUTXO(){
	blockchain:=BlockchainObject()
	defer blockchain.DB.Close()
	utxoSet:=UTXOSet{Blockchain:blockchain}
	utxoSet.ResetUTXOSet()
}