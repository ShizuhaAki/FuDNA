package sort

type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func Sort(data Interface) {
	if data.Len() <= 1 {
		return
	}
	quickSort(data, 0, data.Len()-1)
}

func quickSort(data Interface, low, high int) {
	if low < high {
		pi := partition(data, low, high)

		quickSort(data, low, pi-1)
		quickSort(data, pi+1, high)
	}
}

func partition(data Interface, low, high int) int {
	i := low - 1 // Index of smaller element

	for j := low; j < high; j++ {
		if data.Less(j, high) {
			i++
			data.Swap(i, j)
		}
	}

	data.Swap(i+1, high)

	return i + 1
}

type IntSlice []int

func (s IntSlice) Len() int           { return len(s) }
func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func Ints(a []int) {
	Sort(IntSlice(a))
}
