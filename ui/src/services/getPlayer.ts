import { summonerSchema } from "./types";

export async function getPlayer(name: string, tag: string) {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-name-tag/${name}/${tag}`;
  const res = await fetch(url);
  const data = await res.json();
  return summonerSchema.parse(data);
}
