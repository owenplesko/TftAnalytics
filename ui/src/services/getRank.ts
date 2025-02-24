import { rankSchema } from "./types";

type Params = { summonerId: string; region: string };

export async function getRank({ summonerId, region }: Params) {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/leaderboard/${region}/${summonerId}`
    const res = await fetch(url);
    const data = await res.json();
    return rankSchema.parse(data);
}
