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

## 下载

点击下载功能

```js
fileDownload(row) {
    if (row.attId) {
        // 根据id进行下载
        downloadFile(row.attId).then((res) => {
            // 获取响应体中的信息,里面有文件名
            const str = res.headers['content-disposition']
            if (str) {
                const index = str.lastIndexOf('=')
                // 获取文件名
                const filename = window.decodeURI(str.substring(index + 1, str.length))
                this.handleFileDownloadRes(res, filename)
            } else {
                this.$message.error('文件信息不存在')
            }
        })
    }
},
handleFileDownloadRes(res, filename) {
    if (!res.data) {
        this.$message.error('文件信息不存在')
        return
    }
    if (window.navigator && window.navigator.msSaveOrOpenBlob) {
        // 检测是否在IE浏览器打开
        window.navigator.msSaveOrOpenBlob(new Blob([res.data]), filename)
    } else {
        // 谷歌、火狐浏览器
        let url = ''
        if (
            window.navigator.userAgent.indexOf('Chrome') >= 1 ||
            window.navigator.userAgent.indexOf('Safari') >= 1
        ) {
            url = window.webkitURL.createObjectURL(new Blob([res.data]))
        } else {
            url = window.URL.createObjectURL(new Blob([res.data]))
        }
        const link = document.createElement('a')
        const iconv = require('iconv-lite')
        iconv.skipDecodeWarning = true // 忽略警告
        link.style.display = 'none'
        link.href = url
        link.setAttribute('download', filename)
        document.body.appendChild(link)
        link.click()
        link.remove()
        window.URL.revokeObjectURL(url)
    }
}
```

