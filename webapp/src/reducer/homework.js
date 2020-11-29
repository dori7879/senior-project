import { CREATE_HOMEWORK, FETCH_HOMEWORK, SUBMIT_HOMEWORK } from "../actions/types";

const teacher_link = JSON.parse(localStorage.getItem("teacher_link"));
const student_link = JSON.parse(localStorage.getItem("student_link"));

const initialState = teacher_link && student_link
  ? { teacher_link, student_link }
  : { teacher_link: null, student_link: null };

export default function (state = initialState, action){

    const { type, payload } = action;
  
    switch (type) {
      case CREATE_HOMEWORK:
        return {
          ...state,
          payload
        };
      case FETCH_HOMEWORK:
        return {
          ...state
        }
      case SUBMIT_HOMEWORK:
        return{
          ...state
        }
      default:
        return state;
    }
  
  
}