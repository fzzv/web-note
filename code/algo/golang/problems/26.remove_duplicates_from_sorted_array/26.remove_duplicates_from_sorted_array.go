package removeduplicates

// 删除有序数组中的重复项
/*
解析：
- 题意：给定非降序数组，原地删除重复项，使每个元素只出现一次，返回新长度。
- 思路：双指针（快读 slow写）。write 指向已保留区间的末尾，read 向前扫描；当发现新值（与 nums[write] 不同）时，将其写到 write+1 并推进 write。
- 特点：保持相对顺序，操作仅发生在原切片上。
- 复杂度：时间 O(n)，空间 O(1)。
*/
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	write := 0
	for read := 1; read < len(nums); read++ {
		if nums[read] != nums[write] {
			write++
			nums[write] = nums[read]
		}
	}
	return write + 1
}

func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slot := 1
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slot] = nums[fast]
			slot++
		}
	}
	return slot
}
