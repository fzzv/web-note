package movezeroes

// 移动零
/*
题意：给定整数数组，原地将所有 0 移到末尾，同时保持非零元素的相对顺序。
思路：双指针。write 指向下一个应放置非零元素的位置，read 遍历数组：
- 如果 nums[read] 非零，则写到 nums[write]，并在 read 与 write 不同的情况下将 read 位置置 0，随后 write++。
- 如果 nums[read] 为 0，跳过。
遍历结束后，write 左侧已按原顺序存放所有非零元素，右侧自然为 0。
复杂度：时间 O(n)，空间 O(1)，只在原切片上操作。
*/
func MoveZeroes(nums []int) {
	write := 0
	for read, v := range nums {
		if v != 0 {
			if read != write {
				nums[write] = v
				nums[read] = 0
			}
			write++
		}
	}
}

func moveZeroes(nums []int) {
	if len(nums) == 0 {
		return
	}
	slow := 0
	for fast := range nums {
		// 只处理非 0 元素
		if nums[fast] != 0 {
			nums[slow] = nums[fast]
			slow++
		}
	}
	// slow 之后的元素全部补 0
	for i := slow; i < len(nums); i++ {
		nums[i] = 0
	}
}
