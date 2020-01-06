package BLC

//输出交易管理

//输出结构
type TxOutput struct{
	value 	int//金额 收钱的金额
	ScriptPubkey 	string//用户名(UTXO 的所有者) 收钱的人
}