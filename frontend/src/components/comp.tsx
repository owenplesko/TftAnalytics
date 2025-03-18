import { MatchComp, type CompData, type SummonerMatch } from "@/services/types";
import UnitIcon from "./unitIcon";
import { formatSeconds, formatStageNumber, formatTimeSince } from "@/lib/utils";
import TraitIcon from "./trait";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { getMatchComps } from "@/services/getMatchComps";
import { IconChevronUp } from "@tabler/icons-react";
import { Link } from "@tanstack/react-router";

const SummonerMatchCard: React.FC<{
  summonerMatch: SummonerMatch;
}> = ({ summonerMatch }) => {
  const [expanded, setExpanded] = useState(false);
  const queryClient = useQueryClient();

  return (
    <>
      <div className="overflow-hidden rounded border text-card-foreground">
        <div className="flex">
          <div className="grid flex-grow grid-cols-[75px_75px_50px_auto] items-center gap-2 bg-card p-4">
            <PlacementSection summonerMatch={summonerMatch} />
            <TacticianSection compData={summonerMatch.compData} />
            <CombatStatSection compData={summonerMatch.compData} />
            <UnitSection compData={summonerMatch.compData} />
          </div>
          <button
            className="flex items-end border-l bg-secondary text-secondary-foreground hover:bg-secondary/80"
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
    <ul className="bg-card/80">
      {matchComps
        .sort((a, b) => a.compData.placement - b.compData.placement)
        .map(({ tftSummoner, compData }) => (
          <li
            key={tftSummoner.puuid}
            className="grid grid-cols-[25px_50px_75px_50px_auto] items-center gap-2 border-t px-4 py-2"
          >
            <PlacementText placement={compData.placement} />
            <TacticianSection compData={compData} />
            <Link
              className="overflow-hidden overflow-ellipsis whitespace-nowrap text-left text-sm hover:opacity-75"
              title={`${tftSummoner.name}#${tftSummoner.tag}`}
              to="/player/$name/$tag"
              params={{ name: tftSummoner.name, tag: tftSummoner.tag }}
            >
              <span>{tftSummoner.name}</span>
              <span className="text-muted-foreground">#{tftSummoner.tag}</span>
            </Link>
            <CombatStatSection compData={compData} />
            <UnitSection compData={compData} />
          </li>
        ))}
    </ul>
  );
};

const PlacementText: React.FC<{
  placement: number;
}> = ({ placement }) => {
  return (
    <span
      className={`text-2xl font-bold ${{ 1: "text-gold", 2: "text-silver", 3: "text-bronze" }[placement] ?? "text-muted-foreground"}`}
    >
      {placement}
    </span>
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

const PlacementSection: React.FC<{ summonerMatch: SummonerMatch }> = ({
  summonerMatch,
}) => {
  return (
    <div className="flex flex-col">
      <PlacementText placement={summonerMatch.compData.placement} />
      <span>
        {
          // @ts-ignore
          // ignore because we have nullish coalesce
          queueLabels[summonerMatch.tftMatch.queueId] ?? "Unknown"
        }
      </span>
      <span className="text-xs text-muted-foreground">
        {formatTimeSince(new Date(summonerMatch.tftMatch.matchDate))}
      </span>
      <span className="text-xs text-muted-foreground">
        {`${formatSeconds(summonerMatch.compData.timeEliminated)}
          •
          ${formatStageNumber(summonerMatch.compData.lastRound)}`}
      </span>
    </div>
  );
};

const CombatStatSection: React.FC<{ compData: CompData }> = ({ compData }) => {
  return (
    <div className="grid grid-cols-[1rem_auto] items-center gap-x-1 text-neutral-300">
      <img src="/icon/combat.png" />
      {compData.totalDamageToPlayers}
      <img src="/icon/gold.png" />
      {compData.goldLeft}
    </div>
  );
};

const TacticianSection: React.FC<{ compData: CompData }> = ({ compData }) => {
  return (
    <div className="relative">
      <img
        className="rounded border"
        src={`/companion/${compData.companion.itemId}.png`}
      />
      <span className="absolute bottom-[-0.75rem] right-1/2 h-5 w-5 translate-x-1/2 rounded bg-secondary text-center text-xs leading-5 text-secondary-foreground">
        {compData.level}
      </span>
    </div>
  );
};

const UnitSection: React.FC<{ compData: CompData }> = ({ compData }) => {
  return (
    <div className="flex flex-col gap-2">
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
      <ul className="flex flex-row gap-2">
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
