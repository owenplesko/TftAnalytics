import asyncio
import io
import os
from PIL import Image
from typing import Dict
import aiohttp
from tqdm import tqdm

from datatypes import Item, SetDatum
from constants import COMMUNITY_DRAGON_BASE_URL
from utils import DownloadResult, DownloadResultTracker, NotFoundException, fetch_bytes, file_exists

DIRECTORY_PATH = "assets/item"

async def get_item_icons(session: aiohttp.ClientSession, items_dict: Dict[str, Item], setData: SetDatum):
    os.makedirs(DIRECTORY_PATH, exist_ok=True)
    
    items = [i for i in setData.items if "_Item_" in i and "Debug" not in i and "Grant" not in i]
    
    tracker = DownloadResultTracker()
    tasks = [asyncio.create_task(get_item_icon(session, items_dict[item_name])) for item_name in items]

    with tqdm(total=len(tasks), desc=f"{setData.mutator} items") as pbar:
        for task in asyncio.as_completed(tasks):
            result = await task
            tracker.push(result)
            pbar.set_postfix_str(tracker)
            pbar.update(1)
            
async def get_item_icon(session: aiohttp, item: Item) -> DownloadResult:
    url = f"{COMMUNITY_DRAGON_BASE_URL}/game/{item.icon.replace('.tex', '.png').lower()}"
    filename = f"{DIRECTORY_PATH}/{item.apiName.lower()}.webp"
    
    if await file_exists(filename):
        return DownloadResult.SKIPPED

    try:
        data = await fetch_bytes(session, url)
    except NotFoundException:
        return DownloadResult.NOT_FOUND
    except Exception:
        return DownloadResult.ERROR

    img = Image.open(io.BytesIO(data))
    img = img.resize((64, 64))
    img.save(filename, format="WEBP", quality=80)

    return DownloadResult.SUCCESS