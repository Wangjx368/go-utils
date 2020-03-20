package math

// MyAbsSortList 自定义的根据绝对值进行从小到大排序
type MyAbsSortList []int

func (a MyAbsSortList) Len() int           { return len(a) }
func (a MyAbsSortList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MyAbsSortList) Less(i, j int) bool { return abs(a[i]) < abs(a[j]) }
func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
