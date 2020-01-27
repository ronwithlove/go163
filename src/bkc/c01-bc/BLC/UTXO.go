package BLC

//UTXO结构管理
type UTXO struct{
	//来自交易的哈希
	TxHash []byte
	//在这个哈希下的输出列表中的索引
	Index  int
	//未花费的交易输出
	Output *TxOutput
}
