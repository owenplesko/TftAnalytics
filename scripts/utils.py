import aiohttp
import aiofiles
import os
import shutil
from typing import Any
from enum import Enum

class NotFoundException(Exception):
    pass

class DownloadResult(Enum):
    SUCCESS = 1
    SKIPPED = 2
    NOT_FOUND = 3
    ERROR = 4
    
class DownloadResultTracker():
    def __init__(self):
        self.state = {}
        for result in DownloadResult:
            self.state[result] = 0
            
    def push(self, res: DownloadResult):
        self.state[res] += 1
        
    def __str__(self):
        return f"✅{self.state[DownloadResult.SUCCESS]} ⏩{self.state[DownloadResult.SKIPPED]} ❓{self.state[DownloadResult.NOT_FOUND]} ❌{self.state[DownloadResult.ERROR]}"

async def file_exists(filename: str) -> bool:
    return os.path.exists(filename)

async def _fetch(session: aiohttp.ClientSession, url: str) -> aiohttp.ClientResponse:
    res = await session.get(url)
    if res.status == 404:
        raise NotFoundException
    if res.status != 200:
        raise Exception(f"failed to fetch: {url} (status {res.status})")
    return res

async def fetch_json(session: aiohttp.ClientSession, url: str) -> Any:
        res = await _fetch(session, url)
        data = await res.json()
        res.close()
        return data

async def fetch_bytes(session: aiohttp.ClientSession, url: str) -> bytes:
    res = await _fetch(session, url)
    data = await res.read()
    res.close()
    return data

def clear_directory(directory: str):
    for filename in os.listdir(directory):
        file_path = os.path.join(directory, filename)
        if os.path.isfile(file_path) or os.path.islink(file_path):
            os.remove(file_path)
        elif os.path.isdir(file_path):
            shutil.rmtree(file_path)


async def download_image(session: aiohttp.ClientSession, url: str, filename: str) -> str:
    if await file_exists(filename):
        return "SKIPPED"

    async with session.get(url) as res:
        if res.status == 404:
            return "NOT_FOUND"
        if res.status != 200:
            raise Exception(f"failed to download\n\turl={url}\n\tstatus={res.status}")

        os.makedirs(os.path.dirname(filename), exist_ok=True)
        data = await res.read()

        async with aiofiles.open(filename, "wb") as f:
            await f.write(data)

        return "SUCCESS"