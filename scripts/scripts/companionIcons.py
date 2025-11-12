import asyncio
import io
import os
from PIL import Image, ImageDraw
from typing import List
import aiohttp
from tqdm import tqdm

from datatypes import Companion
from constants import COMMUNITY_DRAGON_BASE_URL
from utils import DownloadResult, DownloadResultTracker, NotFoundException, fetch_bytes, file_exists

DIRECTORY_PATH = "assets/companion"

async def get_companions(session: aiohttp.ClientSession, companions: List[Companion]):
    os.makedirs(DIRECTORY_PATH, exist_ok=True)

    tracker = DownloadResultTracker()
    tasks = [asyncio.create_task(get_companion(session, companion)) for companion in companions]

    with tqdm(total=len(tasks), desc="companions") as pbar:
        for task in asyncio.as_completed(tasks):
            result = await task
            tracker.push(result)
            pbar.set_postfix_str(tracker)
            pbar.update(1)

async def get_companion(session: aiohttp.ClientSession, companion: Companion) -> DownloadResult:
    url = f"{COMMUNITY_DRAGON_BASE_URL}/plugins/rcp-be-lol-game-data/global/default/assets{companion.loadoutsIcon.split('ASSETS')[1].lower()}"
    filename = f"{DIRECTORY_PATH}/{companion.itemId}.webp"

    if await file_exists(filename):
        return DownloadResult.SKIPPED

    try:
        data = await fetch_bytes(session, url)
    except NotFoundException:
        return DownloadResult.NOT_FOUND
    except Exception:
        return DownloadResult.ERROR

    img = Image.open(io.BytesIO(data))
    img = crop_center_square(img)
    img = img.resize((128, 128))
    img.save(filename, format="WEBP", quality=80)

    return DownloadResult.SUCCESS


def crop_center_square(img: Image.Image) -> Image.Image:
    w, h = img.size
    diameter = h  # circle diameter = image height

    # Compute circle bounding box centered in image
    left = (w - diameter) // 2
    top = 0
    right = left + diameter
    bottom = top + diameter

    # Crop to square bounding box
    square = img.crop((left, top, right, bottom))

    return square
