import {
    CREATE_HOMEWORK,
    FETCH_HOMEWORK,
    SUBMIT_HOMEWORK
} from "./types";

import HomeworkService from "../services/homework-service.js";

export const createHomework = (courseTitle, title, description, files, openDate, closeDate, fullName) => (dispatch) => {
    return HomeworkService.createHomework( courseTitle, title, description, files, openDate, closeDate, fullName).then(
      (data) => {
        dispatch({
          type: CREATE_HOMEWORK,
          payload: { homework: data },
        });
        return Promise.resolve();
      }
    );
  };

export const fetchHomework = (randomStr) => (dispatch) => {
  return HomeworkService.fetchHomework(randomStr).then(
    (data) => {
      dispatch({
        type: FETCH_HOMEWORK
      });
      return Promise.resolve();
    }
  );
} 

export const submitHomework = (fullName, grade, comments, answer, submitDate) => (dispatch) => {
  return HomeworkService.submitHomework(fullName, answer, submitDate, grade, comments)
    .then(
    (data) => {
      dispatch({
        type: SUBMIT_HOMEWORK
      });
      return Promise.resolve();
    }
    ) 
}