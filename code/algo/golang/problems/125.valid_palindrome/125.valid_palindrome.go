package validpalindrome

// 验证回文串
/*
题意：判断一个字符串是否为回文，仅考虑字母和数字字符，忽略大小写。
思路：双指针。left/right 从两端向中间收缩，遇到非字母数字就跳过；比较时统一转小写。
复杂度：时间 O(n)，空间 O(1)，在原字符串上双向扫描。
*/
func IsPalindrome(s string) bool {
	left, right := 0, len(s)-1

	for left < right {
		// 左指针跳过非法字符
		for left < right && !isAlnum(s[left]) {
			left++
		}
		// 右指针跳过非法字符
		for left < right && !isAlnum(s[right]) {
			right--
		}

		if left >= right {
			break
		}

		if toLower(s[left]) != toLower(s[right]) {
			return false
		}
		left++
		right--
	}

	return true
}

func isAlnum(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

func isPalindrome(s string) bool {
	lowerSlice := []byte{}
	sSlice := []byte(s)

	for i := range sSlice {
		if isAlnum(sSlice[i]) {
			lowerSlice = append(lowerSlice, toLower(sSlice[i]))
		}
	}

	left := 0
	right := len(lowerSlice) - 1

	for left < right {
		if lowerSlice[left] != lowerSlice[right] {
			return false
		}
		left++
		right--
	}

	return true
}
