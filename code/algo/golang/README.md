# LeetCode in Go

Practice LeetCode-style problems in Go with table-driven tests.

## Layout
- `problems/<id>.<slug>/`: one folder per problem with solution and `_test.go` files, prefixed by LeetCode id.
- `problems/1.two_sum/`: solution for "Two Sum" (`TwoSum`) plus tests.
- `problems/26.remove_duplicates_from_sorted_array/`: solution for "Remove Duplicates from Sorted Array" (`RemoveDuplicates`，删除有序数组中的重复项) plus tests.
- `problems/27.remove_element/`: solution for "Remove Element" (`RemoveElement`，移除元素) plus tests.
- `problems/125.valid_palindrome/`: solution for "Valid Palindrome" (`IsPalindrome`，验证回文串) plus tests.
- `problems/167.two_sum_ii/`: solution for "Two Sum II - Input array is sorted" (`TwoSum`，两数之和 II - 输入有序数组) plus tests.
- `problems/209.minimum_size_subarray_sum/`: solution for "Minimum Size Subarray Sum" (`MinSubArrayLen`，长度最小的子数组) plus tests.
- `problems/283.move_zeroes/`: solution for "Move Zeroes" (`MoveZeroes`，移动零) plus tests.
- `problems/344.reverse_string/`: solution for "Reverse String" (`ReverseString`，反转字符串) plus tests.
- `problems/414.third_maximum_number/`: solution for "Third Maximum Number" (`ThirdMax`) plus tests.
- `problems/977.squares_of_sorted_array/`: solution for "Squares of a Sorted Array" (`SortedSquares`，有序数组的平方) plus tests.

## Usage
Run all problems from the repo root:

```bash
go test ./...
```

## Adding a new problem
1. Create a folder under `problems/` (e.g., `problems/three_sum/`).
2. Implement the solution in `*.go` using a clear function name.
3. Add table-driven tests in `*_test.go` to cover edge cases.
4. Run `go test ./...` to verify.
