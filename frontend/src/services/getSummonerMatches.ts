import { infiniteQueryOptions } from "@tanstack/react-query";
import { summonerMatchSchema } from "./types";


export const getSummonerMatches = (puuid: string, initialPageParam?: string | null) => infiniteQueryOptions({
  queryKey: ["GET_SUMMONER_MATCHES", puuid],
  queryFn: async ({ pageParam }) => {
    let url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/matches`
    if (pageParam) url += `?before=${pageParam}`

    const res = await fetch(url);
    const data = await res.json();
    return summonerMatchSchema.array().parse(data);
  },
  initialPageParam: initialPageParam,
  getNextPageParam: (lastPage) => lastPage[lastPage.length - 1]?.tftMatch.matchDate
})
