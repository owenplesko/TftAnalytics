import { formatSeconds, sentenceCase } from "@/lib/utils";
import { getSummonerRank } from "@/services/getSummonerRank";
import { getSummonerStats } from "@/services/getSummonerStats";
import { useSuspenseQuery } from "@tanstack/react-query";

const QueueSummary: React.FC<{
  region: string;
  summonerId: string;
  puuid: string;
  queueId: number | null;
  setNumber: number;
}> = ({ region, summonerId, puuid, queueId, setNumber }) => {
  return (
    <div className="mr-auto flex flex-col items-center justify-center gap-2 rounded border bg-card p-4">
      <RankSummary region={region} summonerId={summonerId} />
      <QueueStatsSummary
        puuid={puuid}
        queueId={queueId}
        setNumber={setNumber}
      />
    </div>
  );
};
export default QueueSummary;

const RankSummary: React.FC<{
  region: string;
  summonerId: string;
}> = ({ region, summonerId }) => {
  const query = useSuspenseQuery(getSummonerRank(region, summonerId));

  if (query.isError)
    return <span className="text-destructive">error loading rank</span>;

  if (query.data === null)
    return (
      <div className="flex flex-row gap-2">
        <img width={24} height={24} src="/rank/unranked.svg" />
        <span>Unranked</span>
      </div>
    );

  return (
    <>
      <div className="flex flex-row gap-2">
        <img
          width={24}
          height={24}
          src={`/rank/${query.data.rankData.tier.toLowerCase()}.svg`}
        />
        <span className="text-2xl">{`${sentenceCase(query.data.rankData.tier)} ${query.data.rankData.leaguePoints} LP`}</span>
      </div>
      <span className="text-sm">{`Rank #${query.data.leaderboardPosition.position.toLocaleString()} Top ${query.data.leaderboardPosition.top.toPrecision(3)}%`}</span>
    </>
  );
};

const QueueStatsSummary: React.FC<{
  puuid: string;
  queueId: number | null;
  setNumber: number;
}> = ({ puuid, queueId, setNumber }) => {
  const query = useSuspenseQuery(getSummonerStats(puuid, setNumber, queueId));

  // TODO: add like generic empty stats summary if no stats available
  if (query.data === null) return null;

  return (
    <>
      <span className="text-sm text-muted-foreground">{`${formatSeconds(query.data.totalSecondsIngame, "hours minutes seconds")} played`}</span>
      <div className="flex flex-row gap-2">
        {[
          [
            query.data.avgPlacement.toPrecision(3),
            "Avg Place",
            `${query.data.compCount} games`,
          ],
          [
            `${query.data.top4Rate.toPrecision(3)}%`,
            "Top 4",
            `${query.data.top4Count} games`,
          ],
          [
            `${query.data.top1Rate.toPrecision(3)}%`,
            "Top 1",
            `${query.data.top1Count} games`,
          ],
        ].map(([stat, label, games]) => (
          <div
            key={label}
            className="flex flex-col items-center rounded bg-secondary p-1"
          >
            <span className="text-xl">{stat}</span>
            <span className="text-xs">{label}</span>
            <span className="text-xs text-muted-foreground">{games}</span>
          </div>
        ))}
      </div>
    </>
  );
};
