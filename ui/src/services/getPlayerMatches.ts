import { summonerMatchSchema } from "./types";

export async function getPlayerMatchHistory(puuid: string) {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/matches`;
  const res = await fetch(url);
  const data = await res.json();
  return summonerMatchSchema.array().parse(data);
}
