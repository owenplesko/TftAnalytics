import { infiniteQueryOptions } from "@tanstack/react-query";
import { summonerMatchSchema } from "./types";

export const getPlayerMatchHistory = (puuid: string) => infiniteQueryOptions({
  queryKey: ["GET_PLAYER_MATCH_HISTORY", puuid],
  queryFn: async ({ pageParam }) => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/matches?before=${pageParam}`;
    const res = await fetch(url);
    const data = await res.json();
    return summonerMatchSchema.array().parse(data);
  },
  initialPageParam: new Date().toISOString(),
  getNextPageParam: (lastPage) => lastPage[lastPage.length - 1]?.tftMatch.matchDate
})
