import { createFileRoute, redirect, useRouter } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";

import keys from "@/icons/keys.png";
import { queryKeys } from "@/repositories/queryKeys";
import { postLogin } from "@/repositories/auth/repositories";
import { PostLoginParams } from "@/repositories/auth/types";
import { useEffect } from "react";
import { useAuth } from "@/hooks/useAuth";

type LoginSearch = {
  redirect: string;
};

export const Route = createFileRoute("/login")({
  validateSearch: (search: Record<string, unknown>): LoginSearch => {
    return {
      redirect: (search.redirect as string) || "",
    };
  },
  component: RouteComponent,
  beforeLoad: ({ context, search }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({ to: search.redirect || "/" });
    }
  },
});

function RouteComponent() {
  const auth = useAuth();
  const router = useRouter();

  const loginMutation = useMutation({
    mutationKey: queryKeys.login,
    mutationFn: ({ email, password }: PostLoginParams) =>
      postLogin({ email, password }),
    onSuccess: () => {
      auth.login();
    },
  });

  useEffect(() => {
    if (auth.isAuthenticated) {
      router.invalidate();
    }
  }, [auth.isAuthenticated, router]);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const email = formData.get("email");
    const password = formData.get("password");
    if (!email || !password) return;
    loginMutation.mutate({
      email: email.toString(),
      password: password.toString(),
    });
  };

  return (
    <div className="window w-[400px]">
      <div className="title-bar">
        <div className="title-bar-text">Log On - MSN Messenger Service</div>
      </div>
      <div className="window-body">
        <form onSubmit={handleSubmit}>
          <div className="flex space-x-2.5">
            <div>
              <img src={keys} />
            </div>
            <div className="w-full">
              <p>Enter your e-mail address and password:</p>
              <fieldset disabled={loginMutation.isPending}>
                <div className="field-row">
                  <label className="text-nowrap">Logon Name:</label>
                  <input
                    type="email"
                    name="email"
                    required
                    className="w-full"
                  />
                </div>
                <div className="field-row">
                  <label>Password:</label>
                  <input
                    type="password"
                    name="password"
                    required
                    className="w-full"
                  />
                </div>
              </fieldset>
            </div>
          </div>
          <div className="flex justify-end mt-2.5">
            <button type="submit">
              {loginMutation.isPending ? "Loading..." : "OK"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
