package disjoint

import "sort"

type SummaryRanges struct {
	intervals map[int]int
}

func Constructor() SummaryRanges {
	return SummaryRanges{
		intervals: make(map[int]int),
	}
}

func (this *SummaryRanges) AddNum(value int) {
	if _, ok := this.intervals[value]; ok {
		return
	}
	start, end := value, value

	for s, e := range this.intervals {

		if value >= s && value <= e {
			return
		}

		if e+1 == value {
			start = s
			delete(this.intervals, s)
		} else if s-1 == value {
			end = e
			delete(this.intervals, s)
		}
	}
	this.intervals[start] = end
}

func (this *SummaryRanges) GetIntervals() [][]int {
	keys := make([]int, 0, len(this.intervals))
	for k := range this.intervals {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	res := [][]int{}
	for _, k := range keys {
		res = append(res, []int{k, this.intervals[k]})
	}
	return res
}
