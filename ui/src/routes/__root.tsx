import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { QueryClient } from "@tanstack/react-query";

import { AuthContext } from "@/contexts/AuthContext";

interface RouterContext {
  queryClient: QueryClient;
  auth: AuthContext;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => {
    return (
      <div className="min-h-screen flex justify-center items-center bg-[#3a6ea5]">
        <div className="flex flex-col items-center p-2 w-full">
          <Outlet />
        </div>
      </div>
    );
  },
});
