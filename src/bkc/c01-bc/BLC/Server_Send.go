package BLC

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

//请求发送文件

//发送请求
func sendMessage(to string, msg []byte){
	fmt.Println("向服务器发送请求...")
	//1.连接上服务器
	conn,err:=net.Dial(PROTOCOL,to)
	if nil!=err{
		log.Panicf("connect to server [%s] failed! %v\n",err)
	}
	//要发送的数据
	_, err=io.Copy(conn,bytes.NewReader(msg))
	if nil!=err{
		log.Panicf("add the data to conn failed! %v\n",err)
	}
}

//区块链版本验证
func sendVersion(toAddress string){
	//1.获取当前节点的区块高度
	height:=1
	//2.组装生成version
	versionData:=Version{Height:height,AddrFrom:nodeAddress}
	//3.数据系列化
	data:=gobEncode(versionData)
	//4.将命令与版本组装成完整的请求
	request:=append(commandToBytes(CMD_VERSION),data...)
	//4.发送请求
	sendMessage(toAddress,request)
}

//从指定节点同步数据
func sendGetBlocks(){

}

//发送获取指定区块请求
func sendGetData(){

}

//向其他节点展示
func sendInv(){

}

//发送区块信息
func sendBlock(){

}