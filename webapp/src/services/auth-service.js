import axios from "axios";

const register = ( firstName, lastName, email, password) => {
  return axios.post("api/v1/signup", {
    firstName,
    lastName,
    email,
    password,
  });
};

const login = (email, password) => {
  return axios
    .post("api/v1/login", {
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
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("role");
};

export default {
  register,
  login,
  logout
};
