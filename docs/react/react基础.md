# React基础

## React 三个特点

- 1 声明式  ==> 命令式编程    arr.filter(item => item.price>80)
  - 利用JSX 语法来声明描述动态页面， 数据更新界面自动更新
  - 我们不用亲自操作DOM, 只需要更新数据, 界面就会自动更新
  - React.createElement() 是命令式
- 2 组件化
  - 将一个较大较复杂的界面拆分成几个可复用的部分封装成多个组件， 再组合使用
  - 组件可以被反复使用
- 3 一次学习，随处编写
  - 不仅可以开发 web 应用（react-dom），还可以开发原生安卓或ios应用（react-native）

## React基本使用

### 基本使用步骤

![image-20220515141459622](react%E5%9F%BA%E7%A1%80.assets/image-20220515141459622.png)

1. 引入两个JS文件（ 注意引入顺序 ）

   ```javascript
   <!-- react库, 提供React对象 -->
   <script src="../js/react.development.js"></script>
   <!-- react-dom库, 提供了ReactDOM对象 -->
   <script src="../js/react-dom.development.js"></script>
   ```

2. 在html定义一个根容器标签 

   ```html
   <div id="root"></div>
   ```

3. 创建react元素(类似html元素)

   ```js
   // 返回值：React元素 
   // 参数1：要创建的React元素名称 =》字符串
   // 参数2：元素的属性  =》对象 {id: 'box'} 或者 null
   // 后面参数：该React元素的所有子节点 =》文本或者其他react元素
   const element = React.createElement(
     'h1', 
     {title: '你好, React!'}, 
     'Hello React!'
   )
   ```

4. 渲染 react 元素

   ```js
   // 渲染React元素到页面容器div中
   ReactDOM.render(element, document.getElementById('root'))
   ```

### 理解 React 元素

1. 也称`虚拟 DOM` (virtual DOM) 或`虚拟节点`(virtual Node)

2. 它就是一个普通的 JS 对象, 它不是真实 DOM 元素对象

   虚拟 DOM: 属性比较少  ==> `较轻`的对象

   真实 DOM: 属性特别多  ==> `较重`的对象

3. 但它有一些自己的特点

   虚拟 DOM 可以转换为对应的真实 DOM  => ReactDOM.render方法将虚拟DOM转换为真实DOM再插入页面

   虚拟 DOM 对象包含了对应的真实 DOM 的关键信息属性

   ​    标签名 => type: "h1"

   ​    标签属性 => props: {title: '你好, React!'}

   ​    子节点 => props: {children: 'Hello React!'}

## JSX

### 基本理解和使用

>`问题`:  React.createElement()写起来太复杂了
>
>`解决`:  推荐使用更加简洁的**JSX**
>
>**JSX** 是一种JS 的扩展语法, 用来快速创建 React 元素(虚拟DOM/虚拟节点)
>
>形式上像HTML标签/任意其它标签, 且标签内部是可以套JS代码的

```jsx
const h1 = <h1 className="active">哈哈哈</h1>   
```

> 浏览器并不认识 JSX 所以需要引入babel将jsx 编译成React.createElement的形式
>
> babel编译 JSX 语法的包为：@babel/preset-react 
>
> 运行时编译可以直接使用babel的完整包：babel.js
>
> 线上测试: https://www.babeljs.cn/

```html	
<!-- 必须引入编译jsx的babel库 -->
<script src="../js/babel.min.js"></script>

<!-- 必须声明type为text/babel, 告诉babel对内部的代码进行jsx的编译 -->
<script type="text/babel">
	// 创建React元素 (也称为虚拟DOM 或 虚拟节点)
	const vDom = <h1 title="你好, React2!" className="active">Hello React2!</h1>
    // 渲染React元素到页面容器div中
    ReactDOM.render(vDom, document.getElementById('root'))
</script>
```

> 注意:
>
> ​	必须有结束标签
> ​	整个只能有一个根标签
> ​	空标签可以自闭合

### JSX中使用 JS 表达式

- JSX中使用JS表达式的语法：`{js表达式}`
- 作用: `指定动态的属性值和标签体文本`

