# merkle tree

## API

```go
type IMerkleTree interface {
	// 寻找叶子结点的路径，返回从根到叶子的hash路径
	FindLeafMTreePath(data []byte) [][]byte
	// 输出树的结构
	PrintMTree()
}
```

## Usage

```go
datas := [][]byte{[]byte("111"), []byte("222")}

mTree := NewMTree(datas, nil) // default md5
leafPath := mTree.FindLeafMTreePath([]byte("555"))
...
```