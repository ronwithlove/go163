package BLC

import "crypto/sha256"

//Merkle树实现管理
type MerkleTree struct{
	//根节点最上面那个，他上面和左右两边都没有节点了
	RootNode *MerkleNode
}

//merkle节点结构
type MerkleNode struct{
	//左子节点
	Left *MerkleNode
	//右子节点
	Right *MerkleNode
	//数据（哈希）
	Data []byte

}

//创建Merkle树

//创建Merkle节点
func MakeMerkleNode(left,right *MerkleNode,data[]byte) *MerkleNode{
	node:=&MerkleNode{}
	//判断叶子节点
	if left==nil && right==nil{//最下面一层，他的下面没有子节点，HA
		hash:=sha256.Sum256(data)
		node.Data=hash[:]
	}else{//非叶子节点 HAB,HCD
		prveData:=append(left.Data, right.Data...)
		hash:=sha256.Sum256(prveData)
		node.Data=hash[:]
	}
	//子节点赋值
	node.Left=left
	node.Right=right
	return  node
}
