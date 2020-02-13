package BLC

import "bytes"

//输出交易管理

//输出结构，给谁，给多少钱（如果没有刚好用完，就会找零给自己，就会有2个以上的output）
type TxOutput struct{
	Value 	int//金额 收钱的金额
<<<<<<< HEAD
	//ScriptPubkey 	string
	Ripemd160Hash   []byte	//用户名，准备花钱的人

}

//通过检查出入地址和key是否相等来，验证当前UTXO是否属于指定的地址
//func (txOutput *TxOutput) CheckPubkeyWithAddress(address string)bool{
//	return address==txOutput.ScriptPubkey
//}

//output身份验证
func (TxOutput *TxOutput) UnLockScriptPubkeyWithAddress(address string)bool{
	hash160:=StringToHash160(address)
	return bytes.Compare(hash160,TxOutput.Ripemd160Hash)==0
}

//新建output对象
func NewTxOutput(value int, address string) *TxOutput{
	return &TxOutput{value,StringToHash160(address)}
}
=======
	ScriptPubkey 	string//用户名(UTXO 的所有者) 收钱的人
}

//通过检查出入地址和key是否相等来，验证当前UTXO是否属于指定的地址
func (txOutput *TxOutput) CheckPubkeyWithAddress(address string)bool{
	return address==txOutput.ScriptPubkey
}


>>>>>>> Master
