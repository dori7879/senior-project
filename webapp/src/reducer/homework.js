import { CREATE_HOMEWORK, FETCH_HOMEWORK } from "../actions/types";

const initialState = {homework: null};

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
      default:
        return state;
    }
  
  
}