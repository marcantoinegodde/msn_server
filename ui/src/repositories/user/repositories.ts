import axios from "axios";
import { PostMeResponse } from "./types";

export const postMe = async () => {
  const response = await axios.get<PostMeResponse>(
    `${import.meta.env.VITE_API_URL}/user/me`,
    { withCredentials: true }
  );
  return response.data;
};
