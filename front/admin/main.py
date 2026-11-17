import os

import pandas as pd
import streamlit as st
from api_call import BazaarGoClient


class Admin:
    def __init__(self, client: BazaarGoClient) -> None:
        self.client = client

    def render(self) -> None:
        if "token" not in st.session_state or not st.session_state["token"]:
            self.render_login()
            return
        self.client.token = st.session_state["token"]
        self.render_sidebar()

    def render_login(self) -> None:
        st.title("BaazarGo | Acesso")
        with st.container():
            _col1, col2, _col3 = st.columns([1, 2, 1])
            with col2:
                tab_login, tab_registro = st.tabs(["Login", "Criar Conta"])

                with tab_login, st.form("form_login"):
                    email = st.text_input("E-mail")
                    pw = st.text_input("Senha", type="password")
                    entrar = st.form_submit_button("Entrar", use_container_width=True)

                    if entrar:
                        if not email or not pw:
                            st.warning("Preencha e-mail e senha.")
                        else:
                            try:
                                token = self.client.create_token(email, pw)
                                if token:
                                    st.session_state["token"] = token
                                    st.rerun()
                                else:
                                    st.error("E-mail ou senha inválidos.")
                            except Exception as e:
                                st.error(f"Erro de conexão: {e}")

                with tab_registro:
                    with st.form("form_registro_inicial"):
                        r_nome = st.text_input("Nome")
                        r_email = st.text_input("E-mail")
                        r_cpf = st.text_input("CPF (só números)")
                        r_senha = st.text_input("Senha", type="password")
                        registrar = st.form_submit_button("Registrar", use_container_width=True)

                        if registrar:
                            if not r_nome or not r_email or not r_cpf or not r_senha:
                                st.warning("Preencha todos os campos.")
                            else:
                                try:
                                    res = self.client.register_user(r_nome, r_email, r_senha, r_cpf)
                                    if res.status_code in (200, 201):
                                        st.success("Conta criada! Faça login para continuar.")
                                    else:
                                        st.error(res.json().get("error", "Erro ao criar conta."))
                                except Exception as e:
                                    st.error(f"Erro de conexão: {e}")

    def render_sidebar(self) -> None:
        st.sidebar.title("BaazarGo")
        st.sidebar.caption("Painel de Controle")

        pagina = st.sidebar.radio("Menu", ["Status do Sistema", "Produtos", "Usuários"])

        st.sidebar.markdown("---")
        if st.sidebar.button("Sair (Logout)", use_container_width=True):
            st.session_state["token"] = None
            st.rerun()

        if pagina == "Status do Sistema":
            self.render_status()
        elif pagina == "Produtos":
            self.render_produtos()
        elif pagina == "Usuários":
            self.render_usuarios()

    def render_status(self) -> None:
        st.title("Status do Sistema")
        st.caption("Visão geral dos status do sistema")


        try:
            health = self.client.health_check()
            if health.status_code == 200:
                st.success(" Sistemas operacionais e bancos de dados conectados.")
            else:
                st.error(" Instabilidade nos servidores.")
        except:
            st.error(" Falha ao conectar com a API.")


    def render_produtos(self) -> None:
        st.title("Gerenciamento de Produtos")

        produtos = []
        try:
            res = self.client.get_all_products()
            if res.status_code == 200:
                produtos = res.json().get("data", []) if isinstance(res.json(), dict) else res.json()
        except Exception as e:
            st.error(f"Erro ao buscar produtos: {e}")

        with st.expander("➕ Cadastrar Novo Produto", expanded=False):
            with st.form("form_novo_produto"):
                nome = st.text_input("Nome do produto")
                descricao = st.text_area("Descrição")
                preco = st.number_input("Preço (R$)", min_value=0.0, step=1.0)
                estoque = st.number_input("Estoque", min_value=0, step=1)
                categoria = st.text_input("ID da categoria")
                enviar = st.form_submit_button("Cadastrar", use_container_width=True)

                if enviar:
                    if nome.strip():
                        res = self.client.create_product(nome, descricao, preco, estoque, categoria, {})
                        if res.status_code == 201:
                            st.success("Produto cadastrado!")
                            st.rerun()
                        else:
                            st.error("Erro ao cadastrar.")
                    else:
                        st.warning("O nome é obrigatório.")

        st.divider()
        st.subheader("Produtos Cadastrados")

        if not produtos:
            st.info("Nenhum produto cadastrado.")
            return

        df = pd.DataFrame(produtos)
        if not df.empty:
            df_view = df[["id", "name", "price", "stock"]].copy()
            df_view.columns = ["ID", "Nome", "Preço", "Estoque"]
            st.dataframe(df_view, use_container_width=True)

        with st.expander(" Editar ou Excluir Produto", expanded=False):
            opcoes = {f"{p['name']} (ID: {p['id']})": p for p in produtos}
            produto_selecionado_nome = st.selectbox("Selecione o produto:", list(opcoes.keys()))

            if produto_selecionado_nome:
                prod = opcoes[produto_selecionado_nome]

                with st.form("form_editar_produto"):
                    st.write(f"Editando: **{prod['name']}**")
                    edit_nome = st.text_input("Nome", value=prod.get("name", ""))
                    edit_desc = st.text_area("Descrição", value=prod.get("description", ""))
                    edit_preco = st.number_input("Preço", value=float(prod.get("price", 0.0)))
                    edit_estoque = st.number_input("Estoque", value=int(prod.get("stock", 0)))
                    edit_cat = st.text_input("Categoria ID", value=prod.get("category_id", ""))

                    col_salvar, col_deletar = st.columns(2)
                    btn_salvar = col_salvar.form_submit_button("💾 Salvar Alterações", use_container_width=True)
                    btn_deletar = col_deletar.form_submit_button("🗑️ Excluir Produto", use_container_width=True)

                    if btn_salvar:
                        res = self.client.update_product(
                            product_id=prod["id"],
                            name=edit_nome, description=edit_desc, price=edit_preco,  # pyright: ignore[reportArgumentType]
                            stock=edit_estoque, category_id=edit_cat, specs=prod.get("specs", {})  # pyright: ignore[reportArgumentType]
                        )
                        if res.status_code in (200, 204):
                            st.success("Produto atualizado!")
                            st.rerun()
                        else:
                            st.error(f"Erro ao atualizar: {res.text}")

                    if btn_deletar:
                        res = self.client.delete_product(prod["id"])
                        if res.status_code in (200, 204):
                            st.success("Produto excluído!")
                            st.rerun()
                        else:
                            st.error(f"Erro ao excluir: {res.text}")

    def render_usuarios(self) -> None:
        st.title("Gerenciamento de Usuários")
        st.caption("Cadastre novos administradores ou clientes.")

        with st.form("form_novo_usuario"):
            nome = st.text_input("Nome Completo")
            email = st.text_input("E-mail")
            cpf = st.text_input("CPF (apenas números)")
            senha = st.text_input("Senha", type="password")

            registrar = st.form_submit_button("Criar Usuário", use_container_width=True)

            if registrar:
                if not nome or not email or not cpf or not senha:
                    st.warning("Preencha todos os campos.")
                else:
                    try:
                        res = self.client.register_user(nome, email, senha, cpf)
                        if res.status_code in (200, 201):
                            st.success(f"Usuário {nome} criado com sucesso!")
                        else:
                            st.error(f"Erro: {res.json().get('error', 'Falha ao criar usuário')}")
                    except Exception as e:
                        st.error(f"Erro de conexão: {e}")

if __name__ == "__main__":
    st.set_page_config(page_title="BaazarGo Admin", page_icon="🛒", layout="wide")
    API_URL = os.getenv("API_URL", "http://localhost:8080")
    @st.cache_resource
    def get_api_client():
        return BazaarGoClient(base_url=API_URL)

    client = get_api_client()
    admin = Admin(client)
    admin.render()
