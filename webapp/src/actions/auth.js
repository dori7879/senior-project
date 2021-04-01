import {
  LOGIN_SUCCESS,
  LOGOUT,
  SET_MESSAGE,
} from './types'

import AuthService from '../services/auth'

export const login = ({ Email, Password }) => (dispatch) => {
  return AuthService.login(Email, Password).then(
    (data) => {
      dispatch({
        type: LOGIN_SUCCESS,
        payload: { token: data.Token },
      })

      return Promise.resolve()
    },
    (error) => {
      const message =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.message ||
        error.toString()

      dispatch({
        type: SET_MESSAGE,
        payload: message,
      })

      return Promise.reject()
    }
  )
}

export const logout = () => (dispatch) => {
  AuthService.logout()

  dispatch({
    type: LOGOUT,
  })
}
  