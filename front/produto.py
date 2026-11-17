from dataclasses import dataclass
from datetime import time
from typing import Any


@dataclass
class Produto:
    id:str
    name:str
    description:str
    price:float
    stock:int
    categoryID:str
    specs:dict[str,Any]
    created_at:time
    updated_at:time
