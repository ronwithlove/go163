package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

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