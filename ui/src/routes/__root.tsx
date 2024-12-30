import { ThemeProvider } from "@/components/context/themeProvider";
import { SearchForm } from "@/components/searchForm";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

export const Route = createRootRoute({
  component: () => (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <main className="w-full flex-col min-h-screen flex items-center pt-2">
        <SearchForm />
        <div className="w-[1000px] pt-8 flex-col flex items-center">
          <Outlet />
        </div>
      </main>
      <TanStackRouterDevtools />
    </ThemeProvider>
  ),
});
