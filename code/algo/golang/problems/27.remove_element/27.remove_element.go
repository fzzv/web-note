package removeelement

// 移除元素

/*
解析：
- 目标：原地移除所有等于 val 的元素，返回保留元素的新长度。
- 思路：使用写指针 write 指向下一个保留位置。遍历数组，遇到不等于 val 的元素时，写到 write 位置并递增 write；遇到等于 val 的元素跳过。写指针左侧即为保留区间，元素相对顺序不变。
- 复杂度：时间 O(n)，空间 O(1)，只在原切片上操作。
*/
func RemoveElement(nums []int, val int) int {
	write := 0
	for _, n := range nums {
		if n != val {
			nums[write] = n
			write++
		}
	}
	return write
}

func removeElement(nums []int, val int) int {
	slow := 0

	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != val {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
