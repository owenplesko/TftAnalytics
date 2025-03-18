import { SearchForm } from "@/components/searchForm";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { type QueryClient } from "@tanstack/react-query";

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()(
  {
    component: () => (
      <main className="flex min-h-screen w-full flex-col items-center pt-4">
        <SearchForm />
        <div className="items-center-4 flex w-[1000px] flex-col pt-4">
          <Outlet />
        </div>
      </main>
    ),
  },
);
