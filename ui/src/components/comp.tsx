import { MatchComp, type CompData, type SummonerMatch } from "@/services/types";
import UnitIcon from "./unitIcon";
import { formatSeconds, formatStageNumber, formatTimeSince } from "@/lib/utils";
import TraitIcon from "./trait";
import { useState } from "react";
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { getMatchComps } from "@/services/getMatchComps";
import { IconChevronUp } from "@tabler/icons-react";

const SummonerMatchCard: React.FC<{
  summonerMatch: SummonerMatch;
}> = ({ summonerMatch }) => {
  const [expanded, setExpanded] = useState(false);
  const queryClient = useQueryClient();

  return (
    <>
      <div className="overflow-hidden rounded-sm border">
        <div className="flex items-stretch">
          <div className="grid flex-grow grid-cols-[75px_75px_50px_auto] items-center gap-4 bg-muted/50 p-3">
            <PlacementSection summonerMatch={summonerMatch} />
            <TacticianSection compData={summonerMatch.compData} />
            <CombatStatSection compData={summonerMatch.compData} />
            <UnitSection compData={summonerMatch.compData} />
          </div>
          <button
            className="flex items-end justify-center bg-secondary hover:bg-secondary/80"
            onClick={async () => {
              if (!expanded)
                await queryClient.ensureQueryData(
                  getMatchComps(summonerMatch.tftMatch.id),
                );

              setExpanded(!expanded);
            }}
          >
            <IconChevronUp className={expanded ? "rotate-180" : "rotate-0"} />
          </button>
        </div>
        {expanded && (
          <MatchCard
            matchComps={
              queryClient.getQueryData(
                getMatchComps(summonerMatch.tftMatch.id).queryKey,
              )!
              // justification for using non-null assertion:
              // IF expanded is true THEN ensureQueryData has ran MEANING getQueryData will not be undefined
            }
          />
        )}
      </div>
    </>
  );
};
export default SummonerMatchCard;

const MatchCard: React.FC<{
  matchComps: MatchComp[];
}> = ({ matchComps }) => {
  return (
    <ul className="bg-muted/25">
      {matchComps
        .sort((a, b) => a.compData.placement - b.compData.placement)
        .map(({ tftSummoner, compData }) => (
          <li className="grid grid-cols-[25px_50px_75px_50px_auto] items-center gap-4 border-t px-3 py-2">
            <PlacementText placement={compData.placement} />
            <TacticianSection compData={compData} />
            <div
              className="overflow-hidden overflow-ellipsis whitespace-nowrap text-left text-sm"
              title={`${tftSummoner.name}#${tftSummoner.tag}`}
            >
              <span>{tftSummoner.name}</span>
              <span className="text-neutral-400">#{tftSummoner.tag}</span>
            </div>
            <CombatStatSection compData={compData} />
            <UnitSection compData={compData} />
          </li>
        ))}
    </ul>
  );
};

const placementColor = (placement: number) => {
  switch (placement) {
    case 1:
      return "text-yellow-400";
    case 2:
      return "text-slate-200";
    case 3:
      return "text-yellow-800";
    case 4:
      return "text-neutral-400";
    default:
      return "text-neutral-500";
  }
};

const PlacementText: React.FC<{
  placement: number;
}> = ({ placement }) => {
  return (
    <span className={`text-xl font-bold ${placementColor(placement)}`}>
      {placement}
    </span>
  );
};

const PlacementSection: React.FC<{ summonerMatch: SummonerMatch }> = ({
  summonerMatch,
}) => {
  return (
    <div className="flex flex-col">
      <PlacementText placement={summonerMatch.compData.placement} />
      <span className="text-sm text-neutral-300">
        {
          // @ts-ignore
          // ignore because we have nullish coalesce
          queueLabels[summonerMatch.tftMatch.queueId] ?? "Unknown"
        }
      </span>
      <span className="text-xs text-neutral-400">
        {formatTimeSince(new Date(summonerMatch.tftMatch.matchDate))}
      </span>
      <span className="text-xs text-neutral-400">
        {`${formatSeconds(summonerMatch.compData.timeEliminated)}
   • 
  ${formatStageNumber(summonerMatch.compData.lastRound)}`}
      </span>
    </div>
  );
};

const CombatStatSection: React.FC<{ compData: CompData }> = ({ compData }) => {
  return (
    <div className="grid grid-cols-[16px_auto] items-center gap-x-2 text-neutral-300">
      <img src="/icon/combat.png" />
      <span>{compData.totalDamageToPlayers}</span>
      <img src="/icon/gold.png" />
      <span>{compData.goldLeft}</span>
    </div>
  );
};

const TacticianSection: React.FC<{ compData: CompData }> = ({ compData }) => {
  return (
    <div className="relative">
      <img
        className="rounded-sm border-2"
        src={`/companion/${compData.companion.itemId}.png`}
      />
      <span className="absolute bottom-[-10px] right-1/2 h-5 w-5 translate-x-1/2 rounded-sm bg-neutral-800 text-center text-xs leading-5 text-neutral-300">
        {compData.level}
      </span>
    </div>
  );
};

const UnitSection: React.FC<{ compData: CompData }> = ({ compData }) => {
  return (
    <div className="flex flex-col gap-1">
      <ul className="flex flex-row gap-2">
        {compData.traits
          .sort((a, b) => b.style - a.style)
          .map(
            (trait, i) =>
              trait.style > 0 && (
                // key = index is acceptable because data is static
                <li key={i}>
                  <TraitIcon
                    name={trait.name}
                    style={trait.style}
                    numUnits={trait.numUnits}
                  />
                </li>
              ),
          )}
      </ul>
      <ul className="flex flex-row gap-2 py-2">
        {compData.units.map((unit, i) => (
          // key = index is acceptable because data is static
          <li key={i}>
            <UnitIcon unit={unit} />
          </li>
        ))}
      </ul>
    </div>
  );
};

const queueLabels = {
  1090: "Normal",
  1100: "Ranked",
  1110: "Tutorial",
  1130: "Hyper Roll",
  1160: "Double Up",
  1220: "Tocker's Trials",
  6100: "Revival",
};
