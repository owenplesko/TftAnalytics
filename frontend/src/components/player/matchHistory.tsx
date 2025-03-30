import { getSummonerMatches } from "@/services/getSummonerMatches";
import { useSuspenseInfiniteQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useInView } from "react-intersection-observer";
import SummonerMatchCard from "../comp";

const MatchHistory: React.FC<{ puuid: string }> = ({ puuid }) => {
  const matchesQuery = useSuspenseInfiniteQuery(getSummonerMatches(puuid));
  const matches = matchesQuery.data.pages.flat();

  const [inViewRef, inView] = useInView();

  useEffect(() => {
    if (inView) matchesQuery.fetchNextPage();
  }, [inView, matchesQuery.fetchNextPage]);

  return (
    <>
      <div className="flex w-full flex-col gap-4 py-4">
        {matches.map((match) => (
          <SummonerMatchCard
            key={`${match.compData.puuid}_${match.tftMatch.id}`}
            summonerMatch={match}
          />
        ))}
      </div>
      <div ref={inViewRef} />
    </>
  );
};
export default MatchHistory;
