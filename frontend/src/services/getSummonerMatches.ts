import { infiniteQueryOptions } from "@tanstack/react-query";
import { summonerMatchSchema } from "./types";

export const GET_SUMMONER_MATCHES_KEY = "GET_SUMMONER_MATCHES";

type GetSummonerMatchesParams = {
  puuid: string;
  queueId?: number | null;
  initialPageParam?: string | null
}

export const getSummonerMatches = ({ puuid, queueId, initialPageParam }: GetSummonerMatchesParams) => infiniteQueryOptions({
  queryKey: [GET_SUMMONER_MATCHES_KEY, puuid, queueId],
  queryFn: async ({ pageParam }) => {
    const url = new URL(`${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/matches`)
    if (queueId) url.searchParams.append("queue", queueId.toString())
    if (pageParam) url.searchParams.append("before", pageParam)

    const res = await fetch(url);
    const data = await res.json();
    return summonerMatchSchema.array().parse(data);
  },
  initialPageParam: initialPageParam,
  getNextPageParam: (lastPage) => lastPage[lastPage.length - 1]?.tftMatch.matchDate
})
