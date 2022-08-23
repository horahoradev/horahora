import { apiFetch, formHeader, IAPIFetchOptions } from "#lib/fetch";
import { type IUserRank } from "#lib/account";
import { IProfileData } from "#entities/profile";
import { type IVideo } from "#codegen/schema/001_interfaces";
import { type IVideoDetailed } from "#types/entities";

export interface IHomeData {
  PaginationData: Record<string, unknown>;
  Videos: IVideo[];
}

export async function fetchHome(
  page: number = 1,
  search: string = "",
  order: string = "",
  category: string = ""
) {
  const pathname = "/home";
  const searchParams = new URLSearchParams([
    ["page", String(page)],
    ["search", search],
    ["order", order],
    ["category", category],
  ]);

  const homePage = await apiFetch<IHomeData>({
    pathname,
    searchParams,
  });

  return homePage;
}

export async function fetchProfile(id: number, page: number) {
  const pathname = `/users/${id}`;
  const searchParams = new URLSearchParams([["page", String(page)]]);

  const profile = await apiFetch<IProfileData>({ pathname, searchParams });

  return profile;
}

export interface IPostDetails extends IVideoDetailed {
  RecommendedVideos: IVideo[];
}

export async function getPost(postID: number) {
  const pathname = `/videos/${postID}`;
  const post = await apiFetch<IPostDetails>({ pathname });

  return post;
}

export async function getPostComments(postID: number) {
  const pathname = `/comments/${postID}`;
  const comments = await apiFetch({ pathname });

  return comments;
}

export async function ratePost(postID: number, rating: number) {
  if (postID == 0) {
    return;
    // TODO: throw
    // throw what?
  }
  const pathname = `/rate/${postID}`;
  const formParams = new URLSearchParams([["rating", String(rating)]]);
  const options: IAPIFetchOptions = {
    method: "POST",
    headers: new Headers([formHeader]),
    body: formParams,
  };

  const result = await apiFetch({ pathname, options });

  return result;
}

export async function deletePost(postID: number) {
  const pathname = `/delete/${postID}`;
  const data = await apiFetch({
    pathname,
    options: {
      method: "POST",
    },
  });

  return data;
}

export async function approvePost(postID: number) {
  const pathname = `/approve/${postID}`;
  const data = await apiFetch({
    pathname,
    options: {
      method: "POST",
    },
  });

  return data;
}

export async function banAccount(accountID: number) {
  const pathname = `/ban/${accountID}`;
  const data = await apiFetch({
    pathname,
    options: {
      method: "POST",
    },
  });

  return data;
}

export async function promoteAccountToMod(accountID: number) {
  const pathname = `/setrank/${accountID}/1`;

  const data = await apiFetch({
    pathname,
    options: {
      method: "POST",
    },
  });

  return data;
}

export async function promoteAccountToAdmin(accountID: number) {
  const pathname = `/setrank/${accountID}/2`;

  const data = await apiFetch({
    pathname,
    options: {
      method: "POST",
    },
  });

  return data;
}

export async function getAudits(accountID: number, page: number) {
  const pathname = `/auditevents/${accountID}`;
  const searchParams = new URLSearchParams([["page", String(page)]]);

  const audits = await apiFetch({ pathname, searchParams });

  return audits;
}

export async function upvoteComment(
  commentID: number,
  upvoted_already: boolean
) {
  const pathname = "/comment_upvotes/";
  const formParams = new URLSearchParams([
    ["comment_id", String(commentID)],
    ["user_has_upvoted", String(upvoted_already)],
  ]);

  const data = await apiFetch({
    pathname,
    options: {
      method: "POST",
      headers: new Headers([formHeader]),
      body: formParams,
    },
  });

  return data;
}
