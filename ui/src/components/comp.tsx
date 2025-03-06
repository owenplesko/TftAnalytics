import { SummonerMatch } from "@/services/types";
import UnitIcon from "./unitIcon";
import { formatSeconds, formatStageNumber, formatTimeSince } from "@/lib/utils";
import TraitIcon from "./trait";

const CompSummary: React.FC<{
  summonerMatch: SummonerMatch;
}> = ({ summonerMatch }) => {
  return (
    <>
      <div className="grid grid-cols-[70px_70px_50px_auto] items-center gap-4 rounded-sm border bg-neutral-900 p-3">
        <div className="flex flex-col">
          <span
            className={`text-xl font-bold ${
              // @ts-ignore
              // ignore because we have nullish coalesce
              placementColors[summonerMatch.compData.placement] ??
              "text-neutral-500"
            }`}
          >
            {summonerMatch.compData.placement}
          </span>
          <span className="text-sm text-neutral-300">
            {
              // @ts-ignore
              // ignore because we have nullish coalesce
              queueLabels[summonerMatch.tftMatch.queueId] ?? "Unknown"
            }
          </span>
          <span className="text-xs text-neutral-400">
            {formatTimeSince(summonerMatch.tftMatch.matchDate)}
          </span>
          <span className="text-xs text-neutral-400">
            {`${formatSeconds(summonerMatch.compData.timeEliminated)}
             • 
            ${formatStageNumber(summonerMatch.compData.lastRound)}`}
          </span>
        </div>
        <div className="relative">
          <img
            className="rounded-sm border-2"
            src={`/companion/${summonerMatch.compData.companion.itemId}.png`}
          />
          <span className="absolute bottom-[-10px] right-1/2 h-5 w-5 translate-x-1/2 rounded-sm bg-neutral-800 text-center text-xs leading-5 text-neutral-300">
            {summonerMatch.compData.level}
          </span>
        </div>
        <div className="grid grid-cols-[16px_auto] items-center gap-x-2 text-neutral-300">
          <img src="/icon/combat.png" />
          <span>{summonerMatch.compData.totalDamageToPlayers}</span>
          <img src="/icon/gold.png" />
          <span>{summonerMatch.compData.goldLeft}</span>
        </div>
        <div className="flex flex-col gap-3">
          <ul className="flex flex-row gap-2">
            {summonerMatch.compData.traits
              .sort((a, b) => b.style - a.style)
              .map(
                (trait) =>
                  trait.style > 0 && (
                    <li>
                      <TraitIcon
                        name={trait.name}
                        style={trait.style}
                        tier={trait.tierCurrent}
                      />
                    </li>
                  ),
              )}
          </ul>
          <ul className="flex flex-row gap-2 py-2">
            {summonerMatch.compData.units.map((unit) => (
              <li>
                <UnitIcon unit={unit} />
              </li>
            ))}
          </ul>
        </div>
      </div>
    </>
  );
};
export default CompSummary;

const queueLabels = {
  1090: "Normal",
  1100: "Ranked",
  1110: "Tutorial",
  1130: "Hyper Roll",
  1160: "Double Up",
  1220: "Tocker's Trials",
  6100: "Revival",
};

const placementColors = {
  1: "text-yellow-400",
  2: "text-slate-200",
  3: "text-yellow-800",
};
