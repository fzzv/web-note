# 工作中的一些记录

## 可编辑表格中的日期不可以设置

ant design vue 可编辑表格 日期列需要从上到下依次递增，设置不可选范围

```js
// 设置不可选日期
/**
* tableData 表格数据
* current 当前行数据
* index 下标
*/
function disabledDate(tableData, current, index) {
    // 获取前面所有有效日期
    const prevDates = tableData
    .slice(0, index)
    .map(item => item.date)
    const maxPrevDate = prevDates.length && prevDates.some(item => !!item) ? moment.max(prevDates                                                              .filter(date => !!date)
.map(date => moment(date))) : null
    // 获取后面所有有效日期
    const nextDates = tableData
    .slice(index + 1)
    .map(item => item.date)
    const minNextDate = nextDates.length && nextDates.some(item => !!item) ? moment.min(nextDates
.filter(date => !!date)
.map(date => moment(date))) : null
    let disabled = false
    // 当前日期早于前面的最大日期 → 禁用
    if (maxPrevDate && current.isBefore(maxPrevDate, 'day')) {
        disabled = true
    }
    // 当前日期晚于后面的最小日期 → 禁用
    if (minNextDate && current.isAfter(minNextDate, 'day')) {
        disabled = true
    }
    return disabled
}
```

