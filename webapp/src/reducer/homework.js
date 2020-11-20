import { CREATE_HOMEWORK } from "../actions/types";

const initialState = {};

export default function (state = initialState, action){

    const { type, payload } = action;
  
    switch (type) {
      case CREATE_HOMEWORK:
        return {
          ...state,
          payload
        };
        
      default:
        return state;
    }
  
  
}