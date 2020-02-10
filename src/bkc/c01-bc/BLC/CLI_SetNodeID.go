package BLC

import (
	"fmt"
	"log"
	"os"
)


//设置端口号（环境变量）
func (cli *CLI)SetNodeId(nodeId string)  {
	if nodeId==""{
		fmt.Println("Please set the node id...")
		os.Exit(1)
	}
	err:=os.Setenv("NODE_ID",nodeId)
	if nil!=err{
		log.Fatalf("set env failed! %v\n",err)
	}
}
