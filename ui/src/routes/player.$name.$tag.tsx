import SummonerMatch from "@/components/comp";
import { Button } from "@/components/ui/button";
import { formatTimeSince, sentenceCase } from "@/lib/utils";
import { getPlayer } from "@/services/getPlayer";
import { getPlayerMatchHistory } from "@/services/getPlayerMatches";
import { getRank } from "@/services/getRank";
import { Rank } from "@/services/types";
import { updatePlayer } from "@/services/updatePlayer";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params }) => {
    const player = await getPlayer(params.name, params.tag);
    const matches = await getPlayerMatchHistory(player.puuid);
    const rank = await getRank(player.region, player.summonerId);
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
            <span className="text-muted-foreground">#{player.tag}</span>
          </h1>
          <RankLine rank={rank} />
          <Button variant="outline" onClick={() => updatePlayer(player.puuid)}>
            Update
          </Button>
          <span className="text-sm text-muted-foreground">
            {`Updated ${player.updateTimestamp ? formatTimeSince(player.updateTimestamp) : "never"}`}
          </span>
        </div>
      </div>
      <div className="flex w-full flex-col gap-4 pt-4">
        {matches.map((match) => (
          <SummonerMatch summonerMatch={match} />
        ))}
      </div>
    </>
  );
}

const RankLine: React.FC<{
  rank: Rank | null;
}> = ({ rank }) => {
  if (rank === null)
    return (
      <div className="flex flex-row gap-2">
        <img width={24} height={24} src="/rank/unranked.svg" />
        <span>Unranked</span>
      </div>
    );

  return (
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
  );
};
