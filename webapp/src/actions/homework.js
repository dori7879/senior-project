import {
    CREATE_HOMEWORK,
    FETCH_HOMEWORK
} from "./types";

import HomeworkService from "../services/homework-service.js";

export const createHomework = (courseTitle, title, description, files, openDate, closeDate) => (dispatch) => {
    return HomeworkService.createHomework( courseTitle, title, description, files, openDate, closeDate).then(
      (data) => {
        dispatch({
          type: CREATE_HOMEWORK,
          payload: { homework: data },
        });
  
        return Promise.resolve();
      }
    );
  };

export const fetchHomework = (id) => (dispatch) => {
  return HomeworkService.fetchHomework(id).then(
    (data) => {
      dispatch({
        type: FETCH_HOMEWORK
      });

      return Promise.resolve();
    }
  );
} 