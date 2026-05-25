# React Toolkit

## 使用

安装

```bash
npm install @reduxjs/toolkit
```

index.js

```js
/**
 * 
 * redux toolkit 集中式状态管理
 * 
 */
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
store.dispatch(increment(1))
console.log(store.getState())
```

## RTK 的四大核心杀手锏：

- **1. 告别模板代码（createSlice）：** 它把 Action 和 Reducer 合二为一了。你只需要创建一个 “Slice（切片）”，RTK 会自动在底层帮你生成对应的 Action Creator 和 Reducer。
- **2. 允许你“直接修改”状态（不可变数据的终结者）：** 在原生 JS 和老 Redux 中，修改对象必须用解构赋值（如 return { ...state, count: state.count + 1 }），层级深了极其痛苦。RTK 底层内置了 **Immer.js**，让你能直接写 state.count++ 这种“大逆不道”的代码，而 Immer 会在底层自动帮你转成不可变数据，安全又直观。
- **3. 内置网络请求神器（RTK Query）：** 这是 RTK 极其强大的地方。它自带了一个数据缓存和请求库（类似于 React Query），你只要定义好接口，它会自动帮你生成 React Hooks（如 useGetUsersQuery），自带缓存、加载状态、自动刷新功能，连 useEffect 都不用写了。
- **4. 强大的生态与配置：** 自带开箱即用的 Redux DevTools 浏览器插件支持，调试状态变化如同看电影回放，极其爽快；内置了 redux-thunk 处理异步处理，不需要再额外配置。

## RTK 唯一的缺点：

虽然被极大地简化了，但 Redux 提倡的 **“单向数据流”** 核心概念（Store、Dispatch、Action、Reducer）依然存在。对于小型项目或初学者来说，它的心智负担和概念理解成本依然有点偏高。

## 小案例

目录结构

```js
src
 |--store                     仓库
 |    |- slice                模块
 |    |    |- car.js          购物车模块
 |    |    |- goods.js        商品模块
 |    |- index.js             store
 |- index.js                  测试代码
```

- src->store->index.js

```js
import { configureStore } from "@reduxjs/toolkit"
import goods from './slice/goods.js'
import car from './slice/car.js'
const store = configureStore({
  reducer: {
    goods,
    car
  }
})
export default store
```

- src->store->slice->car.js

```js
import {
  createSlice
} from '@reduxjs/toolkit'

// 添加购物车模块
const carSlice = createSlice({
  name:'car',
  initialState:{
    totalPrice:0,
    carList:[]
  },
  reducers:{
    addCar(state, {payload}){
      const info = state.carList.find(v=>{
        return v.id === payload.id
      })
      if(info){
        info.buyNum += 1
      }else{
        state.carList.unshift({
          ...payload,
          buyNum:1
        })
      }
      state.totalPrice =  state.carList.reduce((pre,cur)=>(pre + cur.buyNum * cur.price),0)
    }
  }
})
export const {addCar} = carSlice.actions
export default carSlice.reducer
```

- src->store->slice->goods.js

```js
import {
  createSlice
} from '@reduxjs/toolkit'
const goodsSlice = createSlice({
  name:'goods',
  initialState:{
    goodsList:[]
  },
  reducers:{
    addGoods(state,{payload}){
      state.goodsList.unshift({
        id:Math.random().toString(30).slice(2),
        ...payload
      })
    }
  }
})
export const {addGoods} = goodsSlice.actions
export default goodsSlice.reducer
```

- src->index.js

```js
import store from './store/index.js'
import { addGoods } from './store/slice/goods.js'
import { addCar } from './store/slice/car.js'

const unsubscribe = store.subscribe(()=>{
  console.log(store.getState())
})
store.dispatch(addGoods({
  name:'华为手机',
  price:1999
}))

// unsubscribe()
store.dispatch(addGoods({
  name:'华为电脑',
  price:9990
}))

store.dispatch(addCar(store.getState().goods.goodsList[0]))
store.dispatch(addCar(store.getState().goods.goodsList[0]))
store.dispatch(addCar(store.getState().goods.goodsList[0]))
store.dispatch(addCar(store.getState().goods.goodsList[1]))
store.dispatch(addCar(store.getState().goods.goodsList[1]))
```

## 异步Slice

在 Redux Toolkit (RTK) 中，**同步 Slice** 和 **异步 Slice** 的核心区别在于：**状态修改的触发源是来自于“纯客户端的本地操作”，还是来自于“不可预测的外部网络请求”。**

在底层实现上，它们的本质区别在于：

- **同步 Slice**：状态改动是**立等可取**的，使用 reducers 属性。
- **异步 Slice**：状态改动需要经历**等待阶段**（加载中、成功、失败），通常需要配合 createAsyncThunk 并使用 extraReducers 属性。

核心区别对比，通过一张表直观地看清它们的差异：

