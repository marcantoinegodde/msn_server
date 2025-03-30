import axios from "axios";

import {
  GetAccountResponse,
  UpdateAccountBody,
  updateAccountPasswordBody,
} from "@/repositories/user/types";

export const getAccount = async () => {
  const response = await axios.get<GetAccountResponse>(
    `${import.meta.env.VITE_API_URL}/user/account`,
    { withCredentials: true }
  );
  return response.data;
};

export const updateAccount = async (data: UpdateAccountBody) => {
  const response = await axios.patch<GetAccountResponse>(
    `${import.meta.env.VITE_API_URL}/user/account`,
    data,
    {
      withCredentials: true,
    }
  );
  return response.data;
};

export const updateAccountPassword = async (
  data: updateAccountPasswordBody
) => {
  await axios.put(
    `${import.meta.env.VITE_API_URL}/user/account/password`,
    data,
    {
      withCredentials: true,
    }
  );
};
