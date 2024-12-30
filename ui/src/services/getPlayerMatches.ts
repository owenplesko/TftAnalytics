import { z } from "zod";

const companionSchema = z.object({
  contentId: z.string(),
  itemId: z.number().int(),
  skinId: z.number().int(),
  species: z.string(),
});

const traitSchema = z.object({
  name: z.string(),
  numUnits: z.number().int(),
  style: z.number().int(),
  tierCurrent: z.number().int(),
  tierMax: z.number().int(),
});

const unitSchema = z.object({
  characterId: z.string(),
  itemNames: z.string().array(),
  rarity: z.number().int(),
  tier: z.number().int(),
});

const compDataSchema = z.object({
  augments: z.string().array(),
  companion: companionSchema,
  goldLeft: z.number().int(),
  lastRound: z.number().int(),
  level: z.number().int(),
  placement: z.number().int(),
  playersEliminated: z.number().int(),
  puuid: z.string(),
  timeEliminated: z.number(),
  totalDamageToPlayers: z.number().int(),
  traits: traitSchema.array(),
  units: unitSchema.array(),
});

const matchSchema = z.object({
  matchId: z.string(),
  gameVersion: z.string(),
  queueId: z.number(),
  gameType: z.string(),
  setNumber: z.number(),
  matchDate: z.string().datetime(),
});

const summonerMatchSchema = matchSchema.merge(
  z.object({ compData: compDataSchema }),
);

type Params = { puuid: string };

export async function GetPlayerMatches({ puuid }: Params) {
  const res = await fetch(`http://localhost:8080/summoner/${puuid}/matches`);
  const data = await res.json();
  return summonerMatchSchema.array().parse(data);
}
