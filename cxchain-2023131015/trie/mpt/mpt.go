package mpt

import (
	"bytes"
	"cxchain-2023131015/kvstore/leveldb"
	"fmt"
)

// MPT 结构体表示一个Merkle Patricia Trie
// Root: 树的根节点
// DB: 底层数据库存储
type MPT struct {
	Root Node
	DB   *leveldb.LevelDBStore
}

// NewMPT 创建一个新的MPT实例
func NewMPT(db *leveldb.LevelDBStore) *MPT {
	return &MPT{
		Root: nil,
		DB:   db,
	}
}

// Loadnode 从数据库加载节点
func (mpt *MPT) Loadnode(key []byte) (Node, error) {
	fmt.Printf("加载节点\n")
	data, err := mpt.DB.Get(key[:])
	if err != nil {
		return nil, err
	}
	var node Node
	nodereturn, err:=DeserializeNode(data, &node)
	if err != nil {
		fmt.Println("loadnode方法：反序列化节点失败")
		return nil, err
	}
	
	fmt.Printf("加载节点成功: key=%v, node=%v\n", key, nodereturn)
	return nodereturn, nil
}

// Savenode 将节点保存到数据库
func (mpt *MPT) Savenode(node Node) error {
	code, err := node.Serialize()
	if err != nil {
		return err
	}
	h := Hash(code)
	mpt.DB.Put(h[:], code)
	fmt.Printf("\n保存节点成功: key=%v, node=%v\n", h, node)
	return nil
}

// findSamePrefix 查找两个字节切片的共同前缀
func findSamePrefix(a, b []byte) []byte {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[:i]
}

// Insert 向MPT中插入键值对
func (m *MPT) Insert(key, value []byte) error {
	nibbles := ToNibbles(key)
	fmt.Printf("插入: key=%v, value=%v, nibbles=%x\n", key, value, nibbles)
	if m.Root == nil {
		m.Root = &LeafNode{
			NodeType: LeafNodeType,
			Key:      nibbles,
			Value:    value,
		}
		fmt.Printf("创建root节点: key=%v, value=%v\n", nibbles, value)
		return m.Savenode(m.Root)
	}
	
	newRoot, err := m.insert(m.Root, nibbles, value)
	if err != nil {
		return err
	}
	m.Root = newRoot
	return m.Savenode(m.Root)
}

// Has 检查MPT中是否存在指定的key
func (m *MPT) Has(key []byte) ([]byte, error) {
	nibbles := ToNibbles(key)
	if m.Root == nil {
		return nil, fmt.Errorf("empty trie")
	}
	value, err := m.has(m.Root, nibbles)

	if err != nil {
		return nil, fmt.Errorf("key not found8")
	}
	fmt.Println("查找成功，value为",value)
	return value, nil
}

