package maps

import "errors"

var (
	ErrWordExists       = errors.New("cannot add word because it already exists")
	ErrWordNotFound     = errors.New("could not find the word you were looking for")
	ErrWordDoesNotExist = errors.New("cannot update word because it does not exist")
)

// 创建一个 Dictionary 类型，它是对 map 的简单封装。
// 这样我们就可以为 Dictionary 类型添加方法了。
type Dictionary map[string]string

var ErrNotFound = errors.New("could not find the word you were looking for")

// 为 Dictionary 类型添加 Search 方法。
// 这个方法接受一个单词，返回它的定义和一个错误。
// 如果单词不存在，返回 ErrNotFound。
func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}
	return definition, nil
}

/*
Map 有一个有趣的特性，不使用指针传递你就可以修改它们。这是因为 map 是引用类型。
这意味着它拥有对底层数据结构的引用，就像指针一样。
它底层的数据结构是 hash table 或 hash map。
*/
func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		d[word] = definition
		return nil
	case nil:
		return ErrWordExists
	default:
		return err
	}
}

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		return ErrWordNotFound
	case nil:
		d[word] = definition
		return nil
	default:
		return err
	}
}

func (d Dictionary) Delete(word string) {
	// map 有一个内置函数 delete。它需要两个参数。第一个是这个 map，第二个是要删除的键。
	delete(d, word)
}
