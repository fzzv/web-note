package sortanarray

import "fmt"

// 排序数组（LeetCode 912）
/*
题目：
- 给定整数数组 nums，请将其按升序排序并返回。
- 约束：1 <= nums.length <= 5*10^4，-5*10^4 <= nums[i] <= 5*10^4。

解析：
- 需求：返回有序数组，允许原地修改。
- 思路：自顶向下归并排序，稳定且保证 O(n log n)：
  - 递归拆分至单元素；
  - 合并时用辅助数组缓存，并按大小写回原数组；
  - 若左半最大值已不大于右半最小值，可跳过合并以减少常数。
- 复杂度：时间 O(n log n)，空间 O(n) 辅助数组。
*/
// func SortArray(nums []int) []int {
// 	if len(nums) <= 1 {
// 		return nums
// 	}

// 	buf := make([]int, len(nums))
// 	mergeSort(nums, buf, 0, len(nums)-1)
// 	return nums
// }

// func mergeSort(nums, buf []int, left, right int) {
// 	if left >= right {
// 		return
// 	}

// 	mid := left + (right-left)/2
// 	mergeSort(nums, buf, left, mid)
// 	mergeSort(nums, buf, mid+1, right)

// 	// 若两段本就有序，直接返回
// 	if nums[mid] <= nums[mid+1] {
// 		return
// 	}

// 	i, j, k := left, mid+1, left
// 	for i <= mid && j <= right {
// 		if nums[i] <= nums[j] {
// 			buf[k] = nums[i]
// 			i++
// 		} else {
// 			buf[k] = nums[j]
// 			j++
// 		}
// 		k++
// 	}

// 	for i <= mid {
// 		buf[k] = nums[i]
// 		i++
// 		k++
// 	}
// 	for j <= right {
// 		buf[k] = nums[j]
// 		j++
// 		k++
// 	}

// 	copy(nums[left:right+1], buf[left:right+1])
// }

func SortArray(nums []int) []int {
	// 1. 递归终止条件（Base Case）
	// 如果数组只有一个元素或者为空，它自然是有序的，直接返回
	if len(nums) <= 1 {
		return nums
	}

	// 2. 找中点，切分
	mid := len(nums) / 2

	// 3. 递归地对左半部分和右半部分进行排序
	// 注意：这里会一直递归下去，直到切成单个元素
	left := SortArray(nums[:mid])
	right := SortArray(nums[mid:])
	fmt.Printf("left: %v\n", left)
	fmt.Printf("right: %v\n", right)

	// 4. 将排好序的两部分“合并”起来
	return merge(left, right)
}

// 辅助函数：负责“合” (拉拉链逻辑)
func merge(left, right []int) []int {
	// 预分配结果数组的空间，避免 append 频繁扩容
	result := make([]int, 0, len(left)+len(right))

	i, j := 0, 0

	// 当两个切片都还有元素时，比较头部元素
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// 把剩下的元素接在后面
	// (因为上面的循环结束后，肯定有一个切片还有剩余，且剩余部分一定是有序且最大的)
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
