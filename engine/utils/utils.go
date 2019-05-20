package utils

func Max(l []int) (max int) {
	max = l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return
}

func Min(l []int) (min int) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
}
