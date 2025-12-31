package twosum

// 两数之和
/*
解析：
- 题意：找到数组中两个数，使其和为 target，返回对应下标。
- 思路：一次遍历 + 哈希。遍历当前元素 n 时，先检查哈希表中是否已有补数 target-n；若有则直接返回两下标；否则把当前值与下标写入表，继续扫描。
- 复杂度：时间 O(n)，空间 O(n) 存储已遍历元素。
*/
func TwoSum(nums []int, target int) []int {
	indexByValue := make(map[int]int, len(nums))

	for i, n := range nums {
		if j, ok := indexByValue[target-n]; ok {
			return []int{j, i}
		}
		indexByValue[n] = i
	}

	return nil
}

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for i := range nums {
		current := nums[i]

		need := target - current

		if index, ok := m[need]; ok {
			return []int{index, i}
		}

		m[current] = i
	}
	return nil
}
