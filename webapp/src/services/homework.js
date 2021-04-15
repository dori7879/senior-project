import axios from 'axios'
import authHeader from './auth-header'
import { BASE_URL } from './index'


const createHomework = (data) => {
  return axios.post(BASE_URL + '/api/v1/homeworks', data, { headers: authHeader() })
}

const fetchStudentHomework = (randomStr) => {
  return axios.get(BASE_URL + '/api/v1/homeworks/shared/' + randomStr + '/student', { headers: authHeader() })
}

const fetchTeacherHomework = (randomStr) => {
  return axios.get(BASE_URL + '/api/v1/homeworks/shared/' + randomStr + '/teacher', { headers: authHeader() })
}

const submitHomework = (data) => {
  return axios.post(BASE_URL + '/api/v1/homeworks/' + data.HomeworkID + '/submissions', data, { headers: authHeader() })
}

const fetchHWSubmission = (id) => {
  return axios.get(BASE_URL + '/api/v1/homeworks/submissions/' + id, { headers: authHeader() })
}

const updateHWSubmission = (data) => {
  console.log(data)
  return axios.patch(BASE_URL + '/api/v1/homeworks/submissions/' + data.SubID, data, { headers: authHeader() })
}

const HomeworkService = {
  createHomework,
  fetchStudentHomework,
  submitHomework,
  fetchTeacherHomework,
  fetchHWSubmission,
  updateHWSubmission,
}

export default HomeworkService