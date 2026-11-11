from typing import Any

import requests

#codigo gerado por IA.
class BazaarGoClient:
    def __init__(self, base_url: str):
        """
        Cliente para a API BazaarGo.
        :url: URL base da API (ex: 'https://api.bazaargo.com')
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

    # ==========================================
    # HEALTH CHECK
    # ==========================================
    def health_check(self) -> requests.Response:
        """Verifica a saúde da API (Postgres, Mongo, Valkey)."""
        url = f"{self.base_url}/health"
        return self.session.get(url, headers=self._get_headers())

    # ==========================================
    # AUTENTICAÇÃO E USUÁRIOS
    # ==========================================
    def register_user(self, name: str, email: str, password: str, cpf: str) -> requests.Response:
        """Registra um novo usuário."""
        url = f"{self.base_url}/auth/register"
        payload = {
            "name": name,
            "email": email,
            "password": password,
            "cpf": cpf
        }
        return self.session.post(url, json=payload, headers=self._get_headers())

    def create_token(self, email: str, password: str) -> str|None:
        """Gera e armazena o token de autenticação (PASETO)."""
        url = f"{self.base_url}/auth/token"
        payload = {
            "email": email,
            "password": password
        }
        response = self.session.post(url, json=payload, headers=self._get_headers())
        if response.status_code == 200:
            data = response.json()
            self.token = data.get("token")
        return self.token

    # ==========================================
    # PRODUTOS
    # ==========================================
    def create_product(self, name: str, description: str, price: float, stock: int, category_id: str, specs: dict[Any, Any]) -> requests.Response:
        """Cria um novo produto (Requer autenticação)."""
        url = f"{self.base_url}/products"
        payload = {
            "name": name,
            "description": description,
            "price": price,
            "stock": stock,
            "category_id": category_id,
            "specs": specs
        }
        return self.session.post(url, json=payload, headers=self._get_headers())

    def get_product_by_id(self, product_id: str, client_id: str) -> requests.Response:
        """Busca um produto por ID (Requer cabeçalho X-Client-ID)."""
        url = f"{self.base_url}/products/{product_id}"
        headers = self._get_headers(client_id=client_id)
        return self.session.get(url, headers=headers)

    def update_product(self, product_id: str, name: str, description: str, price: float, stock: int, category_id: str, specs: dict[Any, Any]) -> requests.Response:
        """Atualiza um produto por ID (Requer autenticação e ID em formato Hex de 24 chars)."""
        url = f"{self.base_url}/products/{product_id}"
        payload = {
            "name": name,
            "description": description,
            "price": price,
            "stock": stock,
            "category_id": category_id,
            "specs": specs
        }
        return self.session.put(url, json=payload, headers=self._get_headers())

    def delete_product(self, product_id: str) -> requests.Response:
        """Remove um produto por ID (Requer autenticação)."""
        url = f"{self.base_url}/products/{product_id}"
        return self.session.delete(url, headers=self._get_headers())

    def search_products(self, name: str) -> requests.Response:
        """Busca produtos pelo nome (Path param)."""
        url = f"{self.base_url}/products/search/{name}"
        return self.session.get(url, headers=self._get_headers())

    def get_recent_products(self, client_id: str) -> requests.Response:
        """Lista produtos recentes (Requer X-Client-ID)."""
        url = f"{self.base_url}/products/recent"
        headers = self._get_headers(client_id=client_id)
        return self.session.get(url, headers=headers)

    # ==========================================
    # CHECKOUT
    # ==========================================
    def checkout(self, phone: str, items: list[dict[str, Any]]) -> requests.Response:
        """Realiza o checkout gerando um link do WhatsApp."""
        url = f"{self.base_url}/checkout"
        payload = {
            "phone": phone,
            "items": items
        }
        return self.session.post(url, json=payload, headers=self._get_headers())


# Exemplo de Uso:
if __name__ == "__main__":
    client = BazaarGoClient("https://api.bazaargo.com")

    # 1. Health Check
    health = client.health_check()
    print("Health:", health.json())

    # 2. Login para obter token
    # client.create_token("user@email.com", "senhaSegura123")

    # 3. Buscar produto com X-Client-ID obrigatório
    # response = client.get_product_by_id("64a7b2c1e4b0000000000001", client_id="meu-cliente-id")
    # print(response.json())
