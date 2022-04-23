import axios from "axios";
import FormData from "form-data";

const API_ENDPOINT = process.env.REACT_APP_API_ENDPOINT;
if (API_ENDPOINT == null)
  throw new Error("No REACT_APP_API_ENDPOINT provided!");

function e(name) {
  return `${API_ENDPOINT}${name}`;
}

export async function getHome(page=1, search="", order="", category="") {
  const res = await axios.get(e(`home?page=${page}&search=${search}&order=${order}&category=${category}`));
  return res.data;
}

export async function getUser(id, page) {
  const path = "users/" + id + "?page=" + page;
  const res = await axios.get(e(path));
  return res.data;
}

export async function getUserdata() {
  const res = await axios.get(e("currentuserprofile/"));
  return res.data;
}

export async function getArchivalSubscriptions() {
  const res = await axios.get(e("archiverequests"));
  return res.data;
}

export async function postArchival(url) {
  let form = new FormData();
  form.append("url", url)

  const res = await axios.post(e("archiverequests"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });

    return res.data;
}

export async function postRegister(data) {
  let form = new FormData();
  form.append("username", data.username);
  form.append("password", data.password);
  form.append("email", data.email);


  const res = await axios.post(e("register"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });
  return res.data;
}

export async function deleteVideo(videoId) {
  const res = await axios.post(e(`delete/${videoId}`));
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

export async function postRating(videoID, rating) {
  if (videoID == 0) {
    return;
    // TODO: throw
  }
  let form = new FormData();
  form.append("rating", rating);

  const res = await axios.post(e(`rate/${videoID}`), form, {
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

export async function getVideo(videoId) {
  const res = await axios.get(e(`videos/${videoId}`));
  return res.data;
}

export async function approveVideo(videoId) {
  const res = await axios.post(e(`approve/${videoId}`));
  return res.data;
}

export async function banUser(userID) {
  const res = await axios.post(e(`ban/${userID}`));
  return res.data;
}
