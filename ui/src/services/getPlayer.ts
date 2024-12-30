import { z } from "zod";

const summonerSchema = z.object({
  puuid: z.string(),
  name: z.string(),
  tag: z.string(),
  summonerId: z.string(),
  profileIconId: z.number(),
  summonerLevel: z.number(),
  fullUpdateTimestamp: z.string().nullable(),
});

type Params = { name: string; tag: string };

export async function getPlayer({ name, tag }: Params) {
  const res = await fetch(`http://localhost:8080/account/${name}/${tag}`);
  const data = await res.json();
  return summonerSchema.parse(data);
}
