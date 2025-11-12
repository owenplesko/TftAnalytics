import asyncio
import aiohttp

from scripts.itemIcons import get_item_icons
from scripts.companionIcons import get_companions
from scripts.metadata import get_companion_data, get_tft_data
from scripts.traitIcons import get_trait_icons
from scripts.unitData import save_unit_data
from scripts.unitIcons import get_unit_icons
from utils import clear_directory

TARGET_SETS = ["TFTSet15"]

async def main():
    clear_directory("./assets")
    
    async with aiohttp.ClientSession() as session:
        tft_data = await get_tft_data(session)
        companions = await get_companion_data(session)

        target_sets = [s for s in tft_data.setData if s.mutator in TARGET_SETS]

        #for set in target_sets:
            #save_unit_data(set)
            #await get_unit_icons(session, set)
            #await get_item_icons(session, tft_data.itemsDict, set)
            #await get_trait_icons(session, set)

        await get_companions(session, companions)

if __name__ == "__main__":
    asyncio.run(main())
