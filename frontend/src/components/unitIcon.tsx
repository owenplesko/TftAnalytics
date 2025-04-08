import { Unit } from "@/services/types";
import { IconStarFilled } from "@tabler/icons-react";
import Repeat from "./util/repeat";
import { cn } from "@/lib/utils";

const rarityStyles: Record<number, string | undefined> = {
  0: "border-common",
  1: "border-uncommon",
  2: "border-rare",
  4: "border-mythic",
  6: "border-legendary",
  8: "border-divine",
};

const tierStyles: Record<number, string | undefined> = {
  2: "text-silver",
  3: "text-gold",
  4: "text-platinum",
};

const UnitIcon: React.FC<{ unit: Unit }> = ({ unit }) => {
  return (
    <div className="relative my-[0.75rem]">
      {unit.tier > 1 && (
        <div
          className={cn(
            "absolute top-[-0.75rem] flex w-full justify-center",
            tierStyles[unit.tier],
          )}
        >
          <Repeat n={unit.tier}>
            <IconStarFilled size={12} />
          </Repeat>
        </div>
      )}
      <img
        className={cn(
          "rounded border",
          rarityStyles[unit.rarity] ?? "border-muted",
        )}
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
