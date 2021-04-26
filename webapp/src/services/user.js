import axios from 'axios'
import authHeader from './auth-header'
import { BASE_API_URL } from './index'

const getProfile = () => {
  return axios.get(BASE_API_URL + '/api/v1/profile', { headers: authHeader() })
}

const fetchUserSuggestions = (substr, isTeacher) => {
  return axios.get(BASE_API_URL + '/api/v1/users/suggestions?email=' + substr + "&isTeacher=" + isTeacher, { headers: authHeader() })
}

const UserService = {
  getProfile,
  fetchUserSuggestions,
}

export default UserService