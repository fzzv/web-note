# Grid 布局

```css
// 占用布局的列数，从第n列开始
grid-column-start: 1; 
// 从第n-1列结束,例如为6，则在第5列结束
grid-column-end: 6;
// 简写方式 grid-column-start / grid-column-end
grid-column: 1 / 6;
// 占用布局的行数，从第n行开始
grid-row-start: 1;
// 从第n-1行结束
grid-row-end: 6;
// 简写方式 grid-row-start / grid-row-end
grid-row: 1 / 6;
// 简写方式 (顺序为grid-row-start/grid-column-start/grid-row-end/grid-row-end)
grid-area: 1 / 2 / 3 / 4;
// 调整顺序
order: 0;
// 布局的列数：px / fr / repeat / auto / em / ch ...
grid-template-columns: 50% 50%; // 两列
grid-template-columns: repeat(5, 20%); // 五列宽度为20%的格子
// 布局的行数
grid-template-rows: 50% 50%; // 两行
// 简写方式 grid-template-columns / grid-template-rows
grid-template: 60% / 50px;
```



`gird`新单位 **fr**：分数fr。每个fr单元分配一份可用空间。例如：两个元素分别设为 1fr 和 3fr，则空间会被分为4等份。第一个元素占 1/4，第二个元素占 3/4。



## 容器查询

语法：

```css
@container <container-condition> {  
    <block-contents>
}
```

- `<container-condition>` 为查询条件，如果匹配到则会应用`<block-contents>` 的样式。它由两部分组成 `<container-name>` 和 `<container-query>`。
  - `<container-name>` 指定作用的容器名称（可选）。如果不写，匹配到的容器则取决于 `<block-contents>` 里的元素
  - `<container-query>`指定查询条件，比如：`(width > 1200px)`、`style(--theme dark)` ......
- `<block-contents>` 容器里子元素的样式，跟`@media` 一样

`container-name: my-container;` 用于指定容器名称

`container-type: inline-size;` 仅查询元素的宽度，size 还可以查询元素的高度和 aspect-ratio

`container: my-container/inline-size;`  为container-name和container-type的简写

> `container-type: inline-size` 可以使用 `min-width` 和 `max-width`，但不可以使用 `min-height` 和 `max-height`

### 单位

cqw: Container query width 容器查询宽度，1cqw等于容器宽度的1%

cqh: Container query height 容器查询高度，1cqh等于容器高度的1%

cqi: Container query inline size 内联大小，1cqi等于容器内联大小的1%

cqb: Container query block size 块大小，1cqb等于容器块大小的1%

cqmin: Container query minimum size 最小的size，cqi/cqb的值，哪个更小取哪个

cqmax: Contaienr query maximum size 最大的size，cqi/cqb的值，哪个更大取哪个

