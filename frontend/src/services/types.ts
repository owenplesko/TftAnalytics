import { z } from "zod";

export const summonerSchema = z.object({
  puuid: z.string(),
  region: z.string(),
  name: z.string(),
  tag: z.string(),
  summonerId: z.string(),
  profileIconId: z.number(),
  summonerLevel: z.number(),
  statsUpdateTimestamp: z.string().datetime({ offset: true }).nullable(),
  matchesBeforeTimestamp: z.string().datetime({ offset: true }).nullable(),
});

export const summonerStatsSchema = z.object({
  queueId: z.number().int().nullable(),
  setNumber: z.number(),
  compCount: z.number().int(),
  avgPlacement: z.number(),
  top4Count: z.number(),
  top4Rate: z.number(),
  top1Rate: z.number(),
  top1Count: z.number()
})

export type SummonerStats = z.infer<typeof summonerStatsSchema>;

export const companionSchema = z.object({
  contentId: z.string(),
  itemId: z.number().int(),
  skinId: z.number().int(),
  species: z.string(),
});

export const traitSchema = z.object({
  name: z.string(),
  numUnits: z.number().int(),
  style: z.number().int(),
  tierCurrent: z.number().int(),
  tierMax: z.number().int(),
});

export const unitSchema = z.object({
  characterId: z.string(),
  itemNames: z.string().array(),
  rarity: z.number().int(),
  tier: z.number().int(),
});

export type Unit = z.infer<typeof unitSchema>;

export const compDataSchema = z.object({
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

export type CompData = z.infer<typeof compDataSchema>

export const matchSchema = z.object({
  id: z.string(),
  gameVersion: z.string(),
  queueId: z.number(),
  gameType: z.string(),
  setNumber: z.number(),
  matchDate: z.string().datetime(),
});

export const summonerMatchSchema = z.object({
  compData: compDataSchema,
  tftMatch: matchSchema,
});

export type SummonerMatch = z.infer<typeof summonerMatchSchema>;

export const matchCompSchema = z.object({
  compData: compDataSchema,
  tftSummoner: summonerSchema
})

export type MatchComp = z.infer<typeof matchCompSchema>

export const rankDataSchema = z.object({
  tier: z.enum([
    "IRON",
    "BRONZE",
    "SILVER",
    "GOLD",
    "PLATINUM",
    "EMERALD",
    "DIAMOND",
    "MASTER",
    "GRANDMASTER",
    "CHALLENGER",
  ]),
  rank: z.enum(["I", "II", "III", "IV"]),
  leaguePoints: z.number(),
});

export const leaderboardPositionSchema = z.object({
  leaderboard: z.string(),
  position: z.number(),
  top: z.number(),
});

export const rankSchema = z.object({
  rankData: rankDataSchema,
  leaderboardPosition: leaderboardPositionSchema,
});

export type Rank = z.infer<typeof rankSchema>;
