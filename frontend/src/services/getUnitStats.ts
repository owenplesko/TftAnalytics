import { queryOptions } from "@tanstack/react-query";
import { unitStatSchema } from "./types";

export const GET_UNIT_STATS_KEY = "GET_MATCH_COMPS"

export const getUnitStats = () => queryOptions({
  queryKey: [GET_UNIT_STATS_KEY],
  queryFn: async () => {
    const url = `${import.meta.env.VITE_BACKEND_URL}/v1/stats/units`;
    const res = await fetch(url);
    const data = await res.json();
    return unitStatSchema.array().parse(data);
  }
})
