import { useId } from "react";
import { createFileRoute } from "@tanstack/react-router";

import { queryKeys } from "@/repositories/queryKeys";
import { getMe } from "@/repositories/user/repositories";
import { useSuspenseQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_auth/_layout/details")({
  loader: (opts) =>
    opts.context.queryClient.ensureQueryData({
      queryKey: queryKeys.me,
      queryFn: getMe,
    }),
  component: RouteComponent,
});

function RouteComponent() {
  const firstNameId = useId();
  const lastNameId = useId();
  const countryId = useId();
  const stateId = useId();
  const cityId = useId();

  const meQuery = useSuspenseQuery({
    queryKey: queryKeys.me,
    queryFn: getMe,
  });

  return (
    <div>
      <form>
        <fieldset>
          <div className="field-row-stacked">
            <label htmlFor={firstNameId}>First Name</label>
            <input
              type="text"
              id={firstNameId}
              defaultValue={meQuery.data.first_name}
            />
          </div>
          <div className="field-row-stacked">
            <label htmlFor={lastNameId}>Last Name</label>
            <input
              type="text"
              id={lastNameId}
              defaultValue={meQuery.data.last_name}
            />
          </div>
          <div className="field-row-stacked">
            <label htmlFor={countryId}>Country</label>
            <select id={countryId}>
              <option>France</option>
              <option>United States</option>
            </select>
          </div>
          <div className="field-row-stacked">
            <label htmlFor={stateId}>State</label>
            <select id={stateId}>
              <option>Texas</option>
              <option>California</option>
            </select>
          </div>
          <div className="field-row-stacked">
            <label htmlFor={cityId}>City</label>
            <input type="text" id={cityId} defaultValue={meQuery.data.city} />
          </div>
        </fieldset>
        <div className="flex justify-end mt-2.5">
          <button type="submit">OK</button>
        </div>
      </form>
    </div>
  );
}
