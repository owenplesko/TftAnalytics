import { getSummoner } from "@/services/getSummoner";
import { getSummonerMatches } from "@/services/getSummonerMatches";
import { getSummonerRank } from "@/services/getSummonerRank";
import { createFileRoute } from "@tanstack/react-router";
import { getSummonerStats } from "@/services/getSummonerStats";
import QueueSummary from "@/components/summoner/queueSummary";
import PlayerHeader from "@/components/summoner/header";
import MatchHistory from "@/components/summoner/matchHistory";

const CURRENT_SET_NUMBER = 13;

export const Route = createFileRoute("/summoner/$name/$tag")({
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

    const data = {
      name,
      tag,
      region,
      summonerId,
      puuid,
      setNumber: CURRENT_SET_NUMBER,
    };
    return data;
  },
});

function Player() {
  const loader = Route.useLoaderData();

  return (
    <div className="flex w-[1000px] flex-col gap-4 pt-4">
      <PlayerHeader name={loader.name} tag={loader.tag} />
      <QueueSummary
        region={loader.region}
        summonerId={loader.summonerId}
        puuid={loader.puuid}
        setNumber={loader.setNumber}
      />
      <MatchHistory puuid={loader.puuid} />
    </div>
  );
}
