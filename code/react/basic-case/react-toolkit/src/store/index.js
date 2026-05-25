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
