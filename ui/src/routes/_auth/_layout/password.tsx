import { createFileRoute } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import axios from "axios";

import password from "@/icons/password.png";
import { updateAccountPassword } from "@/repositories/account/repositories";
import { updateAccountPasswordBody } from "@/repositories/account/types";
import { FieldInfo } from "@/components/FieldInfo";

export const Route = createFileRoute("/_auth/_layout/password")({
  component: RouteComponent,
});

const schema = z
  .object({
    old_password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .max(16, "Password must be at most 16 characters")
      .regex(
        /^[a-zA-Z0-9]*$/,
        "Password must only contain letters and numbers"
      ),
    new_password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .max(16, "Password must be at most 16 characters")
      .regex(
        /^[a-zA-Z0-9]*$/,
        "Password must only contain letters and numbers"
      ),
    new_password_confirmation: z.string(),
  })
  .refine((data) => data.new_password === data.new_password_confirmation, {
    message: "Passwords do not match",
    path: ["new_password_confirmation"],
  });

function RouteComponent() {
  const passwordMutation = useMutation({
    mutationFn: (data: updateAccountPasswordBody) =>
      updateAccountPassword(data),
  });

  const form = useForm({
    defaultValues: {
      old_password: "",
      new_password: "",
      new_password_confirmation: "",
    },
    validators: {
      onChange: schema,
    },
    onSubmit: async ({ formApi, value }) => {
      await passwordMutation.mutateAsync(value);
      formApi.reset();
    },
  });

  return (
    <div className="flex flex-row gap-2.5 w-full">
      <div>
        <img src={password} />
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
              name="old_password"
              children={(field) => {
                return (
                  <>
                    <label htmlFor={field.name}>Old Password:</label>
                    <input
                      type="password"
                      autoComplete="current-password"
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
              name="new_password"
              children={(field) => {
                return (
                  <>
                    <label htmlFor={field.name}>New Password:</label>
                    <input
                      type="password"
                      autoComplete="new-password"
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
              name="new_password_confirmation"
              children={(field) => {
                return (
                  <>
                    <label htmlFor={field.name}>Confirm New Password:</label>
                    <input
                      type="password"
                      autoComplete="new-password"
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
        </fieldset>
        <div className="flex justify-between items-center gap-2 mt-2.5">
          <div className="truncate">
            {passwordMutation.isSuccess && (
              <div>Password updated successfully.</div>
            )}
            {passwordMutation.isError &&
              axios.isAxiosError(passwordMutation.error) && (
                <div className="text-red-700">
                  {passwordMutation.error.response?.status === 401
                    ? "Old password is incorrect."
                    : "Something went wrong."}
                </div>
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
