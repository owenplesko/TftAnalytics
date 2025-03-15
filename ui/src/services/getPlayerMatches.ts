import { queryOptions } from "@tanstack/react-query";
import { summonerMatchSchema } from "./types";

export const getPlayerMatchHistory = (puuid: string) => queryOptions({
  queryKey: ["GET_PLAYER_MATCH_HISTORY", puuid], queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/matches`;
    const res = await fetch(url);
    const data = await res.json();
    return summonerMatchSchema.array().parse(data);
  }
})
