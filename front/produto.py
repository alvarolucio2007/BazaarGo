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
    category_id:str
    specs:dict[str,Any]

@dataclass
class UsuarioResponse:
    id: int
    name: str
    cpf: str
    email: str
