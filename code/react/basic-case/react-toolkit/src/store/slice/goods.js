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
