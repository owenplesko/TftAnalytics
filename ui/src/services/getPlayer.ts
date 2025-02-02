import { z } from "zod";

const summonerSchema = z.object({
  puuid: z.string(),
  region: z.string().nullable(),
  name: z.string().nullable(),
  tag: z.string().nullable(),
  summonerId: z.string().nullable(),
  profileIconId: z.number().nullable(),
  summonerLevel: z.number().nullable(),
  fullUpdateTimestamp: z.string().datetime().nullable(),
  backgroundUpdateTimestamp: z.string().datetime().nullable(),
  flags: z.string().array(),
});

type Params = { name: string; tag: string };

export async function getPlayer({ name, tag }: Params) {
  const res = await fetch(`http://localhost:8080/v1/summoner/by-name-tag/${name}/${tag}`);
  const data = await res.json();
  console.log(data)
  return summonerSchema.parse(data);
}
