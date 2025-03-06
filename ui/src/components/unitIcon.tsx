import { Unit } from "@/services/types";

const getRarityBackground = (rarity: number) => {
  switch (rarity) {
    case 0:
      return "border-neutral-400";
    case 1:
      return "border-green-500";
    case 2:
      return "border-blue-500";
    case 4:
      return "border-purple-500";
    case 6:
      return "border-amber-400";
    case 8:
      return "border-cyan-300";
  }
  return "border-neutral-800";
};

const UnitIcon: React.FC<{ unit: Unit }> = ({ unit }) => {
  return (
    <>
      <div className="relative">
        {unit.tier === 2 || unit.tier === 3 ? (
          <img
            className="absolute top-[-10px]"
            src={`/stars/${unit.tier}.png`}
          />
        ) : null}
        <img
          className={`rounded-sm border-2 ${getRarityBackground(unit.rarity)}`}
          width={48}
          height={48}
          src={`/unit/${unit.characterId}.png`}
        />
        <ul className="absolute bottom-[-10px] flex w-full flex-row justify-center">
          {unit.itemNames.map((item) => (
            <li>
              <img
                className="rounded-sm border border-neutral-800"
                width={16}
                height={16}
                src={`/item/${item}.png`}
              />
            </li>
          ))}
        </ul>
      </div>
    </>
  );
};
export default UnitIcon;
