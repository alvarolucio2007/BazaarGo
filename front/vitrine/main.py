import os
import sys
from pathlib import Path

import requests
import streamlit as st

sys.path.append(str(Path(__file__).resolve().parent.parent))

from api_call import BazaarGoClient


class VitrineBazaarGo:
    def __init__(self, client: BazaarGoClient, phone: str) -> None:
        self.client = client
        self.phone = phone

        if "carrinho" not in st.session_state:
            st.session_state["carrinho"] = []

        self.carrinho = st.session_state["carrinho"]

    def renderizar_principal(self):
        st.title("BaazarGo")
        st.subheader("Sua loja online!")

        tab_loja, tab_novidades, tab_detalhes = st.tabs(["🛍️ Loja & Pesquisa", "🔥 Novidades (Recentes)", "🔍 Buscar por ID"])

        with tab_loja:
            self.buscar_produtos()
            st.divider()
            self.listar_produtos()

        with tab_novidades:
            self.listar_recentes(self.phone)

        with tab_detalhes:
            self._render_busca_id_interface()

        st.divider()

        self.exibir_carrinho()
        self.finalizar_compra()

    def listar_produtos(self):
        try:
            produtos = self.client.list_products()

            if not produtos:
                st.info("Nenhum produto encontrado.")
                return

            st.subheader("Catálogo de Produtos")

            for produto in produtos:
                st.write(f"**{produto.name}**")
                st.write(f"Preço: R$ {produto.price:.2f}")
                st.write(produto.description)

                quantidade = st.number_input(
                    "Quantidade",
                    min_value=1,
                    max_value=produto.stock if produto.stock > 0 else 1,
                    value=1,
                    step=1,
                    key=f"qtd_{produto.id}"
                )

                if st.button(
                    "Adicionar ao carrinho",
                    key=f"adicionar_{produto.id}"
                ):
                    self.adicionar_ao_carrinho(produto, quantidade)

                st.divider()

        except requests.exceptions.RequestException as erro:
            st.error(f"Erro ao buscar produtos: {erro}")

    def buscar_produtos(self):
            nome = st.text_input("Pesquisar produto", placeholder="Digite o nome do produto")

            if st.button("Pesquisar"):
                if nome.strip():
                    try:
                        st.session_state["resultados_pesquisa"] = self.client.search_products(nome)
                        if not st.session_state["resultados_pesquisa"]:
                            st.info("Nenhum produto encontrado.")
                    except requests.exceptions.RequestException as erro:
                        st.error(f"Erro ao pesquisar produtos: {erro}")
                else:
                    st.warning("Digite o nome de um produto.")
                    st.session_state["resultados_pesquisa"] = []

            resultados = st.session_state.get("resultados_pesquisa", [])
            if resultados:
                st.subheader("Resultados da pesquisa")

                for produto in resultados:
                    st.write(f"**{produto.name}**")
                    st.write(f"Preço: R$ {produto.price:.2f}")
                    st.write(produto.description)

                    if st.button(
                        "Adicionar ao carrinho",
                        key=f"add_search_{produto.id}"
                    ):
                        self.adicionar_ao_carrinho(produto, 1)

                    st.divider()

    def listar_recentes(self, client_id: str):
        try:
            produtos = self.client.get_recent_products(client_id)

            if not produtos:
                st.info("Nenhum produto recente encontrado.")
                return

            st.subheader("Novidades")

            for produto in produtos:
                st.write(f"**{produto.name}**")
                st.write(f"Preço: R$ {produto.price:.2f}")
                st.write(produto.description)

                if st.button(
                    "Adicionar ao carrinho",
                    key=f"add_recent_{produto.id}"
                ):
                    self.adicionar_ao_carrinho(produto, 1)

                st.divider()

        except requests.exceptions.RequestException as erro:
            st.error(f"Erro ao carregar novidades: {erro}")

    def visualizar_produto(self, product_id: str, client_id: str):
        try:
            produto = self.client.get_product_by_id(product_id, client_id)

            st.subheader(produto.name)
            st.write(f"Preço: R$ {produto.price:.2f}")
            st.write(f"Descrição: {produto.description}")
            st.write(f"Estoque: {produto.stock}")

            if produto.specs:
                st.write("Especificações:")
                st.json(produto.specs)

        except requests.exceptions.RequestException as erro:
            st.error(f"Erro ao consultar produto: {erro}")

    def _render_busca_id_interface(self):
        st.subheader("🔍 Consultar Detalhes por ID")
        st.caption("Atende à rota `/products/{id}` exigida pelo contrato OpenAPI.")
        try:
            todos = self.client.list_products()
            if not todos:
                st.info("Nenhum produto cadastrado.")
                return

            opcoes = {f"{p.name} (ID: {p.id})": p.id for p in todos}
            escolha = st.selectbox("Selecione o produto:", list(opcoes.keys()))

            if st.button("Buscar Dados por ID"):
                id_alvo = opcoes[escolha]
                self.visualizar_produto(id_alvo, self.phone)
        except Exception as e:
            st.error(f"Erro ao buscar: {e}")

    def adicionar_ao_carrinho(self, produto, quantidade: int = 1):
        if quantidade <= 0:
            st.warning("A quantidade deve ser maior que zero.")
            return

        item_existente = next((item for item in self.carrinho if item["name"] == produto.name), None)

        if item_existente:
            item_existente["quantity"] += quantidade
        else:
            item = {
                "name": produto.name,
                "price": produto.price,
                "quantity": quantidade
            }
            self.carrinho.append(item)

        st.success(f"{quantidade}x {produto.name} adicionado(s) ao carrinho!")

    def exibir_carrinho(self):
        st.subheader("🛒 Carrinho de Compras")

        if not self.carrinho:
            st.info("Seu carrinho está vazio.")
            return

        total = 0

        for _, item in enumerate(self.carrinho):
            subtotal = float(item["price"]) * item["quantity"]
            total += subtotal

            st.write(f"**{item['name']}**")
            st.write(f"Quantidade: {item['quantity']}")
            st.write(f"Subtotal: R$ {subtotal:.2f}")
            st.divider()

        st.write(f"### Total: R$ {total:.2f}")

    def finalizar_compra(self):
        if st.button("Finalizar pedido"):
            if not self.carrinho:
                st.warning("Seu carrinho está vazio.")
                return

            try:
                whatsapp_link = self.client.checkout(
                    phone=self.phone,
                    items=self.carrinho
                )

                if whatsapp_link:
                    st.link_button(
                        "Continuar pelo WhatsApp",
                        whatsapp_link
                    )

                self.carrinho.clear()

            except requests.exceptions.RequestException as erro:
                st.error(f"Erro ao finalizar compra: {erro}")

if __name__ == "__main__":
    st.set_page_config(page_title="Vitrine BazaarGo", page_icon="🛒", layout="wide")
    API_URL = os.getenv("API_URL", "http://api.localhost")
    @st.cache_resource
    def get_api_client():
        return BazaarGoClient(base_url=API_URL)

    client = get_api_client()
    vitrine = VitrineBazaarGo(client, "99999999999")  # Change this phone number in "prod" (as if one day this will ever see prod)
    vitrine.renderizar_principal()
