package main

// 创建并返回一个内存中的玩家分数存储实例
// 初始化一个空的 map[string]int 来保存玩家名称与对应的得分
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

// 定义一个内存中的玩家分数存储结构体
// 其中 key 是玩家名字，value 是该玩家的得分
type InMemoryPlayerStore struct {
	store map[string]int
}

// 记录指定玩家的获胜次数
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

// 获取指定玩家当前的获胜次数
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league
}
