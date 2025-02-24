import { Button } from "@/components/ui/button";
import { getPlayer } from "@/services/getPlayer";
import { getPlayerMatchHistory } from "@/services/getPlayerMatches";
import { getRank } from "@/services/getRank";
import { updatePlayer } from "@/services/updatePlayer";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params }) => {
    const player = await getPlayer(params);
    const rank = await getRank({
      summonerId: player.summonerId,
      region: player.region,
    });
    const matches = await getPlayerMatchHistory(player);
    return { player, rank, matches };
  },
});

function Player() {
  const { player, rank, matches } = Route.useLoaderData();

  return (
    <>
      <div className="w-full flex items-center gap-4 border-b pb-4">
        <img
          width={124}
          height={124}
          src={
            player.profileIconId
              ? `/profileicon/profileicon${player.profileIconId}.png`
              : "/profileicon/profileicon29.png"
          }
        />
        <div className="flex flex-col items-start gap-2">
          <h1 className="text-3xl font-semibold">
            <span>{player.name}</span>
            <span className="text-muted">#{player.tag}</span>
          </h1>
          <span className="text-sm">
            {rank.rankData.tier} {rank.rankData.rank}{" "}
            {rank.rankData.leaguePoints} #{rank.leaderboardPosition.position}{" "}
            Top {rank.leaderboardPosition.top.toPrecision(2)}%
          </span>
          <Button
            variant="outline"
            onClick={() => updatePlayer({ puuid: player.puuid })}
          >
            Update
          </Button>
          <span className="text-sm text-muted-foreground">
            Updated 2 hours ago
          </span>
        </div>
      </div>
      <span>{matches.length}</span>
    </>
  );
}