// insert 内部方法，递归插入键值对
func (mpt *MPT) insert(node Node, nibbles, value []byte) (Node, error) {

	switch n := node.(type) {
	case *LeafNode:
		Prefix := findSamePrefix(n.Key, nibbles)
		//Prefix输出
		fmt.Println("共同前缀=", Prefix)
		if bytes.Equal(n.Key, nibbles) {
			n.Value = value
			if err := mpt.Savenode(n); err != nil {
				return nil, err
			}
			return n, nil
		}

		branch := &BranchNode{
			NodeType: BranchNodeType,
			Child:    [17][]byte{},
		}

		if len(n.Key) > len(Prefix) {
			keyleft := n.Key[len(Prefix)]
			keyright := n.Key[len(Prefix)+1:]
			//keyleft与keyright输出
			fmt.Printf("keyleft=%v, keyright=%v\n", keyleft, keyright)
			NewLeafNode := &LeafNode{
				NodeType: LeafNodeType,
				Key:      keyright,
				Value:    n.Value,
			}
			err := mpt.Savenode(NewLeafNode)
			if err != nil {
				return nil, err
			}
			branch.Child[keyleft] = NewLeafNode.GetHash()
			fmt.Printf("Create branch LeafNode: key=%v, value=%v\n", keyright, n.Value)
		}

		if len(nibbles) > len(Prefix) {
			nibblesleft := nibbles[len(Prefix)]
			nibblesright := nibbles[len(Prefix)+1:]
			//nibblesleft与nibblesright输出
			fmt.Printf("nibblesleft=%v, nibblesright=%v\n", nibblesleft, nibblesright)
			NewLeafNode := &LeafNode{
				NodeType: LeafNodeType,
				Key:      nibblesright,
				Value:    value,
			}
			err := mpt.Savenode(NewLeafNode)
			if err != nil {
				return nil, err
			}
			branch.Child[nibblesleft] = NewLeafNode.GetHash()
			fmt.Printf("Create branch LeafNode: key=%v, value=%v\n", nibblesright, n.Value)
		}

		if len(Prefix) > 0 {
			err := mpt.Savenode(branch)
			if err != nil {
				return nil, err
			}
			ex := &ExtensionNode{
				NodeType: ExtensionNodeType,
				Key:      Prefix,
				Value:    branch.GetHash(),
			}
			if err := mpt.Savenode(ex); err != nil {
				return nil, err
			}
			return ex, nil
		}

		err := mpt.Savenode(branch)
		if err != nil {
			return nil, err
		}
		//branch输出
		fmt.Printf("Create branch: %v", branch)
		return branch, nil
	case *ExtensionNode:
		Prefix := findSamePrefix(n.Key, nibbles)
		// Create a branch node to split at the diverging point
		branch := &BranchNode{
			NodeType: BranchNodeType,
			Child:    [17][]byte{},
		}
		if len(Prefix) != len(n.Key) {
			if len(n.Key[len(Prefix):]) > 0 {
				idx := n.Key[len(Prefix)]
				ex := &ExtensionNode{
					NodeType: ExtensionNodeType,
					Key:      n.Key[len(Prefix)+1:],
					Value:    n.Value,
				}
				if err := mpt.Savenode(ex); err != nil {
					return nil, err
				}
				branch.Child[idx] = ex.GetHash()
			}

			if len(nibbles[len(Prefix):]) > 0 {
				idx := nibbles[len(Prefix)]
				LeafNode := &LeafNode{
					NodeType: LeafNodeType,
					Key:      nibbles[len(Prefix)+1:],
					Value:    value,
				}
				if err := mpt.Savenode(LeafNode); err != nil {
					return nil, err
				}
				branch.Child[idx] = LeafNode.GetHash()
			}

			if len(Prefix) > 0 {
				if err := mpt.Savenode(branch); err != nil {
					return nil, err
				}
				ext := &ExtensionNode{
					NodeType: ExtensionNodeType,
					Key:      Prefix,
					Value:    branch.GetHash(),
				}
				if err := mpt.Savenode(ext); err != nil {
					return nil, err
				}
				return ext, nil
			}

			if err := mpt.Savenode(branch); err != nil {
				return nil, err
			}
			return branch, nil
		}

		Lnode, err := mpt.Loadnode(n.Value)
		if err != nil {
			return nil, err
		}
		newChild, err := mpt.insert(Lnode, nibbles[len(Prefix):], value)
		if err != nil {
			return nil, err
		}
		n.Value = newChild.GetHash()
		if err := mpt.Savenode(n); err != nil {
			return nil, err
		}
		return n, nil

	case *BranchNode:
		fmt.Println("进入分支节点处理")
		if len(nibbles) == 0 {
			fmt.Println("分支节点：nibbles为空，设置值到Child[16]")
			var hash []byte
			hash = Hash(value)
			n.Child[16] = hash
			if err := mpt.Savenode(n); err != nil {
				return nil, err
			}
			return n, nil
		}

		idx := nibbles[0]
		fmt.Printf("分支节点：处理索引 %d，剩余nibbles长度 %d\n", idx, len(nibbles)-1)
		child, err := mpt.Loadnode(n.Child[idx])
		if err != nil {
			fmt.Println("分支节点：加载子节点失败")
			return nil, err
		}
		if child == nil {
			fmt.Println("分支节点：子节点为空，创建新叶子节点")
			// Create a new leaf node
			leaf := &LeafNode{
				NodeType: LeafNodeType,
				Key:      nibbles[1:],
				Value:    value,
			}
			if err := mpt.Savenode(leaf); err != nil {
				return nil, err
			}
			n.Child[idx] = leaf.GetHash()
		} else {
			fmt.Println("分支节点：子节点存在，递归插入")
			// Insert into existing child
			newChild, err := mpt.insert(child, nibbles[1:], value)
			if err != nil {
				return nil, err
			}
			n.Child[idx] = newChild.GetHash()
		}
		if err := mpt.Savenode(n); err != nil {
			return nil, err
		}
		return n, nil
	default:
		return nil, fmt.Errorf("unknown")
	}
}

