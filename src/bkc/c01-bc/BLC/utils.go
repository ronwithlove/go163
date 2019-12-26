package BLC

import (
	"bytes"
	"encoding/binary"
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