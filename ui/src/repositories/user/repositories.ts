import axios from "axios";

import { GetAccountResponse } from "@/repositories/user/types";

export const getAccount = async () => {
  const response = await axios.get<GetAccountResponse>(
    `${import.meta.env.VITE_API_URL}/user/account`,
    { withCredentials: true }
  );
  return response.data;
};
