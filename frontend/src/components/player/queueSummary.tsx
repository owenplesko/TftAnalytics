import { sentenceCase } from "@/lib/utils";
import { getSummonerRank } from "@/services/getSummonerRank";
import { getSummonerStats } from "@/services/getSummonerStats";
import { useSuspenseQuery } from "@tanstack/react-query";

const QueueSummary: React.FC<{
  region: string;
  summonerId: string;
  puuid: string;
  queueId?: number | null;
  setNumber: number;
}> = ({ region, summonerId, puuid, queueId, setNumber }) => {
  return (
    <div className="flex flex-col rounded border bg-card p-4">
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
    <div className="flex flex-row gap-2">
      <img
        width={24}
        height={24}
        src={`/rank/${query.data.rankData.tier.toLowerCase()}.svg`}
      />
      <span>{`${sentenceCase(query.data.rankData.tier)}`}</span>
      <span>{`${query.data.rankData.leaguePoints} LP`}</span>
      <span className="text-muted-foreground">{`Rank #${query.data.leaderboardPosition.position}`}</span>
      <span className="text-muted-foreground">{`Top ${query.data.leaderboardPosition.top.toPrecision(2)}%`}</span>
    </div>
  );
};

const QueueStatsSummary: React.FC<{
  puuid: string;
  queueId?: number | null;
  setNumber: number;
}> = ({ puuid, queueId, setNumber }) => {
  const query = useSuspenseQuery(getSummonerStats(puuid, setNumber, queueId));

  // TODO: add like generic empty stats summary if no stats available
  if (query.data === null) return null;

  return (
    <div className="flex flex-col">
      <span>{`${query.data.compCount} Games Played`}</span>
      <span>{`${query.data.avgPlacement.toPrecision(3)} AVP`}</span>
      <span>{`${query.data.top4Rate.toPrecision(3)}% Top 4`}</span>
      <span>{`${query.data.top1Rate.toPrecision(3)}% Top 1`}</span>
    </div>
  );
};