> 1. 可以是js的表达式, 不能是js的语句
>
> 2. 可以是任意基本类型数据值, 但null、undefined和布尔值没有任何显示
>
> 3. 可以是一个js数组, 但不能是js对象
>
> 4. 可以是react元素对象
>
> 5. style属性值必须是一个包含样式的js对象

```jsx
let title = 'I Like You'
const vNode = (
  <div>
    <h3 name={title}>{title.toUpperCase()}</h3>
    <h3>{3}</h3>
    <h3>{null}</h3>
    <h3>{undefined}</h3>
    <h3>{true}</h3>
    <h3>{'true'}</h3>
    <h3>{React.createElement('div', null, 'atguigu')}</h3>
    <h3>{[1, 'abc', 3]}</h3>
    <h3 title={title} id="name" className="ative" style={{color: 'red'}}></h3>
    {/* <h3>{{a: 1}}</h3> */} 
  </div>
)
```

## 样式处理

### 行内样式

- 样式属性名使用小驼峰命名法
- 如果样式是数值，可以省略单位

```js
<h2 style={{color: 'red', fontSize: 30}}>React style</h2>
```

### 类名

- 必须用className, 不能用class
- 推荐, 效率更高些

```js
<h2 className="title">React class</h2>
```

## 事件处理

### 绑定事件

React 元素的事件处理和 DOM 元素的很相似，但是有一点语法上的不同：

- React 事件的命名采用小驼峰式（camelCase），而不是纯小写。比如：onClick、onFocus 、onMouseEnter
- 使用 JSX 语法时你需要传入一个函数作为事件处理函数，而不是一个字符串

```javascript
const div = <div onClick={事件处理函数}></div>
```

### 事件对象 

