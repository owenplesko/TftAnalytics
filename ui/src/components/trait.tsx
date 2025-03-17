import TFT13_Academy from "@/components/traits/TFT13_Academy.svg?react";
import TFT13_Ambassador from "@/components/traits/TFT13_Ambassador.svg?react";
import TFT13_Ambusher from "@/components/traits/TFT13_Ambusher.svg?react";
import TFT13_Bruiser from "@/components/traits/TFT13_Bruiser.svg?react";
import TFT13_Cabal from "@/components/traits/TFT13_Cabal.svg?react";
import TFT13_Challenger from "@/components/traits/TFT13_Challenger.svg?react";
import TFT13_Crime from "@/components/traits/TFT13_Crime.svg?react";
import TFT13_Experiment from "@/components/traits/TFT13_Experiment.svg?react";
import TFT13_Family from "@/components/traits/TFT13_Family.svg?react";
import TFT13_FormSwapper from "@/components/traits/TFT13_FormSwapper.svg?react";
import TFT13_Hextech from "@/components/traits/TFT13_Hextech.svg?react";
import TFT13_HighRoller from "@/components/traits/TFT13_HighRoller.svg?react";
import TFT13_Hoverboard from "@/components/traits/TFT13_Hoverboard.svg?react";
import TFT13_Infused from "@/components/traits/TFT13_Infused.svg?react";
import TFT13_Invoker from "@/components/traits/TFT13_Invoker.svg?react";
import TFT13_JunkerKing from "@/components/traits/TFT13_JunkerKing.svg?react";
import TFT13_Martialist from "@/components/traits/TFT13_Martialist.svg?react";
import TFT13_Pugilist from "@/components/traits/TFT13_Pugilist.svg?react";
import TFT13_Rebel from "@/components/traits/TFT13_Rebel.svg?react";
import TFT13_Scrap from "@/components/traits/TFT13_Scrap.svg?react";
import TFT13_Sniper from "@/components/traits/TFT13_Sniper.svg?react";
import TFT13_Sorcerer from "@/components/traits/TFT13_Sorcerer.svg?react";
import TFT13_Squad from "@/components/traits/TFT13_Squad.svg?react";
import TFT13_Titan from "@/components/traits/TFT13_Titan.svg?react";
import TFT13_Warband from "@/components/traits/TFT13_Warband.svg?react";
import TFT13_Watcher from "@/components/traits/TFT13_Watcher.svg?react";
import TFT13_MissMageTrait from "@/components/traits/TFT13_MissMageTrait.svg?react";
import TFT13_MachineHerald from "@/components/traits/TFT13_MachineHerald.svg?react";
import TFT13_BloodHunter from "@/components/traits/TFT13_BloodHunter.svg?react";

const iconMap = new Map([
  ["TFT13_Academy", TFT13_Academy({})],
  ["TFT13_Ambassador", TFT13_Ambassador({})],
  ["TFT13_Ambusher", TFT13_Ambusher({})],
  ["TFT13_Bruiser", TFT13_Bruiser({})],
  ["TFT13_Cabal", TFT13_Cabal({})],
  ["TFT13_Challenger", TFT13_Challenger({})],
  ["TFT13_Crime", TFT13_Crime({})],
  ["TFT13_Experiment", TFT13_Experiment({})],
  ["TFT13_Family", TFT13_Family({})],
  ["TFT13_FormSwapper", TFT13_FormSwapper({})],
  ["TFT13_Hextech", TFT13_Hextech({})],
  ["TFT13_HighRoller", TFT13_HighRoller({})],
  ["TFT13_Hoverboard", TFT13_Hoverboard({})],
  ["TFT13_Infused", TFT13_Infused({})],
  ["TFT13_Invoker", TFT13_Invoker({})],
  ["TFT13_JunkerKing", TFT13_JunkerKing({})],
  ["TFT13_Martialist", TFT13_Martialist({})],
  ["TFT13_Pugilist", TFT13_Pugilist({})],
  ["TFT13_Rebel", TFT13_Rebel({})],
  ["TFT13_Scrap", TFT13_Scrap({})],
  ["TFT13_Sniper", TFT13_Sniper({})],
  ["TFT13_Sorcerer", TFT13_Sorcerer({})],
  ["TFT13_Squad", TFT13_Squad({})],
  ["TFT13_Titan", TFT13_Titan({})],
  ["TFT13_Warband", TFT13_Warband({})],
  ["TFT13_Watcher", TFT13_Watcher({})],
  ["TFT13_MissMageTrait", TFT13_MissMageTrait({})],
  ["TFT13_MachineHerald", TFT13_MachineHerald({})],
  ["TFT13_BloodHunter", TFT13_BloodHunter({})],
]);

const styleColors = {
  1: "text-bronze bg-bronze border-bronze",
  2: "text-silver bg-silver border-silver",
  3: "text-unique bg-unique border-unique",
  4: "text-gold bg-gold border-gold",
  5: "text-prismatic bg-prismatic border-prismatic",
};

const TraitIcon: React.FC<{
  name: string;
  numUnits: number;
  style: number;
}> = ({ name, numUnits: tier, style }) => {
  const icon = iconMap.get(name);
  return (
    <div
      className={`grid grid-cols-[12px_auto] items-center gap-1 rounded-sm border bg-opacity-15 px-1 text-xs ${
        // @ts-ignore
        // ignore because im a beast
        styleColors[style]
      }`}
    >
      {icon}
      <span>{tier}</span>
    </div>
  );
};

export default TraitIcon;
