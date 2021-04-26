import axios from 'axios'
import authHeader from './auth-header'
import { BASE_API_URL } from './index'


const createQuiz = (data) => {
  return axios.post(BASE_API_URL + '/api/v1/quizzes', data, { headers: authHeader() })
}

const fetchStudentQuiz = (randomStr) => {
  return axios.get(BASE_API_URL + '/api/v1/quizzes/shared/' + randomStr + '/student', { headers: authHeader() })
}

const fetchTeacherQuiz = (randomStr) => {
  return axios.get(BASE_API_URL + '/api/v1/quizzes/shared/' + randomStr + '/teacher', { headers: authHeader() })
}

const submitQuiz = (data) => {
  return axios.post(BASE_API_URL + '/api/v1/quizzes/' + data.QuizID + '/submissions', data, { headers: authHeader() })
}

const fetchQuizSubmission = (id) => {
  return axios.get(BASE_API_URL + '/api/v1/quizzes/submissions/' + id, { headers: authHeader() })
}

const updateQuizSubmission = (data) => {
  return axios.patch(BASE_API_URL + '/api/v1/quizzes/submissions/' + data.SubID, data, { headers: authHeader() })
}

const QuizService = {
  createQuiz,
  fetchStudentQuiz,
  submitQuiz,
  fetchTeacherQuiz,
  fetchQuizSubmission,
  updateQuizSubmission,
}

export default QuizService