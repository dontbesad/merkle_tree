package merkle_tree

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type IMerkleTree interface {
	// 寻找叶子结点的路径
	FindLeafMTreePath(data []byte) [][]byte
	// 输出树的结构
	PrintMTree()
}

type MerkleNode struct {
	LeftNode    *MerkleNode
	RightNode   *MerkleNode
	ParentMTree *MerkleNode
	Data        []byte
	HashData    []byte
}

type MerkleTree struct {
	Root  *MerkleNode
	Datas [][]byte
	Leafs []*MerkleNode
}

func NewMTree(datas [][]byte, hashFn func([]byte) []byte) IMerkleTree {

	// 默认md5
	if hashFn == nil {
		hashFn = func(data []byte) []byte {
			v := md5.Sum(data)
			return v[:]
		}
	}

	var mNodeList = make([]*MerkleNode, len(datas))
	for i := 0; i < len(datas); i++ {
		mNodeList[i] = &MerkleNode{
			LeftNode:  nil,
			RightNode: nil,
			HashData:  hashFn(datas[i]),
			Data:      datas[i],
		}
	}

	return &MerkleTree{
		Root:  buildMTree(hashFn, mNodeList),
		Datas: datas,
		Leafs: mNodeList,
	}
}

func buildMTree(hashFn func([]byte) []byte, mNodeList []*MerkleNode) *MerkleNode {
	// root
	if len(mNodeList) == 1 {
		return mNodeList[0]
	}

	if len(mNodeList)%2 == 1 {
		mNodeList = append(mNodeList, mNodeList[len(mNodeList)-1])
	}

	var newMNodeList []*MerkleNode
	for i := 0; i < len(mNodeList); i++ {
		if i%2 == 0 {
			newMNodeList = append(newMNodeList, &MerkleNode{
				LeftNode:  mNodeList[i],
				RightNode: mNodeList[i+1],
				Data:      append(mNodeList[i].Data, mNodeList[i+1].Data...),
				HashData:  hashFn(append(mNodeList[i].Data, mNodeList[i+1].Data...)),
			})
		}

		mNodeList[i].ParentMTree = newMNodeList[i/2]

	}

	return buildMTree(hashFn, newMNodeList)
}

// bfs
func (m *MerkleTree) PrintMTree() {

	type QueData struct {
		node  *MerkleNode
		level int // 树的层级
	}

	var (
		queue     = []QueData{{node: m.Root, level: 0}} // 默认root[0]结点
		pathStr   = "level[0]"
		currLevel = 0
	)
	for len(queue) > 0 {
		val := queue[0]

		if currLevel < val.level {
			pathStr = pathStr + "\n" + fmt.Sprintf("level[%d]", val.level)
			currLevel = val.level
		}
		pathStr = pathStr + " " + hex.EncodeToString(val.node.HashData)
		// 原始数据
		if val.node.LeftNode == nil && val.node.RightNode == nil {
			pathStr = pathStr + " data(" + string(val.node.Data) + ")"
		}

		if val.node.LeftNode != nil {
			queue = append(queue, QueData{node: val.node.LeftNode, level: val.level + 1})
		}
		if val.node.RightNode != nil {
			queue = append(queue, QueData{node: val.node.RightNode, level: val.level + 1})
		}
		queue = queue[1:len(queue)]
	}

	fmt.Println(pathStr)
}

func (m *MerkleTree) FindLeafMTreePath(data []byte) (path [][]byte) {
	for i := 0; i < len(m.Datas); i++ {
		if hex.EncodeToString(m.Datas[i]) == hex.EncodeToString(data) {
			return m.searchParentMNode(m.Leafs[i])
		}
	}
	return path
}

func (m *MerkleTree) searchParentMNode(node *MerkleNode) (path [][]byte) {
	if node.ParentMTree == nil {
		return [][]byte{node.HashData}
	}
	// fmt.Println("->", hex.EncodeToString(node.HashData))
	return append(m.searchParentMNode(node.ParentMTree), node.HashData)
}
