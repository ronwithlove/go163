package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//参数数量检测函数
func IsValidArgs(){
	if len(os.Args)<2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}

//实现int64转[]byte,函数
func IntToHex(data int64)[]byte{
	buffer:=new(bytes.Buffer)
	err:=binary.Write(buffer,binary.BigEndian,data)//进行转换,将data的binary编码格式写入buffer
	if err!=nil{
		log.Panicf("int to []byte failed! %v\n",err)
	}
	return buffer.Bytes()
}

//标准JSON格式转切片
//Mac Terminal格式：
//./bc send -from '["ron","amy"]' -to '["aaron","norton"]' -amount '["20","10"]'
//./bc send -from '["ron","amy"]' -to '["amy","ron"]' -amount '["5","2"]'
// ./bc send -from '["ron"]' -to '["aaron"]' -amount '["3"]'
//Windows格式：
// bc.exe send -from "[\"ron\"]" to "[\"aaron\"]" -amount "[\"100\"]"
func JSONToSlice(jsonString string)[]string{
	var strSlice []string
	//json
	if err:=json.Unmarshal([]byte(jsonString),&strSlice);nil!=err{
		log.Printf("json to []string failed! %v\n",err)
	}
	return  strSlice
}

//string to hash160
func StringToHash160(address string)[]byte{
	pubKeyHash:=Base58Decode([]byte(address))
	hash160:=pubKeyHash[:len(pubKeyHash)-addressCheckSumLen]
	return hash160
}

//获取节点ID
func GetEnvNodeId()string{
	nodeID:=os.Getenv("NODE_ID")

	if nodeID==""{
		fmt.Println("NODE_ID is not set...")
		os.Exit(1)
	}
	fmt.Printf("NODE_ID:%s\n",nodeID)
	return nodeID
}