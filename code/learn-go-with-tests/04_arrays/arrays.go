package arrays

func Sum(arr []int) int {
	sum := 0
	for i := range arr {
		sum += arr[i]
	}
	return sum
}

/*
func SumAll(numbersToSum ...[]int) (sums []int) {}
相当于

	func SumAll(numbersToSum ...[]int) []int {
		var sums []int
	}

并且在return sums时，sums也是可以省略的
*/
func SumAll(numbersToSum ...[]int) (sums []int) {
	// 通过 len 获取需要创建的切片或数组的长度
	numsLength := len(numbersToSum)
	// make 可以在创建切片的时候指定我们需要的长度和容量
	sums = make([]int, numsLength)
	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}
	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
			continue
		}
		// 获取除第一个元素外的所有元素
		tail := numbers[1:]
		// 切片的容量是 固定 的，但是你可以使用 append 从原来的切片中创建一个新切片
		sums = append(sums, Sum(tail))
	}
	return sums
}