| **特性**         | **同步 Slice**                      | **异步 Slice**                                              |
| ---------------- | ----------------------------------- | ----------------------------------------------------------- |
| **定义属性**     | 写在 reducers 对象里                | 写在 extraReducers 方法里                                   |
| **Action 生成**  | RTK 会自动帮你生成同名的 Action     | 需要手动通过 createAsyncThunk 创建异步 Action               |
| **状态变化数量** | 1 个 Action 对应 1 个确定的状态改变 | 1 个异步请求对应 **3 个状态**（Pending/Fulfilled/Rejected） |
| **典型场景**     | 清空用户输入、切换主题、加减计数器  | 从后端 API 获取用户列表、提交登录表单                       |

同步 Slice：直接修改状态

```jsx
import { createSlice } from '@reduxjs/toolkit';

const counterSlice = createSlice({
  name: 'counter',
  initialState: { value: 0 },
  // 同步操作写在 reducers 里
  reducers: {
    increment: (state) => { state.value += 1; }, // RTK 会自动生成 counter/increment 动作
    decrement: (state) => { state.value -= 1; }
  }
});

export const { increment, decrement } = counterSlice.actions;
export default counterSlice.reducer;
```

异步 Slice：管理请求的“生命周期”

当我们要从服务器获取数据时，这就变成了异步操作。网络请求是有风险和时间的（可能成功、可能失败、可能还在加载）。

因此，异步写法需要两步配合：

1. **在 Slice 外面**用 createAsyncThunk 专门负责发请求。
2. **在 Slice 里面**用 extraReducers 像雷达一样，去监听这个请求触发的 **3 个生命周期状态（Pending / Fulfilled / Rejected）**。

```jsx
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

// ==================== 1. 定义异步 Action (Thunk) ====================
// createAsyncThunk 会自动生成三种 Action 类型：
// fetchUser.pending (开始请求)、fetchUser.fulfilled (成功)、fetchUser.rejected (失败)
export const fetchUserById = createAsyncThunk(
  'user/fetchById',
  async (userId) => {
    const response = await fetch(`https://api.example.com/user/${userId}`);
    return await response.json(); // 这个返回值会变成 fulfilled 状态下的 action.payload
  }
);

// ==================== 2. 创建 Slice ====================
const userSlice = createSlice({
  name: 'user',
  initialState: {
    userInfo: null,
    loading: false,     // 异步专属：用来记录请求状态
    error: null,       // 异步专属：用来记录错误信息
    theme: 'light',     // 同步数据
  },
  
  // 【同步 Slice 区域】：处理即时发生的变化
  reducers: {
    // RTK 会自动生成一个名为 toggleTheme 的 action
    toggleTheme: (state) => {
      // 得益于内置的 Immer.js，可以直接修改 state
      state.theme = state.theme === 'light' ? 'dark' : 'light';
    }
  },

  // 【异步 Slice 区域】：监听异步 Action 的整个生命周期
  extraReducers: (builder) => {
    builder
      // 阶段一：请求刚发送，还没拿到结果
      .addCase(fetchUserById.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      // 阶段二：请求成功，拿到了后端数据
      .addCase(fetchUserById.fulfilled, (state, action) => {
        state.loading = false;
        state.userInfo = action.payload; // 后端返回的数据在这里
      })
      // 阶段三：请求失败（如网络崩溃、404）
      .addCase(fetchUserById.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message;
      });
  }
});

export const { toggleTheme } = userSlice.actions;
export default userSlice.reducer;
```

如何在组件内调用？

```jsx
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { toggleTheme, fetchUserById } from './userSlice';

function UserProfile() {
  const dispatch = useDispatch();
  const { theme, userInfo, loading } = useSelector((state) => state.user);

  // 触发同步：直接调用，Store 立刻变色，组件立刻重绘
  const handleThemeClick = () => {
    dispatch(toggleTheme());
  };

  // 触发异步：通常在 useEffect 或点击事件中触发
  useEffect(() => {
    // 派发异步 Thunk，它不会立刻改变 userInfo，而是先触发 pending，等请求完再触发 fulfilled
    dispatch(fetchUserById('123'));
  }, [dispatch]);

  if (loading) return <div>加载中...</div>;

  return (
    <div className={`app-${theme}`}>
      <button onClick={handleThemeClick}>切换主题</button>
      <h1>用户名：{userInfo?.name}</h1>
    </div>
  );
}
```

## 状态库怎么选

现在的 React 状态库世界，主要分成了三个流派。

| **流派**                          | **代表作**              | **核心哲学 / 特点**                                          | **适用场景**                                           |
| --------------------------------- | ----------------------- | ------------------------------------------------------------ | ------------------------------------------------------ |
| **单向数据流（全局唯一 Store）**  | **Redux Toolkit (RTK)** | 状态集中管理，严格的单向流动，改状态必须发 Action。可预测性极强。 | 大型团队、复杂企业级系统、面试刚需。                   |
| **原子流（Atomic State）**        | **Zustand** / Recoil    | 状态切分成一个一个的小原子，想在哪用就直接用，没有 Redux 那么多概念，代码量极少，极其轻量。 | **当前中小型应用的首选。** 个人项目、敏捷开发。        |
| **响应式流（Proxy/Proxy-based）** | **Valtio** / MobX       | 利用代理机制，像写 Vue 一样直接修改变量，视图自动更新。魔法感最强。 | 喜欢 Vue 编程体验、追求极致开发效率的单人/小团队项目。 |
