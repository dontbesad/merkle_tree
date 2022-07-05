package merkle_tree

import (
	"encoding/hex"
	"testing"
)

func TestMerkleTree(t *testing.T) {

	var datas = [][]byte{[]byte("111"), []byte("222"), []byte("333"), []byte("444"), []byte("555")}

	mTree := NewMTree(datas, nil)
	mTree.PrintMTree()

	leafPath := mTree.FindLeafMTreePath([]byte("555"))
	for i, v := range leafPath {
		t.Logf("level[%d] %+v", i, hex.EncodeToString(v))
	}
}
