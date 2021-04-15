import axios from 'axios'
import authHeader from './auth-header'
import { BASE_URL } from './index'

const getGroup = (id) => {
  return axios.get(BASE_URL + '/api/v1/groups/' + id.toString(), { headers: authHeader() })
}

const acceptGroupShare = (shareHash) => {
  return axios.post(BASE_URL + '/api/v1/groups/' + shareHash + "/accept", {}, { headers: authHeader() })
}

const createGroup = (title) => {
  return axios.post(BASE_URL + '/api/v1/groups', { Title: title}, { headers: authHeader() })
}

const GroupService = {
  getGroup,
  acceptGroupShare,
  createGroup,
}

export default GroupService