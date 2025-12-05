package main

type MaxList struct {
	nums     []int
	max      int
	maxIndex int
}

func (m *MaxList) AddAt(i int) {
	if i < 0 || i > len(m.nums) {
		panic("index out of range")
	}

	m.nums[i]++

	v := m.nums[i]
	if v > m.max {
		m.max = v
		m.maxIndex = i
	}
}

func (m *MaxList) RemoveAt(i int) {
	if i < 0 || i > len(m.nums) {
		panic("index out of range")
	}

	m.nums[i]--

	if i == m.maxIndex {
		m.max = 0
		for idx, v := range m.nums {
			if v > m.max {
				m.max = v
				m.maxIndex = idx
			}
		}
	}
}

func characterReplacement(s string, k int) int {
	L := 0
	R := 0
	longest := 0

	alphabetList := MaxList{nums: make([]int, 26), max: 0, maxIndex: 0}

	for R < len(s) {
		alphabetList.AddAt(int(s[R] - 'A'))

		for (R - L + 1 - alphabetList.max) > k {
			alphabetList.RemoveAt(int(s[L] - 'A'))
			L++
		}

		w := R - L + 1
		if w > longest {
			longest = w
		}
		R++
	}

	return longest
}

func main() {
	result := characterReplacement("AABABBA", 1)
	println(result)
}
