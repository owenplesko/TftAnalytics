import TFT14_Divinicorp from "@/components/trait/TFT14_Divinicorp.svg?react";
import TFT14_StreetDemon from "@/components/trait/TFT14_StreetDemon.svg?react";
import TFT14_Supercharge from "@/components/trait/TFT14_Supercharge.svg?react";
import TFT14_AnimaSquad from "@/components/trait/TFT14_AnimaSquad.svg?react";
import TFT14_Vanguard from "@/components/trait/TFT14_Vanguard.svg?react";
import TFT14_Strong from "@/components/trait/TFT14_Strong.svg?react";
import TFT14_ViegoUniqueTrait from "@/components/trait/TFT14_ViegoUniqueTrait.svg?react";
import TFT14_Cyberboss from "@/components/trait/TFT14_Cyberboss.svg?react";
import TFT14_Cutter from "@/components/trait/TFT14_Cutter.svg?react";
import TFT14_Armorclad from "@/components/trait/TFT14_Armorclad.svg?react";
import TFT14_EdgeRunner from "@/components/trait/TFT14_EdgeRunner.svg?react";
import TFT14_Overlord from "@/components/trait/TFT14_Overlord.svg?react";
import TFT14_Netgod from "@/components/trait/TFT14_Netgod.svg?react";
import TFT14_Swift from "@/components/trait/TFT14_Swift.svg?react";
import TFT14_Controller from "@/components/trait/TFT14_Controller.svg?react";
import TFT14_Bruiser from "@/components/trait/TFT14_Bruiser.svg?react";
import TFT14_Marksman from "@/components/trait/TFT14_Marksman.svg?react";
import TFT14_Immortal from "@/components/trait/TFT14_Immortal.svg?react";
import TFT14_BallisTek from "@/components/trait/TFT14_BallisTek.svg?react";
import TFT14_Techie from "@/components/trait/TFT14_Techie.svg?react";
import TFT14_Thirsty from "@/components/trait/TFT14_Thirsty.svg?react";
import TFT14_Virus from "@/components/trait/TFT14_Virus.svg?react";
import TFT14_Mob from "@/components/trait/TFT14_Mob.svg?react";
import TFT14_HotRod from "@/components/trait/TFT14_HotRod.svg?react";
import TFT14_Suits from "@/components/trait/TFT14_Suits.svg?react";

const iconMap = new Map([
  ["TFT14_Divinicorp", TFT14_Divinicorp({})],
  ["TFT14_StreetDemon", TFT14_StreetDemon({})],
  ["TFT14_Supercharge", TFT14_Supercharge({})],
  ["TFT14_AnimaSquad", TFT14_AnimaSquad({})],
  ["TFT14_Vanguard", TFT14_Vanguard({})],
  ["TFT14_Strong", TFT14_Strong({})],
  ["TFT14_ViegoUniqueTrait", TFT14_ViegoUniqueTrait({})],
  ["TFT14_Cyberboss", TFT14_Cyberboss({})],
  ["TFT14_Cutter", TFT14_Cutter({})],
  ["TFT14_Armorclad", TFT14_Armorclad({})],
  ["TFT14_EdgeRunner", TFT14_EdgeRunner({})],
  ["TFT14_Overlord", TFT14_Overlord({})],
  ["TFT14_Netgod", TFT14_Netgod({})],
  ["TFT14_Swift", TFT14_Swift({})],
  ["TFT14_Controller", TFT14_Controller({})],
  ["TFT14_Bruiser", TFT14_Bruiser({})],
  ["TFT14_Marksman", TFT14_Marksman({})],
  ["TFT14_Immortal", TFT14_Immortal({})],
  ["TFT14_BallisTek", TFT14_BallisTek({})],
  ["TFT14_Techie", TFT14_Techie({})],
  ["TFT14_Thirsty", TFT14_Thirsty({})],
  ["TFT14_Virus", TFT14_Virus({})],
  ["TFT14_Mob", TFT14_Mob({})],
  ["TFT14_HotRod", TFT14_HotRod({})],
  ["TFT14_Suits", TFT14_Suits({})],
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
      className={`grid grid-cols-[0.75rem_auto] items-center gap-1 rounded border bg-opacity-10 px-1 text-xs ${
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
