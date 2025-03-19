import {
  createFileRoute,
  Link,
  linkOptions,
  Outlet,
  useLocation,
} from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/_layout")({
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

  return (
    <div className="window w-1/2">
      <div className="title-bar">
        <div className="title-bar-text">
          {options.find((option) => option.to === location.pathname)?.label} -
          MSN Messenger Service
        </div>
      </div>
      <div className="window-body">
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
