package alv_tree

import "fmt"

// 插入操作
func insert(root *AlvNode, key []byte, node nodeInfo) *AlvNode {
	if root == nil {
		return &AlvNode{data: node, height: 1}
	}

	if string(key) < string(root.data.Key) {
		root.lchild = insert(root.lchild, key, node)
	} else if string(key) > string(root.data.Key) {
		root.rchild = insert(root.rchild, key, node)
	} else {
		if root.data.TTL == -1 {
			// 如果节点过期，则跳过插入
			return root
		}
		root.data = node
		return root
	}

	updateHeight(root)

	bf := balanceFactor(root)

	if bf > 1 && string(key) < string(root.lchild.data.Key) {
		return rightRotate(root)
	}

	if bf < -1 && string(key) > string(root.rchild.data.Key) {
		return leftRotate(root)
	}

	if bf > 1 && string(key) > string(root.lchild.data.Key) {
		root.lchild = leftRotate(root.lchild)
		return rightRotate(root)
	}

	if bf < -1 && string(key) < string(root.rchild.data.Key) {
		root.rchild = rightRotate(root.rchild)
		return leftRotate(root)
	}

	return root
}

// 搜索
func search(root *AlvNode, key []byte) (*nodeInfo, bool) {
	if root == nil {
		return nil, false
	}

	if string(key) < string(root.data.Key) {
		return search(root.lchild, key)
	} else if string(key) > string(root.data.Key) {
		return search(root.rchild, key)
	} else {
		return &root.data, true
	}
}
func minNode(node *AlvNode) *AlvNode {
	current := node
	for current.lchild != nil {
		current = current.lchild
	}
	return current
}

// 删除操作
func delete(root *AlvNode, key []byte) *AlvNode {
	if root == nil {
		return nil
	}

	if string(key) < string(root.data.Key) {
		root.lchild = delete(root.lchild, key)
	} else if string(key) > string(root.data.Key) {
		root.rchild = delete(root.rchild, key)
	} else {
		// 设置 TTL 为 -1 表示过期
		root.data.TTL = -1
	}

	// 无需重新平衡树，因为没有实际删除节点
	return root
}

// 先序遍历
func printPreOrder(node *AlvNode) {
	if node == nil {
		return
	}
	fmt.Printf("%s ", node.data.Key)
	printPreOrder(node.lchild)
	printPreOrder(node.rchild)
}
