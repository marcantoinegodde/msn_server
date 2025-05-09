import axios from "axios";

import {
  GetAccountResponse,
  UpdateAccountBody,
  updateAccountPasswordBody,
} from "@/repositories/account/types";

export const getAccount = async () => {
  const response = await axios.get<GetAccountResponse>(
    `${import.meta.env.VITE_API_URL}/account`,
    { withCredentials: true }
  );
  return response.data;
};

export const updateAccount = async (data: UpdateAccountBody) => {
  const response = await axios.patch<GetAccountResponse>(
    `${import.meta.env.VITE_API_URL}/account`,
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
  await axios.put(`${import.meta.env.VITE_API_URL}/account/password`, data, {
    withCredentials: true,
  });
};
