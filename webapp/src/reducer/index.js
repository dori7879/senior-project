import auth from "./auth";
import { combineReducers } from "redux";
import homework from "./homework";
import message from "./message";

export default combineReducers({
  auth,
  message,
  homework
});
