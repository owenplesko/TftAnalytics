import { queryOptions } from "@tanstack/react-query";
import { rankSchema } from "./types";

export const getRank = (region: string, summonerId: string) => queryOptions({
  queryKey: ["GET_RANK", region, summonerId], queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/leaderboard/${region}/${summonerId}`;
    const res = await fetch(url);
    if (res.status == 404) return null;

    const data = await res.json();
    return rankSchema.parse(data);
  }
})
