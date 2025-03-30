import { queryOptions } from "@tanstack/react-query";
import { rankSchema } from "./types";

export const GET_SUMMONER_RANK_KEY = "GET_SUMMONER_RANK";

export const getSummonerRank = (region: string, summonerId: string) => queryOptions({
  queryKey: [GET_SUMMONER_RANK_KEY, region, summonerId], queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/leaderboard/${region}/${summonerId}`;
    const res = await fetch(url);
    if (res.status == 404) return null;

    const data = await res.json();
    return rankSchema.parse(data);
  }
})
