/* eslint-disable import/no-anonymous-default-export */
import {
  LOGIN_SUCCESS,
  LOGOUT,
} from '../actions/types'

const token = JSON.parse(localStorage.getItem('token'))
const role = JSON.parse(localStorage.getItem('role'))

const initialState = token
  ? { isLoggedIn: true, token: token, role: role }
  : { isLoggedIn: false, token: null, role: null }

export default function auth(state = initialState, action) {
  const { type, payload } = action

  switch (type) {
    case LOGIN_SUCCESS:
      return {
        ...state,
        isLoggedIn: true,
        token: payload.token,
      }
    case LOGOUT:
      return {
        ...state,
        isLoggedIn: false,
        token: null,
        role: null
      }
    default:
      return state
  }
}
  