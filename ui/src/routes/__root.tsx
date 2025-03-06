import { ThemeProvider } from "@/components/context/themeProvider";
import { SearchForm } from "@/components/searchForm";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

export const Route = createRootRoute({
  component: () => (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <main className="flex min-h-screen w-full flex-col items-center pt-2">
        <SearchForm />
        <div className="flex w-[1000px] flex-col items-center pt-8">
          <Outlet />
        </div>
      </main>
      <TanStackRouterDevtools />
    </ThemeProvider>
  ),
});
