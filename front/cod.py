
import streamlit as st
import pandas as pd

st.set_page_config(
    page_title="BaazarGo Admin",
    layout="wide"
)

AZUL = "#CADCFD"
ROSA = "#C08664"
BEGE = "#F7DBA4"

st.title(" BaazarGo Admin")
st.caption("Painel administrativo da loja")

pagina = st.sidebar.radio(
    "Navegação",
    [
        "Dashboard",
        "Produtos",
        "Pedidos",
        "Clientes",
        "Categorias",
        "Relatórios",
        "Configurações"
    ]
)

st.sidebar.markdown("---")
st.sidebar.write("Administrador")

# DASHBOARD
if pagina == "Dashboard":

    st.header("Dashboard")

    col1, col2, col3, col4 = st.columns(4)

    col1.metric("Faturamento", "R$ 0,00")
    col2.metric("Pedidos", "0")
    col3.metric("Produtos", "0")
    col4.metric("Estoque baixo", "0")

    st.subheader("Últimos pedidos")

    st.dataframe(
        pd.DataFrame(
            columns=[
                "Pedido",
                "Cliente",
                "Data",
                "Valor",
                "Status"
            ]
        ),
        use_container_width=True
    )

# PRODUTOS
elif pagina == "Produtos":

    st.header("Produtos")

    with st.form("produto"):

        nome = st.text_input("Nome")
        descricao = st.text_area("Descrição")
        preco = st.number_input("Preço", min_value=0.0)
        estoque = st.number_input("Estoque", min_value=0)
        categoria = st.text_input("Categoria")

        cadastrar = st.form_submit_button("Cadastrar")

    if cadastrar:
        if nome.strip():
            st.success("Produto preenchido!")
        else:
            st.warning("Informe o nome do produto.")

    st.subheader("Lista de produtos")

    st.dataframe(
        pd.DataFrame(
            columns=[
                "ID",
                "Nome",
                "Preço",
                "Estoque",
                "Categoria"
            ]
        ),
        use_container_width=True
    )

# PEDIDOS
elif pagina == "Pedidos":

    st.header("Pedidos")

    st.dataframe(
        pd.DataFrame(
            columns=[
                "Pedido",
                "Cliente",
                "Data",
                "Valor",
                "Status"
            ]
        ),
        use_container_width=True
    )

# CLIENTES
elif pagina == "Clientes":

    st.header("Clientes")

    st.dataframe(
        pd.DataFrame(
            columns=[
                "ID",
                "Nome",
                "E-mail",
                "Pedidos"
            ]
        ),
        use_container_width=True
    )

# CATEGORIAS
elif pagina == "Categorias":

    st.header("Categorias")

    categoria = st.text_input("Nome da categoria")

    if st.button("Criar categoria"):
        if categoria.strip():
            st.success("Categoria criada!")
        else:
            st.warning("Informe o nome.")

    st.dataframe(
        pd.DataFrame(
            columns=["ID", "Nome"]
        ),
        use_container_width=True
    )

# RELATÓRIOS
elif pagina == "Relatórios":

    st.header("Relatórios")

    st.info("Relatórios de vendas e produtos.")

# CONFIGURAÇÕES
elif pagina == "Configurações":

    st.header("Configurações")

    nome = st.text_input(
        "Nome da loja",
        value="BaazerGo"
    )

    email = st.text_input("E-mail administrativo")

    if st.button("Salvar"):
        st.success("Configurações preenchidas!")