React 根据 [W3C 规范](https://www.w3.org/TR/DOM-Level-3-Events/)来自定义的合成事件, 与原生事件不完全相同

- 处理好了浏览器的兼容性问题

- 阻止事件默认行为不能使用return false,  必须要调用: event.preventDefault()

- 有自己特有的属性, 比如: nativeEvent --原生事件对象

- `<input>`的change监听在输入过程中触发， 而原生是在失去焦点才触发
  - 原理：内部绑定的是原生input事件

### 合成事件与原生事件的 3 大核心区别

我们可以通过一个对比表格直观地看出它们的差异：

| **特性**         | **原生 DOM 事件**              | **React 合成事件**                   |
| ---------------- | ------------------------------ | ------------------------------------ |
| **命名方式**     | 全小写（如 onclick, onchange） | 小驼峰（如 onClick, onChange）       |
| **事件处理函数** | 接收字符串（"handleClick()"）  | 接收函数引用（{handleClick}）        |
| **绑定位置**     | 真实的 DOM 节点上              | 应用的根节点（如 #root 或 document） |

> **核心避坑指南：**
>
> 在原生 JS 中，如果你想阻止表单默认提交，可以写 return false;
>
> 但在 React 中，**return false 没有任何效果**，你必须显式地调用 e.preventDefault()。

随着 React 版本的迭代，合成事件的底层机制发生过重大的改变。

1. 绑定位置的改变（React 16 vs 17）

- **React 16 及以前：** 所有的事件都被委托绑定在浏览器的 document 节点上。这导致如果页面上同时存在多个 React 实例（或者混用微前端、其他框架），事件冒泡很容易失控。
- **React 17 及以后：** 事件不再绑定到 document，而是绑定到你渲染 React 应用的**根 DOM 容器**（即 ReactDOM.createRoot 挂载的那个 #root 节点）。这使得多版本 React 共存或嵌入其他技术栈时更加安全。

2. 事件池（Event Pooling）的废弃

- **在 React 16 以前**，为了节省内存，React 会重用事件对象。这意味着当你的事件处理器执行完后，e 里面的属性都会被清空（变成 null）。如果你在 setTimeout 异步代码里去拿 e.target.value，直接会报空指针错误，必须手动执行 e.persist() 来脱离对象池。
- **从 React 17 开始**，由于现代浏览器性能大幅提升，**事件池机制被完全废弃**。现在你可以随意在异步回调中访问 e 对象，再也不用写 e.persist() 了。

如何拿到原生事件？

使用`e.nativeEvent`

```jsx
function App() {
  const handleClick = (e) => {
    console.log(e);             // React 的合成事件对象 (SyntheticEvent)
    console.log(e.nativeEvent);  // 浏览器原生的鼠标事件对象 (MouseEvent)
    
    // 注意：
    // e.stopPropagation() 阻止的是 React 事件树的冒泡
    // e.nativeEvent.stopImmediatePropagation() 才能真正阻止绑定在根节点上的原生事件冒泡
  };

  return <button onClick={handleClick}>点击我</button>;
}
```

### React 的事件代理流程

简单用一句话总结它的核心：**“你在组件里写的事件监听，并没有绑定在真实的 DOM 节点上，而是全量委托给了 React 的应用根节点（Root Container），由它统一拦截、组装并分发。”**

以现代 React（v17/v18+）为例，核心工作流：从点击到执行的 5 个步骤

假设你点击了页面上的一个 ，React 底层会经历以下五个阶段：

![event delegation](react%E5%9F%BA%E7%A1%80.assets/react_17_event_delegation.png)

**1.事件初始化（应用启动时）：**准备阶段。

在你的 React 应用刚加载、执行 createRoot(root).render() 时，React 就已经把几乎所有的原生事件（如 click, keyup 等）以监听器的形式注册到了 **#root 容器**上。此时，你的组件甚至还没开始渲染。

**2.触发浏览器原生事件：**捕获与目标阶段。

用户点击了真实的 。浏览器开始执行标准的事件流：

1. **捕获阶段**：从 window 一路向下，经过 #root 节点，最终到达 。
2. **目标阶段**：真正触发点击的按钮。

**3.浏览器原生冒泡与拦截：**冒泡至 Root 节点。

事件开始从  向上冒泡。

当它冒泡到 **#root 容器**时，立刻被 React 在第一步就埋伏好的统一监听器（名叫 dispatchEvent 的函数）成功拦截。

**4.合成事件对象与模拟事件流：**核心核心步骤。

React 的核心处理器开始干活。它会做两件事：

1. **收集路径**：通过点击的真实 DOM 节点，反向查找到对应的 **Fiber 节点**。然后顺着 Fiber 树一路向上爬到顶，收集沿途所有绑定了 onClick 的 React 组件（这一步形成了 React 自己的“捕获”与“冒泡”路径）。
2. **合成事件**：把原生的 PointerEvent 包装成 React 统一的 SyntheticEvent。

**5.批量执行 React 事件回调：**执行阶段。

React 拿着刚刚收集到的组件路径和合成事件对象，开始模拟执行。

- 先从上往下执行收集到的 onClickCapture（捕获）。

- 再从下往上执行收集到的 onClick（冒泡）。

  handleClick 函数就是在这一步被最终调用执行的。

**关键机制：React 是怎么通过 DOM 找到组件的？**

既然所有事件都委托在 #root，React 怎么知道刚才被点击的  对应代码里的哪个组件？

> **秘密在于 DOM 节点上的隐藏属性：** React 在渲染真实 DOM 时，会在 DOM 节点上挂载一个以 __reactFiber$ 开头的特殊属性。这个属性直接指向了该 DOM 对应的 **Fiber 节点**（即 React 的内部虚拟节点，里面存着你写的 onClick 回调函数）。
>
> 当事件冒泡到根节点被拦截时，React 只需要读取 e.target.__reactFiber$... 就能瞬间顺藤摸瓜，拿到你写的所有事件函数。

**深度思考：e.stopPropagation() 到底阻止了什么？**

理解了上面的代理流程，你就能秒懂很多以前想不通的冒泡诡异现象：

1. **在 React 的 onClick 里写 e.stopPropagation()：** 它阻止的**不是**浏览器的原生冒泡，因为当这个函数执行时，浏览器的原生事件**已经**冒泡到 #root 节点了。它阻止的是 React 内部模拟的那个“Fiber 树向上遍历”的循环过程。
2. **原生事件与合成事件的冲突：** 如果你在 document 上绑定了一个原生的 click 监听器，同时在 React 组件里写了 onClick 并调用了 e.stopPropagation()。
   - **结果：** document 上的原生监听器**依然会被触发**。
   - **原因：** React v17+ 的事件只代理到了 #root。事件从 #root 离开继续往上冒泡，最终还是会到达 document。

## 创建组件的两种方式

### 函数组件

```javascript
 function App() {
  // return null
  return <div>App</div>
}

// 函数名就是组件名
ReactDom.render(<App />, document.getElementById('root')) 
```

> 1. 组件名首字母必须大写. 因为react以此来区分组件元素/标签 和 一般元素/标签
>
> 2. 组件内部如果有多个标签,必须使用一个根标签包裹.只能有一个根标签
>
> 3. 必须有返回值.返回的内容就是组件呈现的结构, 如果返回值为 null，表示不渲染任何内容
> 4. 会在组件标签渲染时调用, 但不会产生实例对象（this->undefined）,  不能有状态

> 注意: 后面我们会讲如何在函数组件中定义状态 ==> hooks语法

### 类组件

```javascript
import React from "react"

class App extends React.Component {
  render () {
    return <div>App Component</div>
  }
}

ReactDom.render(<App />, document.getElementById('root'))
```

> 1. 组件名首字母必须大写.
>
> 2. 组件内部如果有多个标签,必须使用一个根标签包裹.只能有一个根标签
>
> 3. 类组件应该继承 React.Component 父类，从而可以使用父类中提供的方法或属性
> 4. 类组件中必须要声明一个render函数, reander返回组件代表组件界面的虚拟DOM元素
>
> 5. 会在组件标签渲染时调用, 产生实例对象(this->组件实例对象),  可以有状态

## 类组件的状态 state

函数组件又叫做无状态组件(不产生实例)，类组件又叫做有状态组件(有实例) 

状态（state）即数据 

函数组件没有state, 只能根据外部传入的数据（props）动态渲染

类组件有自己的state数据，**一旦更新state数据， 界面就会自动更新**

### state的基本使用

- 状态（state）即数据，是组件内部的私有数据，只能在组件内部使用 
- 组件对象的state属性
  - 属性值为对象, 可以在state对象中保存多个数据
  - 初始化state
    - 构造器中: this.state = {xxx: 2}
    - 类体中: state = {}
  - 读取state数据
    - this.state.xxx
  - 更新state数据
    - 不能直接更新state数据
    - 必须 this.setState({ 要修改的属性数据 })

```jsx
class StateTest extends React.Component {

  /* constructor () {
    super() // 必须调用super()
    // 初始化state
    this.state = {
      count: 0,
      xxx: 'abc'
    }
  } */
  // 初始化状态(简洁语法)
  state = {
    count: 0,
    xxx: 'abc'
  }


  render () {
    // 读取state数据
    const {count} = this.state

    return <div onClick={() => {
      // 直接更新状态数据 => 界面不会自动更新  不可用
      // this.state.count = count + 1
      
      // 通过setState()更新state => 界面会自动更新
      this.setState({
        count: count + 1
      })
    }}>点击的次数: {count}</div>
  }
}
```

## 组件的props

### 使用

组件是封闭的，要接收外部数据应该通过 props 来实现 

 props的作用：父组件向子组件传递数据

父向子传入数据：给组件标签添加属性 

子读取父传入的数据：<font color='cornflowerblue'>函数组件</font>通过<font color='red'>参数props</font>接收数据，<font color='cornflowerblue'>类组件</font>通过 <font color='red'>this.props 接收数据 </font>

props的特点

1. 可以给组件传递任意类型的数据 
2. props 是只读的对象，只能读取属性的值，不要修改props 
3. 可以通过...运算符来将对象的多个属性分别传入子组件
4. 如果父组件传入的是动态的state数据, 那一旦父组件更新state数据, 子组件也会更新

- 子组件

  ```js
  // 函数组件
  export function FunProps(props) {
    return <h2>FunProps-个人信息: 姓名: {props.name}, 年龄: {props.age}</h2>
  }
  
  // 类组件
  export class ClassProps extends React.Component {
    render () {
      const { myName, age} = this.props
      return <h2>ClassProps-个人信息: 姓名: {myName}, 年龄: {age}</h2>
    }
  }
  ```

- 父组件【状态数据传递给子组件】

  ```js
  class App extends React.Component {
    state = {
      person: {
        myName: 'tom',
        age: 12
      }
    }
  
    render () {
      const {myName, age} = this.state.person
      return <div>
          <p>人员信息: {myName + ' : ' +age}</p>
          <button onClick={() => {
            this.setState({
              person: { myName: myName+'--', age: age+1}
            })
          }}>更新人员信息</button>
          <br/>
  
          <FunProps name={myName} age={age}/>
          <hr/>
          {/* <ClassProps myName={myName} age={age}/> */}
          <ClassProps {...this.state.person}/>
        </div>
    }
  }
  ```

## 类组件的生命周期

生命周期图谱: https://projects.wojtekmaj.pl/react-lifecycle-methods-diagram/

![react生命周期](react%E5%9F%BA%E7%A1%80.assets/react%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F.png)

### 生命周期三大阶段

#### 挂载阶段

> 流程: constructor  ==> render ==> componentDidMount
>
> 触发: ReactDOM.render(): 渲染组件元素

#### 更新阶段

>流程: render  ==>  componentDidUpdate 
>
>触发: setState() , forceUpdate(), 组件接收到新的props

#### 卸载阶段

>流程: componentWillUnmount
>
>触发: 不再渲染组件

### 生命周期钩子

- constructor: 

  只执行一次: 创建组件对象挂载第一个调用

  用于初始化state属性或其它的实例属性或方法(可以简写到类体中)

- render:

  执行多次: 挂载一次 + 每次state/props更新都会调用

  用于返回要初始显示或更新显示的虚拟DOM界面

- componentDidMount:

  执行一次: 在第一次调用render且组件界面已显示之后调用

  用于初始执行一个异步操作: 发ajax请求/启动定时器等

  应用：

  1. 启动定时器
  2. 订阅消息
  3. 发送ajax请求

- componentDidUpdate:

  执行多次: 组件界面更新(真实DOM更新)之后调用

  用于数据变化后, 就会要自动做一些相关的工作(比如: 存储数据/发请求)

  用得少  => 这次我们先简单了解, 后面需要时再深入说

- componentWillUnmount:

  执行一次: 在组件卸载前调用

  用于做一些收尾工作, 如: 清除定时器、取消订阅

## Hooks

- *Hook* 是 React 16.8 的新增特性。它可以让你在不编写 class 的情况下使用 state 以及其他的 React 特性
- Hook也叫钩子，**本质就是函数**，能让你使用 React 组件的状态和生命周期函数... 
- Hook 语法 基本已经代替了类组件的语法
- 后面的 React 项目就完全是用Hook语法了

### Hook规则: 

1. 只在React组件函数内部中调用 Hook, 不要在组件函数外部调用
2. Hook调用的次数要固定, 所以不要在循环或条件判断中调用

### 官方 Hooks

- [基础 Hook](https://react.docschina.org/docs/hooks-reference.html#basic-hooks)
  - [`useState`](https://react.docschina.org/docs/hooks-reference.html#usestate)
  - [`useEffect`](https://react.docschina.org/docs/hooks-reference.html#useeffect)
  - [`useContext`](https://react.docschina.org/docs/hooks-reference.html#usecontext)
- [额外的 Hook](https://react.docschina.org/docs/hooks-reference.html#additional-hooks)
  - [`useReducer`](https://react.docschina.org/docs/hooks-reference.html#usereducer)
  - [`useCallback`](https://react.docschina.org/docs/hooks-reference.html#usecallback)
  - [`useMemo`](https://react.docschina.org/docs/hooks-reference.html#usememo)
  - [`useRef`](https://react.docschina.org/docs/hooks-reference.html#useref)
  - [`useImperativeHandle`](https://react.docschina.org/docs/hooks-reference.html#useimperativehandle)
  - [`useLayoutEffect`](https://react.docschina.org/docs/hooks-reference.html#uselayouteffect)
  - [`useDebugValue`](https://react.docschina.org/docs/hooks-reference.html#usedebugvalue)

## 收集表单数据

**受控组件：** 数据的“管家”是 **React 状态（State）**。

**非受控组件：** 数据的“管家”是 **DOM 元素自己**。

### 非受控组件

表单项不与state数据相向关联, 需要手动读取表单元素的值

借助于 useRef，使用原生 DOM 方式来获取表单元素值 

useRef 的作用：用于获取 DOM元素

```jsx
<form>
  <h2>登陆页面</h2>
  用户名: <input type="text"/> <br/>
  密  码: <input type="password"/> <br/>
  <input type="submit" value="登 陆"/>
</form>
```

```jsx
import React, { useRef } from 'react'

/* 
非受控组件:
  包含表单组件
  在输入过程中, 不将输入数据收集到state数据中, 只是提交的回调中手动读取input中的输入值
  表单项输入数据不与state数据相关联
编码过程
  1. 使用useRef创建用于存储input元素的容器对象(内部使用current属性存储)
  2. 将ref容器通过ref属性交给表单项标签 => 渲染时内部会将对应的input元素保存到ref容器的current属性上
  3. 点击提交按钮时, 通过ref容器的current属性得到input DOM元素 => 就可以读取其value了
不足:
  不够自动 / 不方便进行实时的数据检验
*/
export default function FormTest () {

  const nameRef = useRef()
  const pwdRef = useRef()
  console.log(nameRef) // {current: undefined}

  // 点击登陆的回调
  const login = (event) => {
    console.log(nameRef)

    // 阻止事件的默认行为 => 不提交表单
    event.preventDefault()
    // 得到输入框
    const nameInput = nameRef.current
    const pwdInput = pwdRef.current

    // 得到输入框的值
    const name = nameInput.value
    const pwd = pwdInput.value

    // 发送登陆的请求
    alert(`发送登陆的请求 name=${name}, pwd=${pwd}`)
  }

  return (
    <form>
      <h2>登陆页面(非受控组件)</h2>
      用户名: <input ref={nameRef}  type="text"/> <br/>
      密  码: <input ref={pwdRef} type="password"/> <br/>
      <input type="submit" value="登 陆" onClick={login}/>
    </form>
  )
}
```

### 受控组件

组件中的表单项根据状态数据动态初始显示和更新显示, 当用户输入时实时同步到状态数据中

也就是实现了页面表单项与state数据的双向绑定

**实现方式**

1. 在 state 中添加一个状态，作为表单元素的value值（控制表单元素值的来源） 
2. 给表单元素绑定 change 事件，将 表单元素的值 设置为 state 的值（控制表单元素值的变化） 

```jsx
import React, { useState } from 'react'

/* 
受控组件
  在输入过程, 实时收集到state数据中 / 界面也可以根据state数据进行显示
  表单项与state数据进行 双向同步 => 数据双向绑定  state <===> 页面的input
编码过程
  1. 使用useState定义一个state数据，作为表单元素的value值（界面根据state动态显示） 
  2. 给表单元素绑定 change 事件，将 表单元素的值 设置为 state 的值（界面输入变化时, 保存到state)
数据双向绑定
  state 到 页面 的绑定 => 将state数据指定为input的value
  页面 到 state 的绑定 => 给input绑定change事件, 在回调中将输入的最新值更新到state
好处:
  实时自动收集数据 => 需要数据时非常轻松
  方便进行实时的数据检验
*/
export default function FormTest2 () {
  // 定义state
  const [name, setName] = useState('admin')
  const [pwd, setPwd] = useState('123')

  const handleSubmit = (e) => {
    // 点击提交按钮的默认行为就是提交表单, 但不想自动提交表单 => 阻止一下事件的默认行为
    e.preventDefault()
    alert(`准备提交登陆的ajax请求 name=${name}, pwd=${pwd}`)
  }

  // 当用户名输入发生改变的回调
  const handleNameChange = (e) => {

    // 将最新输入的值更新到name状态
    const name = e.target.value
    setName(e.target.value)

    // 对name进行实时检验: 不能超过6位
    if (name.length>6) {
      alert('用户名不能超过6位')
    }
  }
  
  // 当密码输入发生改变的回调
  const handlePwdChange = (event) => {
    // 将最新输入的值更新到pwd状态
    setPwd(event.target.value)
  }

  return (
    <div>

      <h3>登陆页面(受控组件)</h3>
      <form action='/xxx'>
        {/* 2. 给表单元素绑定 change 事件，将 表单元素的值 设置为 state 的值 */}
        用户名: <input type="text" value={name} onChange={handleNameChange}/><br/>
        密码: <input type="text" value={pwd} onChange={handlePwdChange}/><br/>
        <input type="submit" value='登陆' onClick={handleSubmit}/>
      </form>

      <button onClick={() => { // 更新state, 界面会自动更新
        setName(name + '--')
        setPwd(pwd + '--')
      }}>更新状态数据</button>

    </div>
  )
}
```

**优化: 使用同一个事件函数处理多个事件**

> 方式一: 柯里化函数
>
> 方式二: 包裹箭头函数

```js
import React, { useState } from 'react'

/* 
优化: 将2个事件函数优化为1个
方式一: 柯里化函数
方式二: 包裹箭头函数
*/
export default function FormTest3 () {
  // 定义state
  const [name, setName] = useState('admin')
  const [pwd, setPwd] = useState('123')

  const handleSubmit = (e) => {
    e.preventDefault()
    alert(`准备提交登陆的ajax请求 name=${name}, pwd=${pwd}`)
  }

  /* 
  方式一: 使用柯里化函数(也是一个高阶函数)
  */
  const handleChange = (setFn) => {
    return (event) => {
      setFn(event.target.value)
    }
  }

  /* 
  方式二: 包裹箭头函数: 在外部包一个事件回调函数, 我们在其中调用传递特定参数
  */
  const handleChange2 = (event, setFn) => {
    setFn(event.target.value)
  }

  return (
    <div>

      <h3>登陆页面(受控组件)</h3>
      <form action='/xxx'>
        用户名: <input type="text" value={name} onChange={handleChange(setName)}/><br/>
        密码: <input type="text" value={pwd} onChange={handleChange(setPwd)}/><br/>
        <input type="submit" value='登陆' onClick={handleSubmit}/>
      </form>

      <form action='/xxx'>
        用户名: <input type="text" value={name} onChange={event => handleChange2(event, setName)}/><br/>
        密码: <input type="text" value={pwd} onChange={event => handleChange2(event, setPwd)}/><br/>
        <input type="submit" value='登陆' onClick={handleSubmit}/>
      </form>

      <button onClick={() => { // 更新state, 界面会自动更新
        setName(name + '--')
        setPwd(pwd + '--')
      }}>更新状态数据</button>
    </div>
  )
}
```

### 对比

| **特性 / 需求**            | **受控组件 (Controlled)** | **非受控组件 (Uncontrolled)** |
| -------------------------- | ------------------------- | ----------------------------- |
| **数据存储位置**           | React State (内存)        | 真实 DOM 节点                 |
| **更新频率**               | 每次敲击键盘都重新渲染    | 只有在需要提取时才访问        |
| **实时单项验证**           | 🟢 支持（极其简单）        | ❌ 无法优雅实现                |
| **根据输入动态禁用按钮**   | 🟢 支持                    | ❌ 无法实时联动                |
| **文件上传 (type="file")** | ❌ 很难实现                | 🟢 官方推荐做法                |

**绝大多数场景的黄金法则：**

- **优先使用受控组件**：普通的文本框、下拉菜单、单选/复选框，全部用受控组件。这符合 React “数据驱动视图”的核心哲学。
- **少数特殊场景使用非受控组件**：文件上传、快速移植老旧的原生 JS 插件、或者对极致输入性能有追求的超大型表单。

## 组件通讯

>react组件通讯有三种方式.分别是：props, context, pubsub

### props  

>单向数据流: 非函数属性通过标签属性, 由外层组件逐级传递给内层组件
>
>父子间通信
>祖孙间通信
>兄弟间通信

![props传值](react%E5%9F%BA%E7%A1%80.assets/props%E4%BC%A0%E5%80%BC.png)

![image-20220524080728541](react%E5%9F%BA%E7%A1%80.assets/image-20220524080728541.png)

### context  (了解)

>与任意后代直接通信
>
>一般应用中不使用, 但一些插件库内部会使用context封装, 如: react-redux

![context传值](react%E5%9F%BA%E7%A1%80.assets/context%E4%BC%A0%E5%80%BC.png)

- 调用 React. createContext() 创建 context 对象

  ```javascript
  const context = React.createContext() 
  ```

- 在`外部组件`中使用 context 上的 Provider 组件作为父节点, 使用value属性定义要传递的值。Provider【提供者】

  ```javascript
  <context.Provider value={要传递的值}>  提供数据
    <div className="App"> 
      <Child1 /> 
    </div> 
  </context.Provider>
  ```

- 在`任意后代组件`中, 通过 React 的`useContext`读取数据

  ```js
  function Child () {
    const data = useContext(context)
  	return <div>{data}</div>
  }
  ```


### pubsub 发布订阅

>pubsub不是react特有的，是一种技术，可以在任何js项目中使用。vue react
>
>兄弟/任意组件间直接通信
>
>发布订阅机制: publish / subscribe
>pubsub-js是一个用JS编写的库。
>利用订阅发布模式, 当一个组件的状态发生了变化，可以通知其他任意组件更新这些变化

![pubsub传值](react%E5%9F%BA%E7%A1%80.assets/pubsub%E4%BC%A0%E5%80%BC.png)

- 导入

  ```javascript
  import PubSub from "pubsub-js" // 导入的PubSub是一个对象.提供了发布/订阅的功能
  ```

- pubsub-js 提供的方法

  ```javascript
  // 订阅消息
  // 参数一: 消息名
  // 参数二: 用于接收数据的函数
  // token 订阅消息返回的令牌(用于取消订阅)
  const token = PubSub.subscribe('消息名', function (msg, data) {
      console.log( msg, data );
  });
  
  // 发布消息
  // 参数一: 消息名
  // 参数二: 要传递的数据
  PubSub.publish('消息名', 'hello world!');
  
  // 取消指定的订阅
  PubSub.unsubscribe(token);
  
  // 取消某个话题的所有订阅
  PubSub.unsubscribe(消息名);
  
  // 清除所有话题
   PubSub.clearAllSubscriptions()
  /*
  div.addEventListener('click', (event) => {})
  我们点击div => 浏览器自动帮我分发事件: 事件名, 包含事件相关数据的事件对象
  div.removeEventListener('click')
  */
  ```

## Fragment

>doucmentFragment: 是原生DOM中, 内存中可以用来保存多个DOM节点对象的容器
>
>如果将这个fragment添加到页面中, 它本身不会进入页面, 它的所有子节点会进行页面
>
>react组件中只能有一个根组件.
>
>之前使用div包裹的方式会给html结构增加很多无用的层级
>
>为了解决这个问题,可以使用React.Fragment

### 测试DocumentFragment

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>测试DocumentFragment</title>
</head>
<body>
  <div id="test"></div>

  <script>
    const testDiv = document.getElementById('test')

    const h1 = document.createElement('h1')
    h1.innerHTML = '我是标题'
    const p = document.createElement('p')
    p.innerHTML = '我是内容'

    const fragment = document.createDocumentFragment()
    fragment.appendChild(h1)
    fragment.appendChild(p)

    testDiv.appendChild(fragment)

  </script>
</body>
</html>
```

### 不使用React.Fragment

```jsx
function Hello(){
    return (
      // 渲染到页面之后,这个div就是一个多余的
      <div>
        <h1>fragment</h1>
        <p>hello react</p>
      </div>
    ) 
}
```

### 使用React.Fragment

```jsx
function Hello(){
    return (
      // 这样就只会渲染h1和p
      <React.Fragment>
        <h1>fragment</h1>
        <p>hello react</p>
      </React.Fragment>
    ) 
}
```

### 使用简写(无名标签 <>)

```jsx
function Hello(){
    return (
      // 这是React.Fragment的简写形式
      <>
        <h1>fragment</h1>
        <p>hello react</p>
      </>
    ) 
}
```

### DocumentFragment (了解)

> <React.Fragment> 内部就是使用 DocumentFragment 实现的
>
> DocumentFragment 是也是一种 DOM 节点, 它有几个特点
>
> 	1. 它只存在于内存中, 它本身是不会进入页面显示的
> 	2. 它专门用来存放任意多个节点
> 	3. 如果将它添加到页面标签中, 那进入页面的是它的所有子节点

```html
<div id="test"></div>

<script>
  // 得到页面的空div
  const testDiv = document.getElementById('test')

  // 创建h1标签, 并指定内容
  const h1 = document.createElement('h1')
  h1.innerHTML = '我是标题'
  // 创建p标签, 并指定内容
  const p = document.createElement('p')
  p.innerHTML = '我是内容'

  // 创建fragment容器, 将h1和p添加为它的子节点
  const fragment = document.createDocumentFragment()
  fragment.appendChild(h1)
  fragment.appendChild(p)

  // 将fragment添加为页面div的子节点 => 但fragment不会进入页面
  testDiv.appendChild(fragment)
</script>
```

