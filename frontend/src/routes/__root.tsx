import { SearchForm } from "@/components/searchForm";
import { Button } from "@/components/ui/button";
import { type QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  Link,
  Outlet,
} from "@tanstack/react-router";

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()(
  {
    component: () => (
      <>
        <header className="sticky top-0 z-10 flex w-full bg-background px-4 py-2">
          <nav className="mr-auto flex items-start gap-2">
            <Link to="/units">
              <Button variant="link">Units</Button>
            </Link>
          </nav>
          <SearchForm />
        </header>
        <main className="flex min-h-screen w-full flex-col items-center pt-12">
          <Outlet />
        </main>
      </>
    ),
  },
);