// has 内部方法，递归查找key对应的value
func (mpt *MPT) has(node Node, nibbles []byte) ([]byte, error) {
    fmt.Println("has方法:path:", nibbles)
    switch n := node.(type) {
    case *LeafNode:
        fmt.Println("has方法：进入叶子节点")
        fmt.Printf("has方法：叶子节点key长度 %d, 查找key长度 %d\n", len(n.Key), len(nibbles))
        
        if len(n.Key) == 0 {
            fmt.Println("has方法：叶子节点key为空，直接返回value")
            return n.Value, nil
        }
        
        if len(nibbles) < len(n.Key) {
            fmt.Println("has方法：查找key长度不足，key不存在")
            return nil, fmt.Errorf("key not found1")
        }
        
        if !bytes.Equal(n.Key, nibbles[:len(n.Key)]) {
            fmt.Println("has方法：key前缀不匹配，key不存在")
            return nil, fmt.Errorf("key not found2")
        }
        
        if len(nibbles) != len(n.Key) {
            fmt.Println("has方法：key长度不匹配，key不存在")
            return nil, fmt.Errorf("key not found3")
        }
        
        fmt.Println("has方法：找到匹配的叶子节点")
        return n.Value, nil
        
    case *ExtensionNode:
        fmt.Println("has方法：进入扩展节点")
        fmt.Printf("has方法：扩展节点路径长度 %d, 查找key长度 %d\n", len(n.Key), len(nibbles))
        
        if len(nibbles) < len(n.Key) {
            fmt.Println("has方法：查找key长度不足，key不存在")
            return nil, fmt.Errorf("key not found4")
        }
        
        if !bytes.Equal(n.Key, nibbles[:len(n.Key)]) {
            fmt.Println("has方法：路径前缀不匹配，key不存在")
            return nil, fmt.Errorf("key not found5")
        }
        
        fmt.Println("has方法：路径匹配，加载子节点继续查找")
        child, err := mpt.Loadnode(n.Value)
        if err != nil {
            fmt.Println("has方法：加载子节点失败")
            return nil, err
        }
        return mpt.has(child, nibbles[len(n.Key):])
        
    case *BranchNode:
        fmt.Println("has方法：进入分支节点")
        if len(nibbles) == 0 {
            fmt.Println("has方法：nibbles为空，返回Child[16]的值")
            return n.Child[16], nil
        }
        idx := nibbles[0]
        fmt.Printf("has方法：处理索引 %d，剩余nibbles长度 %d\n", idx, len(nibbles)-1)
        if idx >= 16 {
            fmt.Printf("has方法：无效的nibble值 %d\n", idx)
            return nil, fmt.Errorf("invalid nibble value: %d", idx)
        }
        child, err := mpt.Loadnode(n.Child[idx])
        fmt.Println(n.Child[idx])
        if err != nil {
            fmt.Println("has方法：加载子节点失败")
            return nil, fmt.Errorf("key not found6")
        }
        if child == nil {
            fmt.Println("has方法：子节点为空")
            return nil, fmt.Errorf("key not found7")
        }
        fmt.Println("has方法：递归查找子节点")
        return mpt.has(child, nibbles[1:])

    default:
        return nil, fmt.Errorf("unknown node type")
    }
}
