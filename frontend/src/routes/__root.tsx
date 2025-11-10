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
      <div className="min-h-screen">
        <header className="sticky top-0 z-10 flex w-full border-b bg-background p-2">
          <nav className="mr-auto flex items-start gap-2">
            <Link to="/">
              <Button variant="link">Home</Button>
            </Link>
            <Link to="/units">
              <Button variant="link">Units</Button>
            </Link>
            <Link to="/items">
              <Button variant="link">Items</Button>
            </Link>
            <Link to="/traits">
              <Button variant="link">Traits</Button>
            </Link>
          </nav>
          <SearchForm />
        </header>
        <main className="flex w-full flex-col items-center pt-12">
          <Outlet />
        </main>
      </div>
    ),
  },
);
