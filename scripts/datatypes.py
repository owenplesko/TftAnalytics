from typing import Any, Dict, List, Optional
from pydantic import BaseModel

class Item(BaseModel):
    apiName: str
    icon: str

    class Config:
        extra = "ignore"

class Champion(BaseModel):
    ability: Any
    apiName: str
    characterName: str
    cost: int
    icon: Optional[str]
    name: Optional[str]
    role: Optional[str]
    squareIcon: Optional[str]
    stats: Any
    tileIcon: Optional[str]
    traits: List[str]


class Trait(BaseModel):
    apiName: str

    class Config:
        extra = "ignore"

class SetDatum(BaseModel):
    number: int
    mutator: str
    champions: List[Champion]
    items: List[str]
    traits: List[Trait]

    class Config:
        extra = "ignore"


class Companion(BaseModel):
    name: str
    itemId: int
    loadoutsIcon: str

    class Config:
        extra = "ignore"


class TftData(BaseModel):
    setData: List[SetDatum]
    items: List[Item]

    class Config:
        extra = "ignore"

    itemsDict: Dict[str, Item] = {}
