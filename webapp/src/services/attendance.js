import axios from 'axios'
import authHeader from './auth-header'
import { BASE_API_URL } from './index'


const createAttendance = (data) => {
  return axios.post(BASE_API_URL + '/api/v1/attendances', data, { headers: authHeader() })
}

const fetchStudentAttendance = (randomStr) => {
  return axios.get(BASE_API_URL + '/api/v1/attendances/shared/' + randomStr + '/student', { headers: authHeader() })
}

const fetchTeacherAttendance = (randomStr) => {
  return axios.get(BASE_API_URL + '/api/v1/attendances/shared/' + randomStr + '/teacher', { headers: authHeader() })
}

const submitAttendance = (data) => {
  return axios.post(BASE_API_URL + '/api/v1/attendances/' + data.AttendanceID + '/submissions', data, { headers: authHeader() })
}

const fetchHWSubmission = (id) => {
  return axios.get(BASE_API_URL + '/api/v1/attendances/submissions/' + id, { headers: authHeader() })
}

const updateHWSubmission = (data) => {
  console.log(data)
  return axios.patch(BASE_API_URL + '/api/v1/attendances/submissions/' + data.SubID, data, { headers: authHeader() })
}

const AttendanceService = {
  createAttendance,
  fetchStudentAttendance,
  submitAttendance,
  fetchTeacherAttendance,
  fetchHWSubmission,
  updateHWSubmission,
}

export default AttendanceService