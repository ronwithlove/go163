package BLC

//输出交易管理

//输出结构，给谁，给多少钱（如果没有刚好用完，就会找零给自己，就会有2个以上的output）
type TxOutput struct{
	value 	int//金额 收钱的金额
	ScriptPubkey 	string//用户名(UTXO 的所有者) 收钱的人
}