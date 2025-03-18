import { queryOptions } from "@tanstack/react-query";
import { summonerSchema } from "./types";

export const getPlayer = (name: string, tag: string) => queryOptions({
  queryKey: ["GET_PLAYER", name, tag], queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-name-tag/${name}/${tag}`;
    const res = await fetch(url);
    const data = await res.json();
    return summonerSchema.parse(data);
  }
})
