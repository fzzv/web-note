package twosumii

// 两数之和 II - 输入有序数组
/*
解析：
- 题意：给定有序数组和目标值，返回两个数的 1-based 下标，使其和为目标。每个输入只存在唯一解。
- 思路：双指针。left 从头，right 从尾：
  - 若 nums[left]+nums[right]==target，返回 [left+1, right+1]。
  - 若和小于 target，增大 left；若和大于 target，减小 right。
- 有序性保证指针单调收敛，时间 O(n)，空间 O(1)。
*/
func TwoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		if sum > target {
			right--
		} else if sum < target {
			left++
		} else {
			return []int{left + 1, right + 1}
		}
	}
	return nil
}
