import { queryOptions } from "@tanstack/react-query";
import { summonerSchema } from "./types";

export const GET_SUMMONER_KEY = "GET_SUMMONER";

export const getSummoner = (name: string, tag: string) => queryOptions({
  queryKey: [GET_SUMMONER_KEY, name, tag],
  queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-name-tag/${name}/${tag}`;
    const res = await fetch(url);
    const data = await res.json();
    return summonerSchema.parse(data);
  }
})
