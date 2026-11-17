
import streamlit as st
import requests

from api_call import BazaarGoClient


class Vitrine:
    def __init__(self, client: BazaarGoClient) -> None:
        self.client = client

    def renderizar_principal(self):
        st.title("BaazarGo")
        st.subheader("Sua loja online!")

        self.buscar_produtos()

        st.divider()

        self.listar_produtos()

        st.divider()

        self.exibir_carrinho()

        self.finalizar_compra()

        try:
            res = self.client.get_all_products()
            res.raise_for_status()

            dados = res.json()
            produtos = dados.get("data", []) if isinstance(dados, dict) else dados

            if not produtos:
                st.info("Nenhum produto encontrado.")
                return

            for produto in produtos:
                st.write(f"**{produto['name']}**")
                st.write(f"Preço: R$ {produto['price']}")
                st.write(produto["description"])
                st.divider()

        except Exception as erro:
            st.error(f"Não foi possível carregar os produtos: {erro}")

def __init__(self, client: BazaarGoClient) -> None:
    self.client = client

    if "carrinho" not in st.session_state:
        st.session_state["carrinho"] = []

    self.carrinho = st.session_state["carrinho"]

    # Listagem de produtos
def listar_produtos(self):
    try:
        resposta = self.client.get_all_products()
        resposta.raise_for_status()

        dados = resposta.json()
        produtos = dados.get("data", []) if isinstance(dados, dict) else dados

        if not produtos:
            st.info("Nenhum produto encontrado.")
            return

        for produto in produtos:
            st.write(f"**{produto['name']}**")
            st.write(f"Preço: R$ {produto['price']}")
            st.write(produto['description'])

            if st.button(
                "Adicionar ao carrinho",
                key=f"adicionar_{produto['id']}"
            ):
                self.adicionar_ao_carrinho(produto)

            st.divider()

    except requests.exceptions.RequestException as erro:
        st.error(f"Erro ao buscar produtos: {erro}")

    # Pesquisa de produtos
def buscar_produtos(self):
    nome = st.text_input("Pesquisar produto", placeholder="Digite o nome do produto")

    if st.button("Pesquisar"):
        if nome.strip():
            try:
                resposta = self.client.search_products(nome)
                resposta.raise_for_status()

                produtos = resposta.json()

                if not produtos:
                    st.info("Nenhum produto encontrado.")
                    return

                st.subheader("Resultados da pesquisa")

                for produto in produtos:
                    st.write(f"**{produto['name']}**")
                    st.write(f"Preço: R$ {produto['price']}")
                    st.write(produto['description'])
                    st.divider()

            except requests.exceptions.RequestException as erro:
                st.error(f"Erro ao pesquisar produtos: {erro}")
        else:
            st.warning("Digite o nome de um produto.")

    # Produtos recentes
def listar_recentes(self,client_id: str):
    try:
        resposta = self.client.get_recent_products(client_id)
        resposta.raise_for_status()

        dados = resposta.json()
        produtos = dados.get("data", []) if isinstance(dados, dict) else dados

        if not produtos:
            st.info("Nenhum produto recente encontrado.")
            return

        st.subheader("Novidades")

        for produto in produtos:
            st.write(f"**{produto['name']}**")
            st.write(f"Preço: R$ {produto['price']}")
            st.write(produto['description'])
            st.divider()

    except requests.exceptions.RequestException as erro:
        st.error(f"Erro ao carregar novidades: {erro}")

    # Detalhes do produto
def visualizar_produto(self, product_id, client_id: str):
    try:
        resposta = self.client.get_product_by_id(product_id, client_id)
        resposta.raise_for_status()

        produto = resposta.json()

        if isinstance(produto, dict) and "data" in produto:
            produto = produto["data"]

        st.subheader(produto["name"])
        st.write(f"Preço: R$ {produto['price']}")
        st.write(f"Descrição: {produto['description']}")
        st.write(f"Estoque: {produto['stock']}")

        if produto.get("specs"):
            st.write("Especificações:")
            st.json(produto["specs"])

    except requests.exceptions.RequestException as erro:
        st.error(f"Erro ao consultar produto: {erro}")

    # Carrinho
def adicionar_ao_carrinho(self, produto: dict, quantidade: int = 1):
    if quantidade <= 0:
        st.warning("A quantidade deve ser maior que zero.")
        return

    item = {
        "name": produto["name"],
        "price": produto["price"],
        "quantity": quantidade
    }

    self.carrinho.append(item)

    st.success(f"{produto['name']} adicionado ao carrinho!")
    
    # Finalização da compra
def exibir_carrinho(self):
    st.subheader("🛒")

    if not self.carrinho:
        st.info("Seu carrinho está vazio.")
        return

    total = 0

    for item in self.carrinho:
        subtotal = float(item["price"]) * item["quantity"]
        total += subtotal

        st.write(f"**{item['name']}**")
        st.write(f"Quantidade: {item['quantity']}")
        st.write(f"Subtotal: R$ {subtotal:.2f}")
        st.divider()

    st.write(f"### Total: R$ {total:.2f}")

def finalizar_compra(self):
    st.subheader("Finalizar compra")

    telefone = st.text_input(
        "Seu telefone com DDD",
        placeholder="Ex: 86999999999"
    )

    if st.button("Finalizar pedido"):
        if not self.carrinho:
            st.warning("Seu carrinho está vazio.")
            return

        if not telefone.strip():
            st.warning("Informe seu telefone.")
            return

        try:
            resposta = self.client.checkout(
                phone=telefone,
                items=self.carrinho
            )

            resposta.raise_for_status()

            dados = resposta.json()

            st.success("Pedido enviado com sucesso!")

            if isinstance(dados, dict) and dados.get("whatsapp_link"):
                st.link_button(
                    "Continuar pelo WhatsApp",
                    dados["whatsapp_link"]
                )

            self.carrinho.clear()

        except requests.exceptions.RequestException as erro:
            st.error(f"Erro ao finalizar compra: {erro}")