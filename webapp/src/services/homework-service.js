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
        return response.data;
      })
  }
const fetchHomework = ( id ) => {
  return axios
    .get(
      "/api/v1/homework/"+id
    )
    .then((response) => {
      return response.data;
    })
}

  
  export default {
    createHomework, fetchHomework
  };