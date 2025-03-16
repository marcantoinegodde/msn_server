import { useState, useEffect, createContext } from "react";
import { useQuery } from "@tanstack/react-query";

import { queryKeys } from "@/repositories/queryKeys";
import { postMe } from "@/repositories/user/repositories";
import { PostMeResponse } from "@/repositories/user/types";

export interface AuthContext {
  isAuthenticated: boolean;
  user: PostMeResponse | undefined;
  login: () => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContext | null>(null);

function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<PostMeResponse>();
  const isAuthenticated = !!user;

  const meQuery = useQuery({
    queryKey: queryKeys.me,
    queryFn: postMe,
  });

  const login = () => {
    meQuery.refetch();
  };

  const logout = () => {
    setUser(undefined);
  };

  useEffect(() => {
    if (meQuery.isSuccess) {
      setUser(meQuery.data);
    }
  }, [meQuery.isSuccess, meQuery.data]);

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export { AuthProvider, AuthContext };
