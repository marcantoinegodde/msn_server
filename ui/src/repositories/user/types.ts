export type GetAccountResponse = {
  email: string;
  first_name: string;
  last_name: string;
  country: string;
  state: string;
  city: string;
};

export type UpdateAccountParams = {
  first_name?: string;
  last_name?: string;
  country?: string;
  state?: string;
  city?: string;
};
