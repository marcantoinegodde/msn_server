import { useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import {
  queryOptions,
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";

import { queryKeys } from "@/repositories/queryKeys";
import {
  deletePasskey,
  getPasskeys,
  postWebauthnRegisterBegin,
  postWebauthnRegisterFinish,
} from "@/repositories/webauthn/repositories";

const checkUVPAA = async () => {
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
};

const getPasskeysQueryOptions = queryOptions({
  queryKey: queryKeys.passkeys,
  queryFn: () => getPasskeys(),
});

export const Route = createFileRoute("/_auth/_layout/passkeys")({
  loader: async ({ context }) => {
    const isUVPAA = await checkUVPAA();
    context.queryClient.ensureQueryData(getPasskeysQueryOptions);
    return isUVPAA;
  },
  component: RouteComponent,
});

function RouteComponent() {
  const [selectedPasskey, setSelectedPasskey] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const isUVPAA = Route.useLoaderData();

  const getPasskeysQuery = useSuspenseQuery(getPasskeysQueryOptions);

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
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.passkeys });
    },
  });

  const passkeyDeleteMutation = useMutation({
    mutationFn: (id: string) => deletePasskey(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.passkeys });
      setSelectedPasskey(null);
    },
  });

  return (
    <div className="flex flex-col gap-4">
      <div>
        <span>Passkeys Support: </span>
        <span className="font-bold">
          {isUVPAA ? "Available" : "Not available"}
        </span>
      </div>
      <fieldset>
        <legend>Registered Passkeys</legend>
        {getPasskeysQuery.data.length > 0 ? (
          getPasskeysQuery.data.map((passkey) => (
            <div key={passkey.id} className="field-row">
              <input
                id={passkey.id}
                type="radio"
                name="passkey"
                value={passkey.id}
                onChange={() => {
                  setSelectedPasskey(passkey.id);
                }}
              />
              <label htmlFor={passkey.id}>{passkey.name}</label>
            </div>
          ))
        ) : (
          <span className="font-bold">No credential found.</span>
        )}
      </fieldset>
      <div className="flex gap-1 justify-end">
        <button
          className="cursor-pointer disabled:cursor-not-allowed"
          disabled={!selectedPasskey}
          onClick={() =>
            selectedPasskey && passkeyDeleteMutation.mutate(selectedPasskey)
          }
        >
          Delete
        </button>
        <button
          className="cursor-pointer disabled:cursor-not-allowed"
          disabled={!isUVPAA}
          onClick={() => passkeyRegisterBeginMutation.mutate()}
        >
          Create a passkey
        </button>
      </div>
    </div>
  );
}
