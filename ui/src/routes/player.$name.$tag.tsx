import { Button } from "@/components/ui/button";
import { getPlayer } from "@/services/getPlayer";
import { GetPlayerMatches } from "@/services/getPlayerMatches";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params }) => {
    const player = await getPlayer(params);
    const matches = await GetPlayerMatches(player);
    return { player, matches };
  },
});

function Player() {
  const { player, matches } = Route.useLoaderData();
  return (
    <>
      <div className="w-full flex items-center gap-4 border-b pb-4">
        <img
          width={124}
          height={124}
          src={`/profileicon/profileicon${player.profileIconId}.png`}
        />
        <div className="flex flex-col items-start gap-2">
          <h1 className="text-3xl font-semibold">
            <span>{player.name}</span>
            <span className="text-muted">#{player.tag}</span>
          </h1>
          <Button variant="outline">Update</Button>
          <span className="text-sm text-muted-foreground">
            Updated 2 hours ago
          </span>
        </div>
      </div>
      <span>{JSON.stringify(matches)}</span>
    </>
  );
}
