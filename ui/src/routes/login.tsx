import { useEffect, useId } from "react";
import { createFileRoute, redirect, useRouter } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

import keys from "@/icons/keys.png";
import { postLogin } from "@/repositories/auth/repositories";
import {
  postWebauthnLoginBegin,
  postWebauthnLoginFinish,
} from "@/repositories/webauthn/repositories";
import { PostLoginBody } from "@/repositories/auth/types";
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
  loader: async () => {
    if (
      PublicKeyCredential &&
      PublicKeyCredential.isConditionalMediationAvailable
    ) {
      try {
        return await PublicKeyCredential.isConditionalMediationAvailable();
      } catch (e) {
        console.error(e);
      }
    }
  },
});

function RouteComponent() {
  const emailId = useId();
  const passwordId = useId();
  const { isAuthenticated, login } = useAuth();
  const { invalidate } = useRouter();

  const isCMA = Route.useLoaderData();

  const { mutate: passkeyLoginBeginMutate } = useMutation({
    mutationFn: postWebauthnLoginBegin,
    onSuccess: async (
      data,
      variables: { abortController: AbortController }
    ) => {
      const options = PublicKeyCredential.parseRequestOptionsFromJSON(data);
      const credential = await navigator.credentials.get({
        publicKey: options,
        mediation: "conditional",
        signal: variables.abortController.signal,
      });
      if (!(credential instanceof PublicKeyCredential)) {
        throw new TypeError();
      }
      passkeyLoginFinishMutation.mutate(credential);
    },
  });

  const passkeyLoginFinishMutation = useMutation({
    mutationFn: (data: PublicKeyCredential) => postWebauthnLoginFinish(data),
    onSuccess: () => {
      login();
    },
  });

  const loginMutation = useMutation({
    mutationFn: ({ email, password }: PostLoginBody) =>
      postLogin({ email, password }),
    onSuccess: () => {
      login();
    },
  });

  useEffect(() => {
    if (isAuthenticated) {
      invalidate();
    }
  }, [isAuthenticated, invalidate]);

  useEffect(() => {
    const abortController = new AbortController();
    if (isCMA) {
      passkeyLoginBeginMutate({ abortController: abortController });
    }
    return () => {
      abortController.abort();
    };
  }, [isCMA, passkeyLoginBeginMutate]);

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
                  <label htmlFor={emailId} className="text-nowrap">
                    Logon Name:
                  </label>
                  <input
                    type="email"
                    autoComplete="username webauthn"
                    id={emailId}
                    name="email"
                    required
                    className="w-full"
                  />
                </div>
                <div className="field-row">
                  <label htmlFor={passwordId}>Password:</label>
                  <input
                    type="password"
                    autoComplete="current-password webauthn"
                    id={passwordId}
                    name="password"
                    required
                    className="w-full"
                  />
                </div>
              </fieldset>
            </div>
          </div>
          <div className="flex justify-between items-center gap-2 mt-2.5">
            <div className="truncate">
              {loginMutation.isError &&
                axios.isAxiosError(loginMutation.error) && (
                  <div className="text-red-700">
                    {loginMutation.error.response?.status === 401
                      ? "Invalid email or password."
                      : "Something went wrong."}
                  </div>
                )}
            </div>
            <button type="submit" className="cursor-pointer">
              {loginMutation.isPending ? "Loading..." : "OK"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
