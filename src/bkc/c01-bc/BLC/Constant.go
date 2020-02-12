package BLC

//网络服务常量管理
//协议
const PROTOCOL="tcp"

//命令长度
const CMMAND_LENGTH=12

//命令分类
const(
	CMD_VERSION="version"//验证当前节点末端区块是否是最新区块
	CMD_GETBLOCKS="getbloks"//从最长链上获取区块
	CMD_INV="inv"//向其他节点展示当前节点有哪些区块
	CMD_GETDATA="getdata"//请求指定区块
	CMD_BLOCK="block"//接收到新的区块之后，进行处理
 	)