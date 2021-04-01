import axios from 'axios'
import authHeader from './auth-header'
import { BASE_URL } from './index'

const getProfile = () => {
  return axios.get(BASE_URL + '/api/v1/profile', { headers: authHeader() })
}



export default {
  getProfile,
}