import axios from "axios";
import FormData from "form-data";

import { API_ENDPOINT } from "#environment/vars";

function e(name: string) {
  return `${API_ENDPOINT}${name}`;
}

export async function getHome(
  page: number = 1,
  search: string = "",
  order: string = "",
  category: string = ""
) {
  const res = await axios.get(
    e(`home?page=${page}&search=${search}&order=${order}&category=${category}`)
  );
  return res.data;
}

export async function getUser(id: number, page: number) {
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

export async function postArchival(url: string) {
  let form = new FormData();
  form.append("url", url);

  const res = await axios.post(e("archiverequests"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });

  return res.data;
}

export async function deleteVideo(videoId: number) {
  const res = await axios.post(e(`delete/${videoId}`));
  return res.data;
}

export async function postRating(videoID: number, rating: number) {
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

interface IPasswordResetData {
  old_password: string;
  new_password: string;
}

export async function postPasswordReset(data: IPasswordResetData) {
  let form = new FormData();
  form.append("old_password", data.old_password);
  form.append("new_password", data.new_password);

  const res = await axios.post(e("password-reset"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });
  return res.data;
}

export async function getVideo(videoId: number) {
  const res = await axios.get(e(`videos/${videoId}`));
  return res.data;
}

export async function approveVideo(videoId: number) {
  const res = await axios.post(e(`approve/${videoId}`));
  return res.data;
}

export async function banUser(userID: number) {
  const res = await axios.post(e(`ban/${userID}`));
  return res.data;
}

export async function setUserMod(userID: number) {
  const res = await axios.post(e(`setrank/${userID}/1`));
  return res.data;
}

export async function setUserAdmin(userID: number) {
  const res = await axios.post(e(`setrank/${userID}/2`));
  return res.data;
}

export async function getComments(videoID: number) {
  const res = await axios.get(e(`comments/${videoID}`));
  return res.data;
}

export async function getAudits(userID: number, page: number) {
  const res = await axios.get(e(`auditevents/${userID}?page=${page}`));
  return res.data;
}

export async function getDownloadsInProgress() {
  const res = await axios.get(e(`downloadsinprogress`));
  return res.data;
}

interface IPostCommentData {
  video_id: number;
  content: string;
  parent: string;
}

export async function postComment(data: IPostCommentData) {
  let form = new FormData();
  form.append("video_id", data.video_id);
  form.append("content", data.content);
  form.append("parent", data.parent);

  const res = await axios.post(e("comments/"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });
  return res.data;
}

export async function deleteArchivalRequest(download_id: number) {
  let form = new FormData();
  form.append("download_id", download_id);

  const res = await axios.post(e("delete-archiverequest"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });
  return res.data;
}

export async function retryArchivalRequest(download_id: number) {
  let form = new FormData();
  form.append("download_id", download_id);

  const res = await axios.post(e("retry-archiverequest"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });
  return res.data;
}

export async function upvoteComment(
  commentID: number,
  upvoted_already: boolean
) {
  let form = new FormData();
  form.append("comment_id", commentID);
  form.append("user_has_upvoted", upvoted_already);

  const res = await axios.post(e("comment_upvotes/"), form, {
    headers: {
      "content-type": "multipart/form-data",
    },
  });

  return res.data;
}
