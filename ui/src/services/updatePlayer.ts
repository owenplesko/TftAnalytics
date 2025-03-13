export async function updatePlayer(puuid: string) {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/summoner/by-puuid/${puuid}/update`;
  await fetch(url);
}
