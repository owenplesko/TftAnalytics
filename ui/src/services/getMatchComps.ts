import { queryOptions } from "@tanstack/react-query";
import { matchCompSchema } from "./types";

export const getMatchComps = (matchId: string) => queryOptions({
    queryKey: ["GET_MATCH_COMPS", matchId],
    queryFn: async () => {
        const url = `${import.meta.env.VITE_BACKEND_URL}/v1/match/${matchId}/comps`;
        const res = await fetch(url);
        const data = await res.json();
        return matchCompSchema.array().parse(data);
    },
})
