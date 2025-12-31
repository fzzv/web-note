package sortedsquares

// 有序数组的平方
/*
解析：
- 题意：给定非递减数组，返回每个元素平方后按非递减顺序排列的数组。
- 思路：双指针从两端向中间。最小绝对值在中间，最大平方在两端。用 left/right 指向首尾，idx 从末尾向前填充结果：
  - 比较 nums[left]^2 与 nums[right]^2，将较大者放入 res[idx]，对应指针向内收缩。
  - 如此保证 res 由大到小填充，最终得到升序结果。
- 复杂度：时间 O(n)，空间 O(n) 结果数组（题目允许返回新数组）。
*/
func SortedSquares(nums []int) []int {
	numsLen := len(nums)
	res := make([]int, numsLen)
	left, right, position := 0, numsLen-1, numsLen-1
	// 递增排序的数组，需要从右侧先填充最大值
	for left <= right {
		leftValue := nums[left] * nums[left]
		rightValue := nums[right] * nums[right]

		if leftValue > rightValue {
			res[position] = leftValue
			left++
		} else {
			res[position] = rightValue
			right--
		}
		position--
	}
	return res
}
