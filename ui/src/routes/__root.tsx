import { createRootRoute, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
  component: () => {
    return (
      <div className="min-h-screen flex justify-center items-center bg-[#3a6ea5]">
        <div className="p-2">
          <Outlet />
        </div>
      </div>
    );
  },
});
