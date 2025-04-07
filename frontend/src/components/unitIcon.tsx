import { Unit } from "@/services/types";

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

const UnitIcon: React.FC<{ unit: Unit }> = ({ unit }) => {
  return (
    <div className="relative my-[0.625rem]">
      {unit.tier === 2 || unit.tier === 3 ? (
        <img
          className="absolute top-[-0.625rem]"
          src={`/stars/${unit.tier}.png`}
        />
      ) : null}
      <img
        className={`rounded border ${getRarityBackground(unit.rarity)}`}
        width={48}
        height={48}
        src={`/unit/${unit.characterId}.webp`}
      />
      <ul className="absolute bottom-[-0.625rem] flex w-full flex-row justify-center">
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
