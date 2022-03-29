import axios from "axios";

const api = axios.create({
  baseURL: process.env.REACT_APP_BASE_URL,
});

api.interceptors.request.use((config) => {
  const auth = JSON.parse(localStorage.getItem("user"));

  if (auth) {
    config.headers.Authorization = `Bearer ${auth.access_token}`;
  }

  return config;
});


export default api;