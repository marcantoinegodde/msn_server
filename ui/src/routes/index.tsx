import { createFileRoute } from "@tanstack/react-router";

import LoginWindows from "@/components/LoginWindow";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return <LoginWindows />;
}
