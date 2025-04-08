import { Unit } from "@/services/types";
import { IconStarFilled } from "@tabler/icons-react";
import Repeat from "./util/repeat";
import { cn } from "@/lib/utils";

const getRarityBackground = (rarity: number) => {
  switch (rarity) {
    case 0:
      return "border-common";
    case 1:
      return "border-uncommon";
    case 2:
      return "border-rare";
    case 4:
      return "border-mythic";
    case 6:
      return "border-legendary";
    case 8:
      return "border-divine";
  }
  return "border-muted";
};

const starColors: Record<number, string | undefined> = {
  2: "text-silver",
  3: "text-gold",
  4: "text-blue-400",
};

const UnitIcon: React.FC<{ unit: Unit }> = ({ unit }) => {
  return (
    <div className="relative my-[0.75rem]">
      {unit.tier >= 1 && (
        <div
          className={cn(
            "absolute top-[-0.75rem] flex w-full justify-center",
            starColors[unit.tier],
          )}
        >
          <Repeat n={unit.tier}>
            <IconStarFilled size={12} />
          </Repeat>
        </div>
      )}
      <img
        className={`rounded border ${getRarityBackground(unit.rarity)}`}
        width={48}
        height={48}
        src={`/unit/${unit.characterId}.webp`}
      />
      <ul className="absolute bottom-[-0.75rem] flex w-full flex-row justify-center">
        {unit.itemNames.map((item, i) => (
          // key = index is acceptable because data is static
          <li key={i}>
            <img
              className="rounded border"
              width={16}
              height={16}
              src={`/item/${item}.webp`}
            />
          </li>
        ))}
      </ul>
    </div>
  );
};
export default UnitIcon;
