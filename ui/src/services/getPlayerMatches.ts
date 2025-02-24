import { summonerMatchSchema } from "./types";

type Params = { puuid: string };

export async function getPlayerMatchHistory({ puuid }: Params) {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/matches`
  const res = await fetch(url);
  const data = await res.json();
  return summonerMatchSchema.array().parse(data);
}
