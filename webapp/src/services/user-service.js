import authHeader from "./auth-header";
import axios from "axios";


const getAllUsers = () => {
  return axios.get("/api/v1/users");
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
