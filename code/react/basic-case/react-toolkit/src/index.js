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
