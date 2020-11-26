import axios from "axios";

const createHomework = ( courseTitle, title, description, files, openDate, closeDate) => {
    return axios
      .post("/api/v1/homework-page", {
        course_title: courseTitle, 
        title: title, 
        content: description, 
        opened_at: openDate, 
        closed_At: closeDate
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