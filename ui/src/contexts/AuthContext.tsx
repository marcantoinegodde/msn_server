import { useState, useEffect, createContext } from "react";
import { useQuery } from "@tanstack/react-query";

import { queryKeys } from "@/repositories/queryKeys";
import { getAccount } from "@/repositories/account/repositories";
import { GetAccountResponse } from "@/repositories/account/types";

export interface AuthContext {
  isAuthenticated: boolean;
  user: GetAccountResponse | undefined;
  login: () => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContext | null>(null);

function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<GetAccountResponse>();
  const isAuthenticated = !!user;

  const accountQuery = useQuery({
    queryKey: queryKeys.account,
    queryFn: getAccount,
    retry: false,
  });

  const login = () => {
    accountQuery.refetch();
  };

  const logout = () => {
    accountQuery.refetch();
  };

  useEffect(() => {
    if (accountQuery.isSuccess) {
      setUser(accountQuery.data);
    }
    if (accountQuery.isError) {
      setUser(undefined);
    }
  }, [accountQuery.isSuccess, accountQuery.isError, accountQuery.data]);

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export { AuthProvider, AuthContext };
