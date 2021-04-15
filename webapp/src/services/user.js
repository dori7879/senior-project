import axios from 'axios'
import authHeader from './auth-header'
import { BASE_API_URL } from './index'

const getProfile = () => {
  return axios.get(BASE_API_URL + '/api/v1/profile', { headers: authHeader() })
}

const UserService = {
  getProfile,
}

export default UserService