import { queryOptions } from "@tanstack/react-query";
import { summonerSchema } from "./types";

export const getSummoner = (name: string, tag: string) => queryOptions({
  queryKey: ["GET_SUMMONER", name, tag], queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-name-tag/${name}/${tag}`;
    const res = await fetch(url);
    const data = await res.json();
    return summonerSchema.parse(data);
  }
})
