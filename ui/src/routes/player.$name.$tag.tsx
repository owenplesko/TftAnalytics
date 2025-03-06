import CompSummary from "@/components/comp";
import { Button } from "@/components/ui/button";
import { sentenceCase } from "@/lib/utils";
import { getPlayer } from "@/services/getPlayer";
import { getPlayerMatchHistory } from "@/services/getPlayerMatches";
import { getRank } from "@/services/getRank";
import { Rank } from "@/services/types";
import { updatePlayer } from "@/services/updatePlayer";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params }) => {
    const player = await getPlayer(params);
    const matches = await getPlayerMatchHistory(player);
    const rank = await getRank({
      summonerId: player.summonerId,
      region: player.region,
    });
    return { player, matches, rank };
  },
});

function Player() {
  const { player, rank, matches } = Route.useLoaderData();

  return (
    <>
      <div className="flex w-full items-center gap-4 border-b pb-4">
        <img
          className="rounded-sm border"
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
          <RankLine rank={rank} />
          <Button
            variant="outline"
            onClick={() => updatePlayer({ puuid: player.puuid })}
          >
            Update
          </Button>
          <span className="text-muted-foreground text-sm">
            Updated 2 hours ago
          </span>
        </div>
      </div>
      <div className="flex w-full flex-col gap-4 pt-4">
        {matches.map((match) => (
          <CompSummary summonerMatch={match} />
        ))}
      </div>
    </>
  );
}

const RankLine: React.FC<{
  rank: Rank | null;
}> = ({ rank }) => {
  return rank ? (
    <>
      <div className="flex flex-row gap-2">
        <img
          width={24}
          height={24}
          src={`/rank/${rank.rankData.tier.toLowerCase()}.svg`}
        />
        <span>{`${sentenceCase(rank.rankData.tier)}`}</span>
        <span>{`${rank.rankData.leaguePoints} LP`}</span>
        <span>{`Rank #${rank.leaderboardPosition.position}`}</span>
        <span>{`Top ${rank.leaderboardPosition.top.toPrecision(2)}%`}</span>
      </div>
    </>
  ) : (
    <>
      <div className="flex flex-row gap-2">
        <img width={24} height={24} src="/rank/unranked.svg" />
        <span>Unranked</span>
      </div>
    </>
  );
};
