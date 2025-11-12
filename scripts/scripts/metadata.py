from typing import List
import aiohttp

from datatypes import Champion, Companion, Item, SetDatum, TftData, Trait
from constants import COMMUNITY_DRAGON_BASE_URL
from utils import fetch_json


async def get_tft_data(session: aiohttp.ClientSession) -> TftData:
    url = f"{COMMUNITY_DRAGON_BASE_URL}/cdragon/tft/en_us.json"
    data = await fetch_json(session, url)

    items = [Item(**item) for item in data["items"]]
    items_dict = {item.apiName: item for item in items}

    set_data = []
    for s in data["setData"]:
        champions = [Champion(**c) for c in s["champions"]]
        traits = [Trait(**t) for t in s["traits"]]
        set_data.append(SetDatum(
            number=s["number"],
            mutator=s["mutator"],
            champions=champions,
            items=s["items"],
            traits=traits
        ))

    return TftData(setData=set_data, items=items, itemsDict=items_dict)

async def get_companion_data(session: aiohttp.ClientSession) -> List[Companion]:
    url = f"{COMMUNITY_DRAGON_BASE_URL}/plugins/rcp-be-lol-game-data/global/default/v1/companions.json"
    data = await fetch_json(session, url)
    return [Companion(**c) for c in data]