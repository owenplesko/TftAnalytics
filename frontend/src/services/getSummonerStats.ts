import { queryOptions } from "@tanstack/react-query";
import { summonerStatsSchema } from "./types";

export const getSummonerStats = (puuid: string) => queryOptions({
    queryKey: ["GET_SUMMONER_STATS", puuid], queryFn: async () => {
        const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/stats`;
        const res = await fetch(url);
        const data = await res.json();
        return summonerStatsSchema.array().parse(data);
    }
})
