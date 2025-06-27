package mpt

import (
	"encoding/json"
	"fmt"
)

type NodeType int

const (
	LeafNodeType NodeType = iota
	ExtensionNodeType
	BranchNodeType
)

type (
	LeafNode struct {
		NodeType NodeType `json:"nodeType"`
		Key      []byte   `json:"key"`
		Value    []byte   `json:"value"`
	}

	ExtensionNode struct {
		NodeType NodeType `json:"nodeType"`
		Key      []byte   `json:"key"`
		Value    []byte   `json:"value"`
	}

	BranchNode struct {
		NodeType NodeType   `json:"nodeType"`
		Child    [17][]byte `json:"child"`
	}
)

type Node interface {
    Serialize() ([]byte, error)          // 序列化为 JSON 字符串
    Deserialize(jsonBytes []byte) error    // 从 JSON 字符串反序列化
    GetNodeType() NodeType
    GetHash() []byte
}

func (n *LeafNode) GetHash() []byte {
    serialized, err := n.Serialize()
    if err != nil {
    return nil
    }
    return Hash(serialized)
}

func (n *ExtensionNode) GetHash() []byte {
    serialized, err := n.Serialize()
    if err != nil {
    return nil
    }
    return Hash(serialized)
}

func (n *BranchNode) GetHash() []byte {
    serialized, err := n.Serialize()
    if err != nil {
    return nil
    }
    return Hash(serialized)
}

func (n *LeafNode) Serialize() ([]byte, error) {
    jsonBytes, err := json.Marshal(n)
    if err != nil {
        return nil, err
    }
    return jsonBytes, nil
}

func (n *LeafNode) Deserialize(jsonBytes []byte) error {
    return deserializeNode(jsonBytes, n)
}

func (n *LeafNode) GetNodeType() NodeType {
    return n.NodeType
}

func (n *ExtensionNode) Serialize() ([]byte, error) {
    jsonBytes, err := json.Marshal(n)
    if err != nil {
        return nil, err
    }
    return jsonBytes, nil
}

func (n *ExtensionNode) Deserialize(jsonBytes []byte) error {
    return deserializeNode(jsonBytes, n)
}

func (n *ExtensionNode) GetNodeType() NodeType {
    return n.NodeType
}

func (n *BranchNode) Serialize() ([]byte, error) {
    jsonBytes, err := json.Marshal(n)
    if err != nil {
        return nil, err
    }
    return jsonBytes, nil
}

func (n *BranchNode) Deserialize(jsonBytes []byte) error {
    return deserializeNode(jsonBytes, n)
}

func (n *BranchNode) GetNodeType() NodeType {
    return n.NodeType
}

func deserializeNode(jsonBytes []byte, node interface{}) error {
    return json.Unmarshal(jsonBytes, node)
}

func DeserializeNode(jsonBytes []byte, node interface{}) (Node, error) {
    var nodeType struct {
        NodeType NodeType `json:"NodeType"`
    }
    
    if err := json.Unmarshal(jsonBytes, &nodeType); err != nil {
        return nil, err
    }

    var result Node
    switch nodeType.NodeType {
    case LeafNodeType:
        leaf := &LeafNode{}
        if err := json.Unmarshal(jsonBytes, leaf); err != nil {
            return nil, err
        }
        result = leaf
    case ExtensionNodeType:
        ext := &ExtensionNode{}
        if err := json.Unmarshal(jsonBytes, ext); err != nil {
            return nil, err
        }
        result = ext
    case BranchNodeType:
        branch := &BranchNode{}
        if err := json.Unmarshal(jsonBytes, branch); err != nil {
            return nil, err
        }
        result = branch
    default:
        return nil, fmt.Errorf("unknown node type")
    }
    
    return result, nil
}