export type GetAccountResponse = {
  email: string;
  first_name: string;
  last_name: string;
  country: string;
  state: string;
  city: string;
};

export type UpdateAccountBody = {
  first_name?: string;
  last_name?: string;
  country?: string;
  state?: string;
  city?: string;
};

export type updateAccountPasswordBody = {
  old_password: string;
  new_password: string;
};
