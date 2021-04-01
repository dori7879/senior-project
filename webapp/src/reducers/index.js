import auth from './auth'
import message from './message'
import { combineReducers } from 'redux'

export default combineReducers({
  auth,
  message,
})
