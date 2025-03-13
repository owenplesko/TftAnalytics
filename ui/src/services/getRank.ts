import { rankSchema } from "./types";

export async function getRank(region: string, summonerId: string) {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/leaderboard/${region}/${summonerId}`;
  const res = await fetch(url);
  if (res.status == 404) return null;

  const data = await res.json();
  return rankSchema.parse(data);
}
