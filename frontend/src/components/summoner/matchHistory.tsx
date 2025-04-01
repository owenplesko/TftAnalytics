import { getSummonerMatches } from "@/services/getSummonerMatches";
import { useSuspenseInfiniteQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useInView } from "react-intersection-observer";
import SummonerMatchCard from "../comp";

const MatchHistory: React.FC<{ puuid: string; queueId: number | null }> = ({
  puuid,
  queueId,
}) => {
  const matchesQuery = useSuspenseInfiniteQuery(
    getSummonerMatches({ puuid, queueId }),
  );
  const matches = matchesQuery.data.pages.flat();

  const [inViewRef, inView] = useInView();

  useEffect(() => {
    if (inView) matchesQuery.fetchNextPage();
  }, [inView, matchesQuery.fetchNextPage]);

  return (
    <>
      {matches.map((match) => (
        <SummonerMatchCard
          key={`${match.compData.puuid}_${match.tftMatch.id}`}
          summonerMatch={match}
        />
      ))}
      <div ref={inViewRef} />
    </>
  );
};
export default MatchHistory;
