package BLC

import "bytes"

//交易输入管理

//输入结构，输入的交易原来有多少钱，这里也有多少钱，所以不需要金额
type TxInput struct{
	TxHash		[]byte	//交易哈希，这笔钱从哪个区块来
	Vout	int	//引用上一笔交易的输出(Output)中的索引号，一个区块里会有好多交易
	Signature   []byte//数字签名
	PublicKey   []byte //公钥
}


//验证应用的地址是否匹配
//func (txIntput *TxInput) CheckPubkeyWithAddress(address string)bool{
//	return address==txIntput.ScriptSig
//}

//传递哈希160进行判断
func (in *TxInput) UnLockRipemd160Hash(ripedmd160Hash []byte)bool{
	inputRipemd160Hash:=Ripemd160Hash(in.PublicKey)//获取input的ripemd160哈希值
	return bytes.Compare(inputRipemd160Hash,ripedmd160Hash)==0//两者相等
}