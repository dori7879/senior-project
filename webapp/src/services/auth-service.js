import axios from "axios";
const API_URL = "http://localhost:"+ process.env.PORT + "/api/v1/";
console.log(process.env.PORT)
const register = ( firstName, lastName, email, password) => {
  return axios.post(API_URL + "signup", {
    firstName,
    lastName,
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
      console.log(response)
      if (response) {
        
        localStorage.setItem("access_token", JSON.stringify(response.access_token));
        localStorage.setItem("refresh_token", JSON.stringify(response.refresh_token));
        localStorage.setItem("role", JSON.stringify(response.role));
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
