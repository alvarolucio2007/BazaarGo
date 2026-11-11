
import streamlit as st

# Configuração da página
st.set_page_config(
    page_title="BaazarGo",
    layout="wide"
)

# Cores da loja
AZUL = "#263B63"
ROSE = "#E8BFA8"
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

.stButton > button {{
    background-color: {AZUL};
    color: white;
    border-radius: 8px;
    border: none;
}}

.stButton > button:hover {{
    background-color: {ROSE};
    color: {AZUL};
}}
</style>
""", unsafe_allow_html=True)

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

# Categorias
st.subheader("Explore nossas categorias")

col1, col2, col3 = st.columns(3)

with col1:
    st.button(" Moda", use_container_width=True)

with col2:
    st.button(" Tecnologia", use_container_width=True)

with col3:
    st.button(" Casa", use_container_width=True)

st.divider()

# Área de produtos
st.subheader("Produtos em destaque")

st.info("Em breve, seus produtos aparecerão aqui!")

st.write("")
st.caption("BaazerGo © 2026 — Sua loja, suas escolhas.")
