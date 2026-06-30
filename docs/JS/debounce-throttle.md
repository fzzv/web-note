# 函数的防抖和节流

## 作用

防抖：触发事件后，不立即执行，等待 N 毫秒。如果 N 毫秒内再次触发，则重新计时。

节流：固定时间内只执行一次。

## 防抖和节流的区别

| 对比     | 防抖           | 节流             |
| -------- | -------------- | ---------------- |
| 英文     | Debounce       | Throttle         |
| 核心     | 只执行最后一次 | 固定时间执行一次 |
| 连续触发 | 不执行         | 周期执行         |
| 停止触发 | 执行一次       | 不一定           |
| 适用场景 | 搜索框         | scroll、resize   |

## 手写函数的防抖和节流

```js
/* 函数的防抖和节流 */
const clearTimer = function clearTimer(timer) {
    if (timer) clearTimeout(timer);
    return null;
};

// 函数防抖：防止“老年帕金森”，用户频繁触发某个行为的时候，我们只识别“一次”「频繁的定义可以自己管控」
const debounce = function debounce(func, wait, immediate = false) {
    // init params
    if (typeof func !== "function") throw new TypeError("func is not a function!");
    if (typeof wait === "boolean") {
        immediate = wait;
        wait = undefined;
    }
    wait = +wait;
    if (isNaN(wait)) wait = 300;
    if (typeof immediate !== "boolean") immediate = false;
    // handler
    let timer = null;
    return function operate(...params) {
        // now:记录是否是立即执行「第一次点击&immediate=true」
        let now = !timer && immediate;
        // 清除之前设置的定时器
        timer = clearTimer(timer);
        timer = setTimeout(() => {
            // 结束边界触发
            if (!immediate) func.call(this, ...params);
            // 清除最后一个定时器
            timer = clearTimer(timer);
        }, wait);
        // 如果是立即执行，则第一次执行operate就把要干的事情做了即可 “开始边界触发”
        if (now) func.call(this, ...params);
    };
};

// 时间戳版适合首次立即响应，定时器版适合保证最后一次执行。实际项目中通常采用两者结合的方式，同时支持 leading 和 trailing
// 时间戳+定时器结合版本
// 函数节流：用户频繁操作的时候，不根据用户的频繁操作度来绝定触发多少次，而是根据设定好的频率进行触发，实现“降频”的效果，相对于防抖来讲，节流是允许触发多次的
const throttle = function throttle(func, wait) {
    // init params
    if (typeof func !== "function") throw new TypeError("func is not a function!");
    wait = +wait;
    if (isNaN(wait)) wait = 300;
    // handler
    let timer = null,
        previous = 0; //记录上一次func触发的时间
    return function operate(...params) {
        let now = +new Date(),
            remaining = wait - (now - previous);
        if (remaining <= 0) {
            // 两次触发的间隔时间超过设定的频率，则立即执行函数
            func.call(this, ...params);
            previous = +new Date();
            timer = clearTimer(timer);
        } else if (!timer) {
            // 间隔时间不足设定的频率，而且还未设置等待的定时器，则设置定时器等待执行函数即可
            timer = setTimeout(() => {
                func.call(this, ...params);
                previous = +new Date();
                timer = clearTimer(timer);
            }, remaining);
        }
    };
};
```

