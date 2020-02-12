package BLC

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//网络服务文件管理

//3000作为引导节点（主节点）的地址
var knownNodes =[]string{"localhost:3000"}

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
		//处理请求
		//单独启动一个goroutine 进行请求处处理
		handleConnection(conn)
	}
}


//请求处理函数
func handleConnection(conn net.Conn){
	request,err:=ioutil.ReadAll(conn)//得到请求
	if nil!=err{
		log.Panicf("Receive a request failed! %v\n",err)
	}
	cmd:=bytesToCommand(request)//接收到的命令反序列话
	fmt.Printf("Receive a Command: %s\n",cmd)//打出来看下
	switch cmd{
	case CMD_VERSION:
		handleVersion()
	case CMD_GETDATA:
		handleGetData()
	case CMD_GETBLOCKS:
		handleGetBlocks()
	case CMD_INV:
		handleInv()
	case CMD_BLOCK:
		handleBlocks()
	default:
		fmt.Println("Unknow command")
	}
}


