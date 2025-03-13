import { matchCompSchema } from "./types";

export async function getMatchComps(matchId: string) {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/match/${matchId}/comps`;
    const res = await fetch(url);
    const data = await res.json();
    return matchCompSchema.array().parse(data);
}
