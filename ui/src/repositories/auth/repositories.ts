import axios from "axios";

import {
  PostLoginParams,
  PostLoginResponse,
  PostLogoutResponse,
} from "@/repositories/auth/types";

export const postLogin = async ({ email, password }: PostLoginParams) => {
  const response = await axios.post<PostLoginResponse>(
    `${import.meta.env.VITE_API_URL}/auth/login`,
    { email, password },
    { withCredentials: true }
  );
  return response.data;
};

export const postLogout = async () => {
  const response = await axios.post<PostLogoutResponse>(
    `${import.meta.env.VITE_API_URL}/auth/logout`,
    {},
    { withCredentials: true }
  );
  return response.data;
};
