package pointers

import (
	"errors"
	"fmt"
)

type Bitcoin int

type Wallet struct {
	balance Bitcoin
}

/*
// 如果不使用指针，会导致钱包余额的拷贝，而不是引用。
func (w Wallet) Balance(amount Bitcoin) {
	w.balance += amount
}
*/

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w Wallet) Balance() Bitcoin {
	return w.balance
}

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return errors.New("cannot withdraw, insufficient funds")
	}
	w.balance -= amount
	return nil
}

/*
类型别名有一个有趣的特性，你还可以对它们声明 方法。
当你希望在现有类型之上添加一些领域内特定的功能时，这将非常有用。
*/
// 当测试出错时，会在数字后加上 BTC 字符串。
// 比如 got 10 BTC want 11 BTC
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
