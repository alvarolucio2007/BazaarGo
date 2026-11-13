import streamlit as st

# Configuração da página
st.set_page_config(
    page_title="BaazarGo",
    layout="wide"
)

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
