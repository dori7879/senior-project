import axios from "axios";

const createHomework = ( courseTitle, title, description, files, openDate, closeDate, fullName) => {
    return axios
      .post("/api/v1/homework-page", {
        course_title: courseTitle, 
        title: title, 
        content: description, 
        opened_at: openDate.toJSON(), 
        closed_At: closeDate.toJSON(),
        teacher_fullname: fullName
      })
      .then((response) => {
        if (response.data) {
          localStorage.setItem("teacher_link", JSON.stringify(response.data.teacher_link));
          localStorage.setItem("student_link", JSON.stringify(response.data.student_link));
        }
        return response.data;
      })
  }
const fetchHomework = (randomStr ) => {
  return axios
    .get(
      `/api/v1/homework-page/student/${randomStr}`
    )
    .then((response) => {
      return response.data;
    })
}

const submitHomework = (fullName, answer, submitDate, randomStr) => {
  return axios 
    .post(`/api/v1/homework-page/student/${randomStr}`,{
      student_fullname: fullName,
      content: answer,
      submitted_at: submitDate
    })
    .then((response) => {
      return response.data;
    }) 
}

  
  export default {
    createHomework, fetchHomework, submitHomework
  };