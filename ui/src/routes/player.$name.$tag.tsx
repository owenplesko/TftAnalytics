import SummonerMatch from "@/components/comp";
import { Button } from "@/components/ui/button";
import { sentenceCase } from "@/lib/utils";
import { getPlayer } from "@/services/getPlayer";
import { getPlayerMatchHistory } from "@/services/getPlayerMatches";
import { getRank } from "@/services/getRank";
import { Rank } from "@/services/types";
import { updatePlayer } from "@/services/updatePlayer";
import { useMutation, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { IconLoader2 } from "@tabler/icons-react";
import TimeSince from "@/components/timeSince";

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params: { name, tag }, context: { queryClient } }) => {
    const { region, puuid, summonerId } = await queryClient.ensureQueryData(
      getPlayer(name, tag),
    );
    await queryClient.ensureQueryData(getRank(region, summonerId));
    await queryClient.ensureQueryData(getPlayerMatchHistory(puuid));
  },
});

function Player() {
  const params = Route.useParams();

  const playerQuery = useSuspenseQuery(getPlayer(params.name, params.tag));
  const player = playerQuery.data;

  const rankQuery = useSuspenseQuery(getRank(player.region, player.summonerId));
  const rank = rankQuery.data;

  const matchesQuery = useSuspenseQuery(getPlayerMatchHistory(player.puuid));
  const matches = matchesQuery.data;

  const updateMutation = useMutation({
    mutationFn: updatePlayer,
    onSuccess: () => {
      playerQuery.refetch();
      rankQuery.refetch();
      matchesQuery.refetch();
    },
  });

  return (
    <>
      <div className="flex w-full items-center gap-4 border-b pb-4">
        <img
          className="rounded-sm border"
          width={124}
          height={124}
          src={`/profileicon/profileicon${player.profileIconId}.png`}
        />
        <div className="flex flex-col items-start gap-2">
          <h1 className="text-3xl font-semibold">
            <span>{player.name}</span>
            <span className="text-muted-foreground">#{player.tag}</span>
          </h1>
          <RankLine rank={rank} />
          <Button
            variant="outline"
            onClick={() => updateMutation.mutate(player.puuid)}
            disabled={updateMutation.isPending}
          >
            {
              {
                idle: "Update",
                pending: (
                  <>
                    <IconLoader2 className="mr-1 animate-spin ease-in-out" />
                    Updating...
                  </>
                ),
                success: "Updated",
                error: "Error",
              }[updateMutation.status]
            }
          </Button>
          <span className="text-sm text-muted-foreground">
            {"Updated "}
            {player.updateTimestamp ? (
              <TimeSince date={new Date(player.updateTimestamp)} />
            ) : (
              "never"
            )}
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
