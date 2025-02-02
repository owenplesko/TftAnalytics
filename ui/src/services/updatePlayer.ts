type Params = { puuid: string };

export async function updatePlayer({ puuid }: Params) {
    await fetch(`http://localhost:8080/v1/summoner/by-puuid/${puuid}/update`);
}
