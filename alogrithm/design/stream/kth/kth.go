package kth

// min-heap
type KthLargest struct {
	nums     []int
	capacity int
}

func Constructor(k int, nums []int) KthLargest {
	kl := &KthLargest{
		nums:     make([]int, 0, k),
		capacity: k,
	}
	for _, num := range nums {
		kl.Add(num)
	}
	return *kl
}

func (kl *KthLargest) Add(val int) int {
	if kl.Len() < kl.capacity {
		kl.nums = append(kl.nums, val)
		kl.upHeap(kl.Len() - 1)
		return kl.peak()
	}

	if val < kl.peak() {
		return kl.peak()
	}

	kl.nums[0] = val
	kl.downHeap(0)
	return kl.peak()
}

func (kl *KthLargest) peak() int {
	return kl.nums[0]
}

func (kl *KthLargest) Len() int {
	return len(kl.nums)
}

func (kl *KthLargest) upHeap(index int) {
	for {
		parent := (index - 1) / 2
		if index == 0 || kl.nums[parent] <= kl.nums[index] {
			break
		}
		kl.nums[parent], kl.nums[index] = kl.nums[index], kl.nums[parent]
		index = parent
	}
}

func (kl *KthLargest) downHeap(index int) {
	n := kl.capacity
	for {
		smallest := index
		left := 2*index + 1
		right := 2*index + 2
		if left < n && kl.nums[left] < kl.nums[smallest] {
			smallest = left
		}
		if right < n && kl.nums[right] < kl.nums[smallest] {
			smallest = right
		}
		if smallest == index {
			break
		}
		kl.nums[smallest], kl.nums[index] = kl.nums[index], kl.nums[smallest]
		index = smallest
	}
}
