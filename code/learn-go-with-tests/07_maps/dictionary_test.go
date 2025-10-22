package maps

import "testing"

/*
声明 map 的方式有点儿类似于数组。不同之处是，它以 map 关键字开头，需要两种类型。
第一个是键的类型，写在 [] 中。
第二个是值的类型，跟在 [] 之后。
map 的键是唯一的，如果键重复，会覆盖之前的值。
map 的键是无序的，不能通过索引访问。

	func TestMap(t *testing.T) {
		dictionary := map[string]string{
			"test": "this is a test",
		}
		got := Search(dictionary, "test")
		want := "this is a test"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	func TestMap(t *testing.T) {
		dictionary := Dictionary{
			"test": "this is a test",
		}
		got := dictionary.Search("test")
		want := "this is a test"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
*/

func assertDefinition(t *testing.T, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if definition != got {
		t.Errorf("got '%s' want '%s'", got, definition)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}

// 搜索字典
func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is a test"}

	t.Run("known word", func(t *testing.T) {
		_, err := dictionary.Search("test")
		want := "this is a test"
		assertDefinition(t, dictionary, "test", want)
		assertError(t, err, nil)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")
		assertError(t, err, ErrNotFound)
	})
}

/*
// 添加单词
func TestAdd(t *testing.T) {
	dictionary := Dictionary{}
	dictionary.Add("favorite", "coding")
	got, err := dictionary.Search("favorite")
	want := "coding"
	assertDefinition(t, got, want)
	assertError(t, err, nil)
}
*/

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "test"
		definition := "this is just a test"

		err := dictionary.Add(word, definition)

		assertError(t, err, nil)
		assertDefinition(t, dictionary, word, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}
		err := dictionary.Add(word, "new test")

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dictionary, word, definition)
	})
}

// 更新单词
func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		newDefinition := "new definition"
		dictionary := Dictionary{word: definition}

		err := dictionary.Update(word, newDefinition)

		assertError(t, err, nil)
		assertDefinition(t, dictionary, word, newDefinition)
	})

	t.Run("new word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{}

		err := dictionary.Update(word, definition)

		// 更精准的错误提示
		// assertError(t, err, ErrWordDoesNotExist)

		assertError(t, err, ErrWordNotFound)
	})
}

// 删除单词
func TestDelete(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		dictionary := Dictionary{word: "test definition"}
		dictionary.Delete(word)
		_, err := dictionary.Search(word)
		assertError(t, err, ErrNotFound)
	})
}
