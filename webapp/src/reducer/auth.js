import {
    LOGIN_FAIL,
    LOGIN_SUCCESS,
    LOGOUT,
    REGISTER_FAIL,
    REGISTER_SUCCESS,
} from "../actions/types";

const access_token = JSON.parse(localStorage.getItem("access_token"));
  
  const initialState = access_token
    ? { isLoggedIn: true, access_token }
    : { isLoggedIn: false, access_token: null };
  
  export default function (state = initialState, action) {
    const { type, payload } = action;
  
    switch (type) {
      case REGISTER_SUCCESS:
        return {
          ...state,
          isLoggedIn: false,
        };
      case REGISTER_FAIL:
        return {
          ...state,
          isLoggedIn: false,
        };
      case LOGIN_SUCCESS:
        return {
          ...state,
          isLoggedIn: true,
          access_token: payload.access_token,
        };
      case LOGIN_FAIL:
        return {
          ...state,
          isLoggedIn: false,
          access_token: null,
        };
      case LOGOUT:
        return {
          ...state,
          isLoggedIn: false,
          access_token: null,
        };
      default:
        return state;
    }
  }
  