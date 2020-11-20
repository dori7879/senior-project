import authHeader from "./auth-header";
import axios from "axios";

const API_URL = " http://reqres.in/api/";

const getAllUsers = () => {
  return axios.get(API_URL + "users");
};
/*
const getUserBoard = () => {
    return axios.get(API_URL + "user", { headers: authHeader() });
  };
  
  const getModeratorBoard = () => {
    return axios.get(API_URL + "mod", { headers: authHeader() });
  };
*/
export default {
  getAllUsers
};
