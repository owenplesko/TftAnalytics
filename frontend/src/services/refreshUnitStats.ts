export const refreshUnitStats = async () => {
  const url = `${import.meta.env.VITE_BACKEND_URL}/v1/stats/units/refresh`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`status code: ${res.status}`)
  }
}
