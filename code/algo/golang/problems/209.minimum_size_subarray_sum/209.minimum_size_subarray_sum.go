package minsubsum

// 长度最小的子数组
/*
解析：
- 题意：给定正整数数组和目标和 target，求和至少为 target 的最短连续子数组长度；若不存在返回 0。
- 思路：滑动窗口。用 left/right 维护窗口和 sum：
  - 右指针扩展累加 sum，直到 sum >= target，尝试收缩左指针以取得最短长度，同时减去移出值。
  - 每次满足条件更新最小长度。
- 因为元素全为正数，窗口收缩不会错过更优解，整体线性扫描。
- 复杂度：时间 O(n)，空间 O(1)。
*/
func MinSubArrayLen(target int, nums []int) int {
	n := len(nums)
	minLen := n + 1

	sum := 0
	left := 0
	for right, v := range nums {
		sum += v
		for sum >= target {
			if currentLen := right - left + 1; currentLen < minLen {
				minLen = currentLen
			}
			sum -= nums[left]
			left++
		}
	}

	if minLen == n+1 {
		return 0
	}
	return minLen
}
