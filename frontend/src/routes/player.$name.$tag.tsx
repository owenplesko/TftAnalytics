import SummonerMatch from "@/components/comp";
import { Button } from "@/components/ui/button";
import { sentenceCase } from "@/lib/utils";
import { getSummoner } from "@/services/getSummoner";
import { getSummonerMatches } from "@/services/getSummonerMatches";
import { getSummonerRank } from "@/services/getSummonerRank";
import { Rank, SummonerStats } from "@/services/types";
import { updateSummoner } from "@/services/updateSummoner";
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
import { getSummonerStats } from "@/services/getSummonerStats";

const CURRENT_SET_NUMBER = 13;

export const Route = createFileRoute("/player/$name/$tag")({
  component: Player,
  loader: async ({ params: { name, tag }, context: { queryClient } }) => {
    const { region, puuid, summonerId } = await queryClient.ensureQueryData(
      getSummoner(name, tag),
    );

    await Promise.all([
      queryClient.ensureQueryData(getSummonerRank(region, summonerId)),
      queryClient.ensureQueryData(getSummonerStats(puuid, CURRENT_SET_NUMBER)),
      queryClient.ensureInfiniteQueryData(getSummonerMatches(puuid)),
    ]);

    return { name, tag };
  },
});

function Player() {
  const loader = Route.useLoaderData();

  const summonerQuery = useSuspenseQuery(getSummoner(loader.name, loader.tag));

  const rankQuery = useSuspenseQuery(
    getSummonerRank(summonerQuery.data.region, summonerQuery.data.summonerId),
  );

  const statsQuery = useSuspenseQuery(
    getSummonerStats(summonerQuery.data.puuid, CURRENT_SET_NUMBER),
  );

  const matchesQuery = useSuspenseInfiniteQuery(
    getSummonerMatches(summonerQuery.data.puuid),
  );
  const matches = matchesQuery.data.pages.flat();

  const updateMutation = useMutation({
    mutationKey: ["UPDATE_SUMMONER", summonerQuery.data.puuid],
    mutationFn: updateSummoner,
    onSuccess: () => {
      summonerQuery.refetch();
      rankQuery.refetch();
      matchesQuery.refetch();
      statsQuery.refetch();
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
          src={`/profileicon/profileicon${summonerQuery.data.profileIconId}.png`}
        />
        <div className="flex items-center gap-1">
          <h1 className="text-3xl font-bold">
            <span>{summonerQuery.data.name}</span>
            <span className="text-muted-foreground">
              #{summonerQuery.data.tag}
            </span>
          </h1>
          <Badge variant="secondary">{summonerQuery.data.region}</Badge>
        </div>
        {rankQuery.isSuccess ? (
          <RankLine rank={rankQuery.data} />
        ) : (
          <span className="text-destructive">error loading rank</span>
        )}
        <Button
          onClick={() => updateMutation.mutate(summonerQuery.data.puuid)}
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
          {summonerQuery.data.updateTimestamp ? (
            <TimeSince date={new Date(summonerQuery.data.updateTimestamp)} />
          ) : (
            "never"
          )}
        </span>
      </div>
      <PlacementLine stats={statsQuery.data} />
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
  stats: SummonerStats[];
}> = ({ stats: allStats }) => {
  return (
    <div className="flex gap-2">
      {allStats.map((stats) => (
        <div className="flex flex-col">
          <span>{`${stats.queueId ?? "all games"}`}</span>
          <span>{`${stats.compCount} Games Played`}</span>
          <span>{`${stats.avgPlacement.toPrecision(3)} AVP`}</span>
          <span>{`${stats.top4Rate.toPrecision(3)}% Top 4`}</span>
          <span>{`${stats.top1Rate.toPrecision(3)}% Top 1`}</span>
        </div>
      ))}
    </div>
  );
};
