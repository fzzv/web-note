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
