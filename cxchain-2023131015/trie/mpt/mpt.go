package mpt

import (
	"bytes"
	"cxchain-2023131015/kvstore/leveldb"
	"fmt"
)

type MPT struct {
	Root Node
	DB   *leveldb.LevelDBStore
}

func NewMPT(db *leveldb.LevelDBStore) *MPT {
	return &MPT{
		Root: nil,
		DB:   db,
	}
}

func (mpt *MPT) Loadnode(key []byte) (Node, error) {
	data, err := mpt.DB.Get(key[:])
	if err != nil {
		return nil, err
	}
	var node Node
	deserializeNode(data, node)
	return node, nil
}

func (mpt *MPT) Savenode(node Node) error {
	code, err := serializeNode(node)
	if err != nil {
		return err
	}
	h := Hash(code)
	mpt.DB.Put(h[:], code)
	return nil
}

func findSamePrefix(a, b []byte) []byte {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[:i]
}

func (m *MPT) Insert(key, value []byte) error {
	nibbles := ToNibbles(key)
	if m.Root == nil {
		m.Root = &LeafNode{
			NodeType: LeafNodeType,
			Key:      nibbles,
			Value:    value,
		}
		return m.Savenode(m.Root)
	}
	
	newRoot, err := m.insert(m.Root, nibbles, value)
	if err != nil {
		return err
	}
	m.Root = newRoot
	return m.Savenode(m.Root)
}

func (m *MPT) Has(key []byte) ([]byte, error) {
	nibbles := ToNibbles(key)
	if m.Root == nil {
		return nil, fmt.Errorf("empty trie")
	}
	
	value, err := m.has(m.Root, nibbles)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (mpt *MPT) insert(node Node, nibbles, value []byte) (Node, error) {

	switch n := node.(type) {
	case *LeafNode:
		Prefix := findSamePrefix(n.Key, nibbles)

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
		}

		if len(nibbles) > len(Prefix) {
			nibblesleft := nibbles[len(Prefix)]
			nibblesright := nibbles[len(Prefix)+1:]
			NewLeafNode := &LeafNode{
				NodeType: LeafNodeType,
				Key:      nibblesright,
				Value:    n.Value,
			}
			err := mpt.Savenode(NewLeafNode)
			if err != nil {
				return nil, err
			}
			branch.Child[nibblesleft] = branch.GetHash()
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
				ex := &ExtensionNode{
					NodeType: ExtensionNodeType,
					Key:      nibbles[len(Prefix)+1:],
					Value:    n.Value,
				}
				if err := mpt.Savenode(ex); err != nil {
					return nil, err
				}
				branch.Child[idx] = ex.GetHash()
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
		if len(nibbles) == 0 {
			var hash []byte
			hash = Hash(value)
			n.Child[16] = hash
			if err := mpt.Savenode(n); err != nil {
				return nil, err
			}
			return n, nil
		}

		idx := nibbles[0]
		child, err := mpt.Loadnode(n.Child[idx])
		if err != nil {
			return nil, err
		}
		if child == nil {
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

func (mpt *MPT) has(node Node, nibbles []byte) ([]byte, error) {
	switch n := node.(type) {
	case *LeafNode:
		// 如果叶子节点的 key 为空，说明这是一个分支节点的直接子节点，直接返回其 value
		if len(n.Key) == 0 {
			return n.Value, nil
		}
		// 如果查找的key长度小于叶子节点的key长度，或者key不匹配，说明key不存在
		if len(nibbles) < len(n.Key) {
			return nil, fmt.Errorf("key not found4")
		}
		// 比较key的前缀
		if !bytes.Equal(n.Key, nibbles[:len(n.Key)]) {

		}
		// 如果查找的key长度与叶子节点的key长度不相等，说明不是完全匹配，key不存在
		if len(nibbles) != len(n.Key) {

			return nil, fmt.Errorf("key not found1")
		}
		return n.Value, nil

	case *ExtensionNode:

		// 如果查找的key长度小于扩展节点的路径长度，或者路径不匹配，说明key不存在
		if len(nibbles) < len(n.Key) {

			return nil, fmt.Errorf("key not found2")
		}
		// 比较路径前缀
		if !bytes.Equal(n.Key, nibbles[:len(n.Key)]) {

			return nil, fmt.Errorf("key not found3")
		}
		child, err := mpt.Loadnode(n.Value)
		if err != nil {
			return nil, err
		}
		return mpt.has(child, nibbles[len(n.Key):])

	case *BranchNode:

		if len(nibbles) == 0 {
			return n.Child[16], nil
		}
		idx := nibbles[0]
		if idx >= 16 {
			return nil, fmt.Errorf("invalid nibble value: %d", idx)
		}
		child, err := mpt.Loadnode(n.Child[idx])
		if err != nil {
			return nil, err
		}
		if child == nil {
			return nil, fmt.Errorf("key not found")
		}
		return mpt.has(child, nibbles[1:])

	default:
		return nil, fmt.Errorf("unknown node type")
	}
}
