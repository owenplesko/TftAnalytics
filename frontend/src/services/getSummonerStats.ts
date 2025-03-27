import { queryOptions } from "@tanstack/react-query";
import { summonerStatsSchema } from "./types";

export const getSummonerStats = (puuid: string, setNumber: number, queueId?: number | null) => queryOptions({
    queryKey: ["GET_SUMMONER_STATS", puuid], queryFn: async () => {
        let url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/stats?set=${setNumber}`;
        if (queueId) url += `&queue=${queueId}`

        const res = await fetch(url);
        if (res.status == 404) return null;

        const data = await res.json();
        return summonerStatsSchema.parse(data);
    }
})
