import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import {
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { z } from "zod";

import userDetails from "@/icons/user_details.png";
import { queryKeys } from "@/repositories/queryKeys";
import { getAccount, updateAccount } from "@/repositories/account/repositories";
import { UpdateAccountBody } from "@/repositories/account/types";
import { FieldInfo } from "@/components/FieldInfo";
import { CountryOptions } from "@/components/CountryOptions";
import { StateOptions } from "@/components/StateOptions";

const schema = z
  .object({
    first_name: z
      .string()
      .min(2, "First name must be at least 2 characters")
      .regex(/^[a-zA-Z' -]*$/, "First name is invalid"),
    last_name: z
      .string()
      .min(2, "Last name must be at least 2 characters")
      .regex(/^[a-zA-Z' -]*$/, "Last name is invalid"),
    country: z.string(),
    state: z.string(),
    city: z.string(),
  })
  .superRefine((val, ctx) => {
    const regex = /^[a-zA-Z' -]*$/;
    if (val.country === "US") {
      if (val.city.length < 2) {
        ctx.addIssue({
          code: z.ZodIssueCode.too_small,
          minimum: 2,
          type: "string",
          inclusive: true,
          message: "City must be at least 2 characters",
          path: ["city"],
        });
      }
      if (!regex.test(val.city)) {
        ctx.addIssue({
          code: z.ZodIssueCode.invalid_string,
          validation: "regex",
          message: "City is invalid",
          path: ["city"],
        });
      }
    }
  });

export const Route = createFileRoute("/_auth/_layout/details")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData({
      queryKey: queryKeys.account,
      queryFn: getAccount,
    }),
  component: RouteComponent,
});

function RouteComponent() {
  const queryClient = useQueryClient();

  const accountQuery = useSuspenseQuery({
    queryKey: queryKeys.account,
    queryFn: getAccount,
  });

  const accountMutation = useMutation({
    mutationFn: (data: UpdateAccountBody) => updateAccount(data),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.account, data);
    },
  });

  const form = useForm({
    defaultValues: {
      first_name: accountQuery.data.first_name,
      last_name: accountQuery.data.last_name,
      country: accountQuery.data.country,
      state: accountQuery.data.state,
      city: accountQuery.data.city,
    },
    validators: {
      onChange: schema,
    },
    onSubmit: async ({ formApi, value }) => {
      await accountMutation.mutateAsync(value);
      formApi.reset();
    },
  });

  return (
    <div className="flex flex-row gap-2.5 w-full">
      <div>
        <img src={userDetails} />
      </div>
      <form
        className="w-full overflow-auto"
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
      >
        <fieldset>
          <div className="field-row-stacked">
            <form.Field
              name="first_name"
              children={(field) => {
                return (
                  <>
                    <label htmlFor={field.name}>First Name:</label>
                    <input
                      type="text"
                      autoComplete="given-name"
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    <FieldInfo field={field} />
                  </>
                );
              }}
            />
          </div>
          <div className="field-row-stacked">
            <form.Field
              name="last_name"
              children={(field) => {
                return (
                  <>
                    <label htmlFor={field.name}>Last Name:</label>
                    <input
                      type="text"
                      autoComplete="family-name"
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    <FieldInfo field={field} />
                  </>
                );
              }}
            />
          </div>
          <div className="field-row-stacked">
            <form.Field
              name="country"
              listeners={{
                onChange: ({ value }) => {
                  if (value === "US") {
                    form.setFieldValue("state", "AL");
                  } else {
                    form.setFieldValue("state", "");
                    form.setFieldValue("city", "");
                  }
                },
              }}
              children={(field) => {
                return (
                  <>
                    <label htmlFor={field.name}>Country/Region:</label>
                    <select
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                    >
                      <CountryOptions />
                    </select>
                  </>
                );
              }}
            />
          </div>
          <form.Subscribe
            selector={(state) => state.values.country === "US"}
            children={(isUS) =>
              isUS && (
                <>
                  <div className="field-row-stacked">
                    <form.Field
                      name="state"
                      children={(field) => {
                        return (
                          <>
                            <label htmlFor={field.name}>State:</label>
                            <select
                              id={field.name}
                              name={field.name}
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) =>
                                field.handleChange(e.target.value)
                              }
                            >
                              <StateOptions />
                            </select>
                          </>
                        );
                      }}
                    />
                  </div>
                  <div className="field-row-stacked">
                    <form.Field
                      name="city"
                      children={(field) => {
                        return (
                          <>
                            <label htmlFor={field.name}>City:</label>
                            <input
                              type="text"
                              autoComplete="address-level2"
                              id={field.name}
                              name={field.name}
                              value={field.state.value}
                              onBlur={field.handleBlur}
                              onChange={(e) =>
                                field.handleChange(e.target.value)
                              }
                            />
                            <FieldInfo field={field} />
                          </>
                        );
                      }}
                    />
                  </div>
                </>
              )
            }
          />
        </fieldset>
        <div className="flex justify-between items-center gap-2 mt-2.5">
          <div className="truncate">
            {accountMutation.isSuccess && (
              <div>Account updated successfully.</div>
            )}
            {accountMutation.isError && (
              <div className="text-red-700">Something went wrong.</div>
            )}
          </div>
          <div className="flex gap-1">
            <form.Subscribe
              selector={(state) => [state.canSubmit, state.isSubmitting]}
              children={([canSubmit, isSubmitting]) => (
                <>
                  <button
                    type="reset"
                    onClick={(e) => {
                      e.preventDefault();
                      form.reset();
                    }}
                    className="cursor-pointer"
                  >
                    Reset
                  </button>
                  <button
                    type="submit"
                    disabled={!canSubmit}
                    className="cursor-pointer"
                  >
                    {isSubmitting ? "Loading..." : "OK"}
                  </button>
                </>
              )}
            />
          </div>
        </div>
      </form>
    </div>
  );
}
