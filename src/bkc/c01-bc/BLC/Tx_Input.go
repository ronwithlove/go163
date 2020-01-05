package BLC

//交易输入管理

//输入结构
type TxInput struct{
	//交易哈希，非当前交易哈希
	TxHash		[]byte
	//引用上一笔交易的输出索引号
	Vout	int
	//用户名
	ScriptSig	string
}
