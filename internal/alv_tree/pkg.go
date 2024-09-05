package alv_tree

// 获取节点的高度
func height(node *AlvNode) int {
	if node == nil {
		return 0
	}
	return node.height
}

// 获取平衡因子
func balanceFactor(node *AlvNode) int {
	if node == nil {
		return 0
	}
	return height(node.lchild) - height(node.rchild)
}

// 更新节点的高度
func updateHeight(node *AlvNode) {
	if node != nil {
		node.height = 1 + max(height(node.lchild), height(node.rchild))
	}
}

// 右旋
func rightRotate(y *AlvNode) *AlvNode {
	x := y.lchild
	T2 := x.rchild

	x.rchild = y
	y.lchild = T2

	updateHeight(y)
	updateHeight(x)

	return x
}

// 左旋
func leftRotate(x *AlvNode) *AlvNode {
	y := x.rchild
	T2 := y.lchild

	y.lchild = x
	x.rchild = T2

	updateHeight(x)
	updateHeight(y)

	return y
}

// 获取较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
