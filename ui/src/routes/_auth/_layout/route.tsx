import {
  createFileRoute,
  Link,
  linkOptions,
  Outlet,
  useLocation,
} from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";

import users from "@/icons/users.png";
import padlock from "@/icons/padlock.png";
import { queryKeys } from "@/repositories/queryKeys";
import { getMe } from "@/repositories/user/repositories";

export const Route = createFileRoute("/_auth/_layout")({
  loader: (opts) => {
    opts.context.queryClient.ensureQueryData({
      queryKey: queryKeys.me,
      queryFn: getMe,
    });
  },
  component: RouteComponent,
});

const options = linkOptions([
  {
    to: "/",
    label: "Home",
  },
  {
    to: "/details",
    label: "User Details",
  },
  {
    to: "/status",
    label: "Server Status",
  },
]);

function RouteComponent() {
  const location = useLocation();
  const { data } = useSuspenseQuery({
    queryKey: queryKeys.me,
    queryFn: getMe,
  });

  return (
    <div className="window w-1/2">
      <div className="title-bar">
        <div className="title-bar-text">
          {options.find((option) => option.to === location.pathname)?.label} -
          MSN Messenger Service
        </div>
      </div>
      <div className="window-body">
        <div className="flex justify-between my-3">
          <fieldset>
            <div className="flex items-center gap-2.5">
              <img height={32} width={32} src={users} />
              <span>{data.email}</span>
            </div>
          </fieldset>
          <button className="cursor-pointer">
            <div className="flex items-center gap-2.5 py-2.5">
              <img height={32} width={32} src={padlock} />
              <span>Logout</span>
            </div>
          </button>
        </div>
        <menu role="tablist">
          {options.map((option, index) => {
            return (
              <li
                key={index}
                role="tab"
                aria-selected={option.to === location.pathname}
              >
                <Link {...option} key={option.to}>
                  {option.label}
                </Link>
              </li>
            );
          })}
        </menu>
        <div className="window" role="tabpanel">
          <div className="window-body">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}
