package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//请求处理文件管理

//version
func handleVersion(request []byte, bc *BlockChain){
	fmt.Println("the request of version handle...")
	var buff bytes.Buffer
	var data Version
	//1.解析请求
	dataBytes:=request[12:]//把request去了，拿到后面的msg
	//2.生成version结构
	decoder:=gob.NewDecoder(&buff)
	buff.Write(dataBytes)
	if err:=decoder.Decode(&data);nil!=err{
		log.Panicf("decode the version struct faile! %v\n",err)
	}
	//3.获取请求方的区块高度
	versionHeight:=data.Height
	//4.获取自身节点的区块高度
	height:=bc.GetHeight()
	//如果当前节点的区块高度大于versionHeight
	//将当前节点版本信息发送给请求节点
	if height> int64(versionHeight){//收到请求的高度大
		sendVersion(data.AddrFrom,bc)//把version发给对方
	}else if height<int64(versionHeight){
		//当前节点区块高度小于发request方，
		sendGetBlocks(data.AddrFrom)//向发送方发起同步数据请求
	}
}

//GetBlocks
func handleGetBlocks(request []byte, bc *BlockChain){

}

//Inv
func handleInv(request []byte, bc *BlockChain){

}

//GetData
func handleGetData(request []byte, bc *BlockChain){

}

//Block
func handleBlocks(request []byte, bc *BlockChain){

}
