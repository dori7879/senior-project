import axios from "axios";
const API_URL = "https://localhost:8080/api/v1/";


const createHomework = ( courseTitle, title, description, files, openDate, closeDate) => {
    return axios
      .post(API_URL + "homework-page", {
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

  
  export default {
    createHomework
  };