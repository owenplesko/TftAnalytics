import data from "./unitData.json"

type UnitData = {
    cost: number
    rarity: number | null
    name: string
    role: string | null
    traits: string[]
}

export const unitData: Record<string, UnitData | undefined> = data