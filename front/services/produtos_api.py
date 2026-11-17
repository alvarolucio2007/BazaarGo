
import requests


class BazaarGoClient:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()

    def listar_produtos(self):
        response = self.session.get(
            f"{self.base_url}/products",
            timeout=10
        )
        response.raise_for_status()
        return response.json()

    def buscar_produtos(self, nome: str):
        response = self.session.get(
            f"{self.base_url}/products/search/{nome}",
            timeout=10
        )
        response.raise_for_status()
        return response.json()

    def buscar_produto(self, produto_id: str, client_id: str):
        response = self.session.get(
            f"{self.base_url}/products/{produto_id}",
            headers={"X-Client-ID": client_id},
            timeout=10
        )
        response.raise_for_status()
        return response.json()

    def produtos_recentes(self, client_id: str):
        response = self.session.get(
            f"{self.base_url}/products/recent",
            headers={"X-Client-ID": client_id},
            timeout=10
        )
        response.raise_for_status()
        return response.json()
