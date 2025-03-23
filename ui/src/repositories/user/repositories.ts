import axios from "axios";

import { GetMeResponse } from "@/repositories/user/types";

export const getMe = async () => {
  const response = await axios.get<GetMeResponse>(
    `${import.meta.env.VITE_API_URL}/user/me`,
    { withCredentials: true }
  );
  return response.data;
};
