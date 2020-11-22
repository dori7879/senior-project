import axios from "axios";

const register = ( firstName, lastName, email, password) => {
  return axios.post("/api/v1/signup", {
    "first_name": firstName,
    "last_name": lastName,
    email,
    password,
  });
};

const login = (email, password) => {
  return axios
    .post("/api/v1/login", {
      email,
      password,
    })
    .then((response) => {
      if (response.data) {
        
        localStorage.setItem("access_token", JSON.stringify(response.data.access_token));
        localStorage.setItem("refresh_token", JSON.stringify(response.data.refresh_token));
        localStorage.setItem("role", JSON.stringify(response.data.role));
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
