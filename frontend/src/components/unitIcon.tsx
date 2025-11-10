import { unitData } from "../data/unitData";
import { IconStarFilled } from "@tabler/icons-react";
import Repeat from "./util/repeat";
import { cn } from "@/lib/utils";

function rarityStyle(rarity: number | null) {
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
    default:
      return "border-border";
  }
}

const starLevelStyles: Record<number, string | undefined> = {
  2: "text-silver",
  3: "text-gold",
  4: "text-platinum",
};

const UnitIcon: React.FC<{
  apiName: string;
  items?: string[];
  starLevel?: number;
}> = ({ apiName, items = [], starLevel = 0 }) => {
  const data = unitData[apiName];

  if (!data) return <span>{apiName}</span>;

  return (
    <div className="relative my-[0.75rem] border-muted">
      {starLevel > 1 && (
        <div
          className={cn(
            "absolute top-[-0.75rem] flex w-full justify-center",
            starLevelStyles[starLevel],
          )}
        >
          <Repeat n={starLevel}>
            <IconStarFilled size={12} />
          </Repeat>
        </div>
      )}
      <img
        className={cn("rounded border", rarityStyle(data.rarity))}
        width={48}
        height={48}
        src={`/assets/unit/${apiName}.webp`}
      />
      <ul className="absolute bottom-[-0.75rem] flex w-full flex-row justify-center">
        {items.map((item, i) => (
          <li key={i}>
            <img
              className="rounded border"
              width={16}
              height={16}
              src={`/assets/item/${item.toLowerCase()}.webp`}
            />
          </li>
        ))}
      </ul>
    </div>
  );
};
export default UnitIcon;
