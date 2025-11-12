import json
from datatypes import SetDatum

def cost_to_rarity(cost):
    match cost:
        case 1:
            return 0
        case 2:
            return 1
        case 3:
            return 2
        case 4:
            return 4
        case 5:
            return 6
        case _:
            return None
            

def save_unit_data(setData: SetDatum):
    data = {}

    for c in setData.champions:
        # Start with only the desired fields
        champ_data = c.model_dump(include={"cost", "name", "role", "traits"})

        champ_data["rarity"] = cost_to_rarity(c.cost)            

        data[c.apiName.lower()] = champ_data

    with open("assets/unitData.json", "w") as f:
        json.dump(data, f, indent=2)
