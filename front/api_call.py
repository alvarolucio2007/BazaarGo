from typing import Any

import requests
from produto import Produto, UsuarioResponse


#codigo gerado por IA.
class BazaarGoClient:
    def __init__(self, base_url: str):
        """
        Cliente para a API BazaarGo.
        :param base_url: URL base da API (ex: 'http://localhost:8080')
        """
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.token: str | None = None

    def _get_headers(self, client_id: str | None = None) -> dict[str, str]:
        headers = {"Content-Type": "application/json"}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        if client_id:
            headers["X-Client-ID"] = client_id
        return headers

    def health_check(self) -> dict[str, bool]:
        """Verifica a saúde da API (Postgres, Mongo, Valkey)."""
        url = f"{self.base_url}/health"
        response = self.session.get(url, headers=self._get_headers())
        response.raise_for_status()
        return response.json()

    def register_user(self, name: str, email: str, password: str, cpf: str) -> UsuarioResponse:
        """Registra um novo usuário e retorna o objeto tipado."""
        url = f"{self.base_url}/auth/register"
        payload = {"name": name, "email": email, "password": password, "cpf": cpf}
        response = self.session.post(url, json=payload, headers=self._get_headers())
        response.raise_for_status()
        data = response.json()
        return UsuarioResponse(**data)

    def create_token(self, email: str, password: str) -> str:
        """Gera e armazena o token de autenticação (PASETO)."""
        url = f"{self.base_url}/auth/token"
        payload = {"email": email, "password": password}
        response = self.session.post(url, json=payload, headers=self._get_headers())
        response.raise_for_status()
        data = response.json()
        self.token = data.get("token")
        return self.token  # pyright: ignore[reportReturnType]

    def create_product(self, name: str, description: str, price: float, stock: int, category_id: str, specs: dict[str, Any]) -> Produto:
        """Cria um novo produto e retorna a instância de Produto."""
        url = f"{self.base_url}/products"
        payload = {
            "name": name,
            "description": description,
            "price": price,
            "stock": stock,
            "category_id": category_id,
            "specs": specs
        }
        response = self.session.post(url, json=payload, headers=self._get_headers())
        response.raise_for_status()
        return Produto(**response.json())

    def list_products(self) -> list[Produto]:
        """Lista todos os produtos cadastrados."""
        url = f"{self.base_url}/products"
        response = self.session.get(url, headers=self._get_headers())
        response.raise_for_status()

        data = response.json()

        items = data.get("data", data) if isinstance(data, dict) else data

        return [Produto(**item) for item in items if isinstance(item, dict)]

    def get_product_by_id(self, product_id: str, client_id: str) -> Produto:
        """Busca um produto por ID (Requer cabeçalho X-Client-ID)."""
        url = f"{self.base_url}/products/{product_id}"
        headers = self._get_headers(client_id=client_id)
        response = self.session.get(url, headers=headers)
        response.raise_for_status()
        return Produto(**response.json())

    def update_product(self, product_id: str, name: str, description: str, price: float, stock: int, category_id: str, specs: dict[str, Any]) -> None:
        """Atualiza um produto por ID (Retorna 204 No Content em caso de sucesso)."""
        url = f"{self.base_url}/products/{product_id}"
        payload = {
            "name": name,
            "description": description,
            "price": price,
            "stock": stock,
            "category_id": category_id,
            "specs": specs
        }
        response = self.session.put(url, json=payload, headers=self._get_headers())
        response.raise_for_status()

    def delete_product(self, product_id: str) -> None:
        """Remove um produto por ID."""
        url = f"{self.base_url}/products/{product_id}"
        response = self.session.delete(url, headers=self._get_headers())
        response.raise_for_status()

    def search_products(self, name: str) -> list[Produto]:
        """Busca produtos pelo nome."""
        url = f"{self.base_url}/products/search/{name}"
        response = self.session.get(url, headers=self._get_headers())
        response.raise_for_status()

        data = response.json()
        if isinstance(data, list) and len(data) > 0 and isinstance(data[0], list):
                    items = data[0]
        elif isinstance(data, dict):
            items = data.get("data", data)
        else:
            items = data
        return [Produto(**item) for item in items if isinstance(item, dict)]

    def get_recent_products(self, client_id: str) -> list[Produto]:
        """Lista produtos recentes (Requer X-Client-ID)."""
        url = f"{self.base_url}/products/recent"
        headers = self._get_headers(client_id=client_id)
        response = self.session.get(url, headers=headers)
        response.raise_for_status()

        data = response.json()
        items = data.get("data", data) if isinstance(data, dict) else data

        return [Produto(**item) for item in items if isinstance(item, dict)]

    def checkout(self, phone: str, items: list[dict[str, Any]]) -> str:
        """Realiza o checkout e retorna o link do WhatsApp gerado."""
        url = f"{self.base_url}/checkout"
        payload = {"phone": phone, "items": items}
        response = self.session.post(url, json=payload, headers=self._get_headers())
        response.raise_for_status()
        return response.json().get("whatsapp_link", "")
