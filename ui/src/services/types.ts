import { z } from "zod";

export const summonerSchema = z.object({
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
    tftMatch: matchSchema
})

export const rankDataSchema = z.object({
    tier: z.enum(["IRON", "BRONZE", "SILVER", "GOLD", "PLATINUM", "EMERALD", "DIAMOND", "MASTER", "GRANDMASTER", "CHALLENGER"]),
    rank: z.enum(["I", "II", "III", "IV"]),
    leaguePoints: z.number()
});

export const leaderboardPositionSchema = z.object({
    leaderboard: z.string(),
    position: z.number(),
    top: z.number()
});

export const rankSchema = z.object({
    rankData: rankDataSchema,
    leaderboardPosition: leaderboardPositionSchema
});
