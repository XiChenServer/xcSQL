package affairs

// WaterMark 结构体
type WaterMark struct {
	heap []Txn // 事务开始和结束时间戳的堆
}

func (w WaterMark) Len() int {
	return len(w.heap)
}

func (w WaterMark) Less(i, j int) bool {
	// 按开始时间戳升序排序
	return w.heap[i].startTime < w.heap[j].startTime
}

func (w WaterMark) Swap(i, j int) {
	w.heap[i], w.heap[j] = w.heap[j], w.heap[i]
}

func (w *WaterMark) Push(x interface{}) {
	w.heap = append(w.heap, x.(Txn))
}

func (w *WaterMark) Pop() interface{} {
	old := w.heap
	n := len(old)
	x := old[n-1]
	w.heap = old[0 : n-1]
	return x
}
