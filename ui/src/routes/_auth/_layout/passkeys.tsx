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
      PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable &&
      PublicKeyCredential.isConditionalMediationAvailable
    ) {
      try {
        return await Promise.all([
          PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable(),
          PublicKeyCredential.isConditionalMediationAvailable(),
        ]);
      } catch (e) {
        console.error(e);
      }
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  const results = Route.useLoaderData();

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
      {results?.every((r) => r === true) ? "Available" : "Not available"}
      <button
        className="cursor-pointer"
        onClick={() => passkeyRegisterBeginMutation.mutate()}
      >
        Create a passkey
      </button>
    </div>
  );
}
