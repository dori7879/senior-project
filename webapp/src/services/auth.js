/* eslint-disable import/no-anonymous-default-export */
import axios from 'axios'
import { BASE_URL } from './index'

const signup = ({ FirstName, LastName, Email, Password, IsTeacher }) => {
  return axios.post(BASE_URL + '/api/v1/signup', {
    FirstName,
    LastName,
    Email,
    Password,
    IsTeacher: (IsTeacher === 'true'),
  })
}

const login = (email, password) => {
  return axios
    .post(BASE_URL + '/api/v1/login', {
      Email: email,
      Password: password,
    })
    .then((response) => {
      if (response.data) {
        localStorage.setItem(
          'token',
          JSON.stringify(response.data.Token)
        )
        if (response.data.role) {
          localStorage.setItem('role', JSON.stringify(response.data.role))
        }
      }
      return response.data
    })
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('role')
}

export default {
  signup,
  login,
  logout,
}
