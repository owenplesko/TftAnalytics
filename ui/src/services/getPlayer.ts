import { summonerSchema } from "./types";

type Params = { name: string; tag: string };

export async function getPlayer({ name, tag }: Params) {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-name-tag/${name}/${tag}`
  const res = await fetch(url);
  const data = await res.json();
  return summonerSchema.parse(data);
}
