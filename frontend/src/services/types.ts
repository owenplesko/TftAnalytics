import { z } from "zod";

export const summonerSchema = z.object({
  puuid: z.string(),
  region: z.string(),
  name: z.string(),
  tag: z.string(),
  summonerId: z.string(),
  profileIconId: z.number(),
  summonerLevel: z.number(),
  matchesBeforeTimestamp: z.string().datetime({ offset: true }).nullable(),
});

export const summonerStatsSchema = z.object({
  compCount: z.number(),
  totalSecondsIngame: z.number(),
  avgPlacement: z.number().nullable(),
  top4Count: z.number(),
  top1Count: z.number()
})

export type SummonerStats = z.infer<typeof summonerStatsSchema>;

export const companionSchema = z.object({
  contentId: z.string(),
  itemId: z.number(),
  skinId: z.number(),
  species: z.string(),
});

export const traitSchema = z.object({
  name: z.string(),
  numUnits: z.number(),
  style: z.number(),
  tierCurrent: z.number(),
  tierMax: z.number(),
});

export const unitSchema = z.object({
  characterId: z.string(),
  itemNames: z.string().array(),
  rarity: z.number(),
  tier: z.number(),
});

export type Unit = z.infer<typeof unitSchema>;

export const compDataSchema = z.object({
  companion: companionSchema,
  goldLeft: z.number(),
  lastRound: z.number(),
  level: z.number(),
  placement: z.number(),
  playersEliminated: z.number(),
  puuid: z.string(),
  timeEliminated: z.number(),
  totalDamageToPlayers: z.number(),
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
  matchDate: z.string().datetime({ offset: true }),
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
