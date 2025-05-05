import { createFileRoute } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";

import {
  postWebauthnRegisterBegin,
  postWebauthnRegisterFinish,
} from "@/repositories/auth/repositories";

export const Route = createFileRoute("/_auth/_layout/passkeys")({
  loader: async () => {
    if (
      PublicKeyCredential &&
      PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable
    ) {
      try {
        return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
      } catch (e) {
        console.error(e);
      }
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  const isUVPAA = Route.useLoaderData();

  const passkeyRegisterBeginMutation = useMutation({
    mutationFn: postWebauthnRegisterBegin,
    onSuccess: async (data) => {
      const options = PublicKeyCredential.parseCreationOptionsFromJSON(data);
      const credential = await navigator.credentials.create({
        publicKey: options,
      });
      if (!(credential instanceof PublicKeyCredential)) {
        throw new TypeError();
      }
      passkeyRegisterFinishMutation.mutate(credential);
    },
  });

  const passkeyRegisterFinishMutation = useMutation({
    mutationFn: (data: PublicKeyCredential) => postWebauthnRegisterFinish(data),
  });

  return (
    <div className="flex flex-col gap-4">
      <div>
        <span>Passkeys Support: </span>
        <span className="font-bold">
          {isUVPAA ? "Available" : "Not available"}
        </span>
      </div>

      <button
        className="cursor-pointer disabled:cursor-not-allowed"
        disabled={!isUVPAA}
        onClick={() => passkeyRegisterBeginMutation.mutate()}
      >
        Create a passkey
      </button>
    </div>
  );
}
