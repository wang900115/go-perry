package heap

type Heap struct {
	data []int
}

func (h *Heap) Insert(val int) {
	h.data = append(h.data, val)
	h.HeapifyUp(len(h.data) - 1)
}

func (h *Heap) HeapifyUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if h.data[parent] <= h.data[index] {
			break
		}
		h.data[parent], h.data[index] = h.data[index], h.data[parent]
		index = parent
	}
}

func (h *Heap) Remove() (int, bool) {
	if len(h.data) == 0 {
		return 0, false
	}
	min := h.data[0]
	h.data[0] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	h.HeapifyDown(0)
	return min, true
}

func (h *Heap) HeapifyDown(index int) {
	left := index*2 + 1
	right := index*2 + 2
	smallest := index
	if left < len(h.data) && h.data[left] < h.data[smallest] {
		smallest = left
	}

	if right < len(h.data) && h.data[right] < h.data[smallest] {
		smallest = right
	}

	if smallest != index {
		h.data[smallest], h.data[index] = h.data[index], h.data[smallest]
		h.HeapifyDown(smallest)
	}
}

func (h *Heap) Peek() (int, bool) {
	if len(h.data) == 0 {
		return 0, false
	}
	return h.data[0], true
}

func (h *Heap) Size() int {
	return len(h.data)
}
