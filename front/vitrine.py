
import streamlit as st

from api_call import BazaarGoClient


class Vitrine:
    def __init__(self, client: BazaarGoClient) -> None:
        self.client = client

    def renderizar_principal(self):
        st.title("BaazarGo")
        st.subheader("Nossos produtos")

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
