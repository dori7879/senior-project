import axios from "axios";
const API_URL = "https://localhost:8080/api/v1/";

const register = (/*frstName, lastName, */ email, password) => {
  return axios.post(API_URL + "signup", {
    //firstName,
    //lastName,
    email,
    password,
  });
};

const login = (email, password) => {
  return axios
    .post(API_URL + "login", {
      email,
      password,
    })
    .then((response) => {
      if (response.data.token) {
        console.log(response)
        localStorage.setItem("user", JSON.stringify(response.data));
      }
      return response.data;
    });
};


const logout = () => {
  localStorage.removeItem("user");
};

export default {
  register,
  login,
  logout
};
