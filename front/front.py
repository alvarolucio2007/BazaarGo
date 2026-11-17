import streamlit as st
import requests
# Configuração da página

API_URL="http://api.localhost"
st.set_page_config(
    page_title="BaazarGo",
    layout="wide"
)
def buscar_produtos():
  try:
        resposta = requests.get(
            f"{API_URL}/products/recent",
            timeout=10
        )

        resposta.raise_for_status()

        dados = resposta.json()

        return dados

  except requests.exceptions.RequestException as erro:
        print(f"Erro ao buscar produtos: {erro}")
        return []


produtos = buscar_produtos()

st.title("BazaarGo - Nossa Vitrine")

if produtos:

    colunas = st.columns(3)

    for i, produto in enumerate(produtos):

        with colunas[i % 3]:

            st.subheader(produto["name"])

            st.write(
                f"R$ {float(produto['price']):.2f}"
            )

            st.write(produto.get("description", ""))

else:
    st.info("Nenhum produto encontrado.")


# Cores da loja
AZUL = "#2995bc"
ROSE = "#8BFA8"
FUNDO = "#F8F5EF"

# Estilo visual
st.markdown(f"""
<style>
.stApp {{
    background-color: {FUNDO};
}}

h1, h2, h3 {{
    color: {AZUL};
}}

p, label {{
    color: {AZUL};
}}

.topo {{
    background-color: {AZUL};
    padding: 25px;
    border-radius: 12px;
    color: white;
    text-align: center;
}}

.topo h1 {{
    color: white;
    font-size: 42px;
}}

.topo p {{
    color: white;
    font-size: 17px;
}}

""", unsafe_allow_html=True)

pagina = st.sidebar.radio(
    "Navegação",
    [
        "Produtos",
        "Carrinho"
    ]
)
if pagina=="Carrinho":
    st.write("teste")
# Cabeçalho
st.markdown("""
<div class="topo">
    <h1>BaazerGo </h1>
    <p>Encontre seus produtos favoritos em um só lugar.</p>
</div>
""", unsafe_allow_html=True)

st.write("")

# Pesquisa
pesquisa = st.text_input(
    "O que você está procurando?",
    placeholder="Digite o nome de um produto..."
)

# Área de produtos
st.markdown('<h3 style="color: black;">Produtos em destaque</h3>', unsafe_allow_html=True)

st.info("Em breve, seus produtos aparecerão aqui!")

st.write("")
st.caption("BaazerGo © 2026 — Sua loja, suas escolhas.")
