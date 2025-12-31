package reversestring

// 反转字符串
// func ReverseString(s []byte) {
// 	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
// 		s[i], s[j] = s[j], s[i]
// 	}
// }

// func ReverseString(s []byte) {
// 	if len(s) == 0 {
// 		return
// 	}
// 	arr := append([]byte{}, s...)
// 	for i := range arr {
// 		s[i] = arr[len(arr)-i-1]
// 	}
// }

/*
解析：
- 题意：原地反转字符数组，函数直接修改入参切片。
- 思路：双指针 left/right 从两端向中间收缩，交换对应元素，直到相遇或交错。
- 特点：仅交换，不创建新切片，满足 O(1) 额外空间；每个元素最多参与一次交换，时间 O(n)。
*/
func ReverseString(s []byte) {
	left := 0
	right := len(s) - 1

	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}
