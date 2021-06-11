import axios from "axios";
import FormData from "form-data";

const API_ENDPOINT = process.env.REACT_APP_API_ENDPOINT;
if (API_ENDPOINT == null)
  throw new Error("No REACT_APP_API_ENDPOINT provided!");

function e(name) {
  return `${API_ENDPOINT}${name}`;
}

export async function getHome() {
  const res = await axios.get(e("home"));
  return res.data;
}

export async function postLogin(data) {
  let form = new FormData();
  form.append("username", data.username);
  form.append("password", data.password);

  const res = await axios.post(e("login"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });
  return res.data;
}

export async function postLogout() {
  const res = await axios.post(e("logout"));
  return res.data;
}
