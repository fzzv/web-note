// 1. 引入模块
import {
  createSlice, // 通过该函数创建数据状态管理模块
  configureStore // 生成仓库
} from '@reduxjs/toolkit'

//2.  创建切片：createSlice 函数返回的是一个对象
const countSlice = createSlice({
  // 模块名字
  name:'counter',
  // 初始状态值
  initialState:{
    count:1
  },
  // 提供方法用于操作状态
  reducers:{
    // 会在countSlice对象的actions属性上增加一个方法 increment
    // state:
    // 如果是第一次执行，那么state值为你的 initialState
    // 如果非第一次执行，那么state值为上一次操作的结果
    // action 是通过dispatch 执行时传递过来的action
    increment(state,action){
      // 在这里可以直接改
      // 1. 不需要对状态进行复制
      // 2. 不需要return
      state.count += action.payload
    }
  }
})

//3. 创建仓库
const store = configureStore({
  reducer:{
    counter:countSlice.reducer
  }
})
console.log(countSlice,store)

// 唯一正确的获取状态数据的途径，getState()返回一个对象，对象的结构与 configureStore传入的参数对象结构相同
// 状态与仓库的关系： state = store.getState()
console.log(store.getState()) // {counter:{count:1}}
console.log(store.getState().counter) // 获得counter模块的状态

//4.  获取状态
const state = store.getState()

// 5. 更改状态
// 错误的方式：不允许直接修改
// state.counter.count = 100;

// 正确的方式：store.dispatch()

// 获取操作counter仓库的 actionCreater 函数 increment. increment(值) 执行的结果会返回一个action {type:'片名/reducer中操作函数名',payload:值}
const {increment} = countSlice.actions
console.log(increment(2)) // {type:'counter/increment',payload:2}
store.dispatch(increment(5))
console.log(store.getState())
