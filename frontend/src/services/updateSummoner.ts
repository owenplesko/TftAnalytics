export const updateSummoner = async (puuid: string) => {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/update`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`status code: ${res.status}`)
  }
}
