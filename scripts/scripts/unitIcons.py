import asyncio
import io
import os
from PIL import Image
import aiohttp
from tqdm import tqdm

from datatypes import Champion, SetDatum
from constants import COMMUNITY_DRAGON_BASE_URL
from utils import DownloadResult, DownloadResultTracker, NotFoundException, fetch_bytes, file_exists

DIRECTORY_PATH = "assets/unit"

async def get_unit_icons(session: aiohttp.ClientSession, setData: SetDatum):
    os.makedirs(DIRECTORY_PATH, exist_ok=True)
    
    tracker = DownloadResultTracker()
    tasks = [asyncio.create_task(get_unit_icon(session, unit)) for unit in setData.champions]
    
    with tqdm(total=len(tasks), desc=f"{setData.mutator} units") as pbar:
        for task in asyncio.as_completed(tasks):
            result = await task
            tracker.push(result)
            pbar.set_postfix_str(tracker)
            pbar.update(1)
            
async def get_unit_icon(session: aiohttp.ClientSession, unit: Champion) -> DownloadResult:
    if(unit.squareIcon == None):
        return DownloadResult.SKIPPED
    
    url = f"{COMMUNITY_DRAGON_BASE_URL}/game/{unit.squareIcon.replace('.tex', '.png').lower()}"
    filename = f"{DIRECTORY_PATH}/{unit.apiName.lower()}.webp"
    
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