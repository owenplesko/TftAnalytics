import SummonerMatch from "@/components/comp";
import { Button } from "@/components/ui/button";
import { sentenceCase } from "@/lib/utils";
import { getPlayer } from "@/services/getPlayer";
import { getPlayerMatchHistory } from "@/services/getPlayerMatches";
import { getRank } from "@/services/getRank";
import { Rank, SummonerStats } from "@/services/types";
import { updatePlayer } from "@/services/updatePlayer";
import {
  useMutation,
  useSuspenseInfiniteQuery,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { IconLoader2 } from "@tabler/icons-react";
import TimeSince from "@/components/timeSince";
import { useInView } from "react-intersection-observer";
import React, { useEffect } from "react";
import { Badge } from "@/components/ui/badge";
import { getPlayerAggregates } from "@/services/getPlayerStats";

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params: { name, tag }, context: { queryClient } }) => {
    const { region, puuid, summonerId } = await queryClient.ensureQueryData(
      getPlayer(name, tag),
    );

    await Promise.all([
      queryClient.ensureQueryData(getRank(region, summonerId)),
      queryClient.ensureQueryData(getPlayerAggregates(puuid)),
      queryClient.ensureInfiniteQueryData(getPlayerMatchHistory(puuid)),
    ]);

    return { name, tag };
  },
});

function Player() {
  const loader = Route.useLoaderData();

  const playerQuery = useSuspenseQuery(getPlayer(loader.name, loader.tag));

  const rankQuery = useSuspenseQuery(
    getRank(playerQuery.data.region, playerQuery.data.summonerId),
  );

  const aggregatesQuery = useSuspenseQuery(
    getPlayerAggregates(playerQuery.data.puuid),
  );
  const rankAggregates = aggregatesQuery.data.find(
    ({ queueId }) => queueId == 1100,
  );

  const matchesQuery = useSuspenseInfiniteQuery(
    getPlayerMatchHistory(playerQuery.data.puuid),
  );
  const matches = matchesQuery.data.pages.flat();

  const updateMutation = useMutation({
    mutationKey: ["UPDATE_PLAYER", playerQuery.data.puuid],
    mutationFn: updatePlayer,
    onSuccess: () => {
      playerQuery.refetch();
      rankQuery.refetch();
      matchesQuery.refetch();
      aggregatesQuery.refetch();
    },
  });

  const [inViewRef, inView] = useInView();

  useEffect(() => {
    if (inView) matchesQuery.fetchNextPage();
  }, [inView, matchesQuery.fetchNextPage]);

  return (
    <>
      <div className="grid w-full grid-cols-[auto_1fr] items-center justify-items-start gap-2 border-b pb-4">
        <img
          className="row-span-4 h-full rounded border"
          src={`/profileicon/profileicon${playerQuery.data.profileIconId}.png`}
        />
        <div className="flex items-center gap-1">
          <h1 className="text-3xl font-bold">
            <span>{playerQuery.data.name}</span>
            <span className="text-muted-foreground">
              #{playerQuery.data.tag}
            </span>
          </h1>
          <Badge variant="secondary">{playerQuery.data.region}</Badge>
        </div>
        {rankQuery.isSuccess ? (
          <RankLine rank={rankQuery.data} />
        ) : (
          <span className="text-destructive">error loading rank</span>
        )}
        <Button
          onClick={() => updateMutation.mutate(playerQuery.data.puuid)}
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
        <span className="text-muted-foreground">
          {"Updated "}
          {playerQuery.data.updateTimestamp ? (
            <TimeSince date={new Date(playerQuery.data.updateTimestamp)} />
          ) : (
            "never"
          )}
        </span>
      </div>
      <PlacementLine aggregates={rankAggregates} />
      <div className="flex w-full flex-col gap-4 py-4">
        {matches.map((match) => (
          <SummonerMatch
            key={`${match.compData.puuid}_${match.tftMatch.id}`}
            summonerMatch={match}
          />
        ))}
      </div>
      <div ref={inViewRef} />
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
      <span className="text-muted-foreground">{`Rank #${rank.leaderboardPosition.position}`}</span>
      <span className="text-muted-foreground">{`Top ${rank.leaderboardPosition.top.toPrecision(2)}%`}</span>
    </div>
  );
};

const PlacementLine: React.FC<{
  aggregates: SummonerStats | undefined;
}> = ({ aggregates }) => {
  if (!aggregates) return null;

  return (
    <div className="flex flex-col">
      <span>{`${aggregates.compCount} Games Played`}</span>
      <span>{`${aggregates.avgPlacement.toPrecision(3)} AVP`}</span>
      <span>{`${aggregates.top4Rate.toPrecision(3)}% Top 4`}</span>
      <span>{`${aggregates.top1Rate.toPrecision(3)}% Top 1`}</span>
    </div>
  );
};
