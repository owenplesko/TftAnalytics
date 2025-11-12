import asyncio
from io import BytesIO
import os
import aiohttp
from tqdm import tqdm
from datatypes import SetDatum, Trait
import xml.etree.ElementTree as ET
from constants import TACTICS_TOOLS_CDN_BASE_URL
from utils import DownloadResult, DownloadResultTracker, NotFoundException, fetch_bytes, file_exists

DIRECTORY_PATH = "assets/trait"

async def get_trait_icons(session: aiohttp.ClientSession, setData: SetDatum):
    os.makedirs(DIRECTORY_PATH, exist_ok=True)
    
    tracker = DownloadResultTracker()
    tasks = [asyncio.create_task(get_trait_icon(session, trait)) for trait in setData.traits]
    
    with tqdm(total=len(tasks), desc=f"{setData.mutator} traits") as pbar:
        for task in asyncio.as_completed(tasks):
            result = await task
            tracker.push(result)
            pbar.set_postfix_str(tracker)
            pbar.update(1)
            
async def get_trait_icon(session: aiohttp.ClientSession, trait: Trait) -> DownloadResult:
    url = f"{TACTICS_TOOLS_CDN_BASE_URL}/trait-icons/{trait.apiName.lower()}.svg"
    filename = f"{DIRECTORY_PATH}/{trait.apiName.lower()}.svg"
    
    if await file_exists(filename):
        return DownloadResult.SKIPPED

    try:
        data = await fetch_bytes(session, url)
    except NotFoundException:
        return DownloadResult.NOT_FOUND
    except Exception:
        return DownloadResult.ERROR

    svg = ET.parse(BytesIO(data))
    setFillColor(svg)
    svg.write(filename)

    return DownloadResult.SUCCESS

def setFillColor(tree):
    root = tree.getroot()
    ET.register_namespace("", "http://www.w3.org/2000/svg")
    for el in root.iter():
        # Only apply to SVG drawable elements
        if el.tag.endswith(("path", "rect", "circle", "polygon", "ellipse", "line", "polyline")):
            el.set("fill", "currentColor")