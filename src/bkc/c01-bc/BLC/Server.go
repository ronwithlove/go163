package BLC

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

//网络服务文件管理

//3000作为引导节点（主节点）的地址
var knownNodes =[]string{"localhost:3000"}

//当前区块版本信息（决定区块是否需要同步）
type Version struct{
//	Version int 	//版本号
	Height int		//当前节点区块高度
	AddrFrom string //当前节点地址
}

//节点服务
var nodeAddress string//全球变量

//启动服务
func startServer(nodeID string){
	fmt.Printf("启动节点[%s]...\n",nodeID)
	//节点地址赋值
	nodeAddress=fmt.Sprintf("localhost:%s",nodeID)
	//1.监听节点
	listen,err:=net.Listen(PROTOCOL,nodeAddress)
	if nil!=err{
		log.Panicf("listen address of %s failed! %v\n",nodeAddress,err)
	}
	defer listen.Close()//别忘记关闭
	//两个节点，主节点负责保存数据，钱包节点负责发送请求，同步数据
	if nodeAddress!=knownNodes[0]{//不是主节点的时候，发送请求，同步数据
	//...
		//sendMessage(knownNodes[0],nodeAddress)
		sendVersion(knownNodes[0])
	}
	for{
		//2.生成连接，接收请求
		conn,err:=listen.Accept()
		if nil!=err{
			log.Panicf("accept connect failed! %v\n",err)
		}
		request,err:=ioutil.ReadAll(conn)//得到请求
		if nil!=err{
			log.Panicf("Receive Message failed! %v\n",err)
		}

		//3.处理请求
		fmt.Printf("Receive a Message:%v\n",request)
		handleConnection()
	}
}

//请求处理函数
func handleConnection(){

}


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
	request:=append(commandToBytes(VERSION),data...)
	//4.发送请求
	sendMessage(toAddress,request)


}