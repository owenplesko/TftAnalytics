import { SearchForm } from "@/components/searchForm";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { type QueryClient } from "@tanstack/react-query";

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()(
  {
    component: () => (
      <main className="flex min-h-screen w-full flex-col items-center pt-2">
        <SearchForm />
        <div className="flex w-[1000px] flex-col items-center pt-8">
          <Outlet />
        </div>
        <TanStackRouterDevtools />
        <ReactQueryDevtools />
      </main>
    ),
  },
);
