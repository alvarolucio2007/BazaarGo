import json
import os
import sys
from pathlib import Path

import pandas as pd
import streamlit as st

sys.path.append(str(Path(__file__).resolve().parent.parent))

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
                                st.error(f"Erro de login: {e}")

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
                                    self.client.register_user(r_nome, r_email, r_senha, r_cpf)
                                    st.success("Conta criada! Faça login para continuar.")
                                except Exception as e:
                                    st.error(f"Erro ao criar conta: {e}")

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
            if all(health.values()):
                st.success("Sistemas operacionais e bancos de dados conectados.")
            else:
                st.warning("Algum serviço pode estar instável.")
                st.json(health)
        except Exception as e:
            st.error(f"Falha ao conectar com a API: {e}")

    def render_produtos(self) -> None:
        st.title("Gerenciamento de Produtos")

        produtos = []
        try:
            produtos = self.client.list_products()
        except Exception as e:
            st.error(f"Erro ao buscar produtos: {e}")

        # --- CADASTRO DE NOVO PRODUTO ---
        with st.expander("➕ Cadastrar Novo Produto", expanded=False):
            with st.form("form_novo_produto"):
                nome = st.text_input("Nome do produto")
                descricao = st.text_area("Descrição")
                preco = st.number_input("Preço (R$)", min_value=0.0, step=1.0)
                estoque = st.number_input("Estoque", min_value=0, step=1)
                categoria = st.text_input("ID da categoria")

                st.subheader("Especificações (Specs)")
                st.caption("Digite em formato JSON válido (ex: {\"Cor\": \"Azul\", \"Voltagem\": \"220V\"})")
                specs_input = st.text_area("Specs (JSON)", value="{}")

                enviar = st.form_submit_button("Cadastrar", use_container_width=True)

                if enviar:
                    if nome.strip():
                        try:
                            # Converte a string do JSON para dicionário Python
                            specs_dict = json.loads(specs_input) if specs_input.strip() else {}

                            self.client.create_product(
                                name=nome,
                                description=descricao,
                                price=preco,
                                stock=int(estoque),
                                category_id=categoria,
                                specs=specs_dict
                            )
                            st.success("Produto cadastrado com sucesso!")
                            st.rerun()
                        except json.JSONDecodeError:
                            st.error("O campo de Specs deve ser um JSON válido (ex: {\"chave\": \"valor\"}).")
                        except Exception as e:
                            st.error(f"Erro ao cadastrar produto: {e}")
                    else:
                        st.warning("O nome é obrigatório.")

        st.divider()
        st.subheader("Produtos Cadastrados")

        if not produtos:
            st.info("Nenhum produto cadastrado.")
            return

        dados_tabela = [
            {
                "ID": p.id,
                "Nome": p.name,
                "Preço": p.price,
                "Estoque": p.stock,
                "Categoria": p.category_id,
            }
            for p in produtos
        ]
        df = pd.DataFrame(dados_tabela)
        if not df.empty:
            st.dataframe(df, use_container_width=True)

        # --- EDIÇÃO E EXCLUSÃO DE PRODUTO ---
        with st.expander("✏️ Editar ou Excluir Produto", expanded=False):
            opcoes = {f"{p.name} (ID: {p.id})": p for p in produtos}
            produto_selecionado_nome = st.selectbox("Selecione o produto:", list(opcoes.keys()))

            if produto_selecionado_nome:
                prod = opcoes[produto_selecionado_nome]

                with st.form("form_editar_produto"):
                    st.write(f"Editando: **{prod.name}**")
                    edit_nome = st.text_input("Nome", value=prod.name)
                    edit_desc = st.text_area("Descrição", value=prod.description)
                    edit_preco = st.number_input("Preço", value=float(prod.price))
                    edit_estoque = st.number_input("Estoque", value=int(prod.stock))
                    edit_cat = st.text_input("Categoria ID", value=prod.category_id)

                    # Converte o dicionário atual de specs para string JSON para preencher o input
                    specs_atual_str = json.dumps(prod.specs, ensure_ascii=False) if prod.specs else "{}"
                    edit_specs_input = st.text_area("Specs (JSON)", value=specs_atual_str)

                    col_salvar, col_deletar = st.columns(2)
                    btn_salvar = col_salvar.form_submit_button("💾 Salvar Alterações", use_container_width=True)
                    btn_deletar = col_deletar.form_submit_button("🗑️ Excluir Produto", use_container_width=True)

                    if btn_salvar:
                        try:
                            specs_dict = json.loads(edit_specs_input) if edit_specs_input.strip() else {}

                            self.client.update_product(
                                product_id=prod.id,
                                name=edit_nome,
                                description=edit_desc,
                                price=edit_preco,
                                stock=int(edit_estoque),
                                category_id=edit_cat,
                                specs=specs_dict
                            )
                            st.success("Produto atualizado com sucesso!")
                            st.rerun()
                        except json.JSONDecodeError:
                            st.error("O campo de Specs deve ser um JSON válido.")
                        except Exception as e:
                            st.error(f"Erro ao atualizar: {e}")

                    if btn_deletar:
                        try:
                            self.client.delete_product(prod.id)
                            st.success("Produto excluído com sucesso!")
                            st.rerun()
                        except Exception as e:
                            st.error(f"Erro ao excluir: {e}")

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
                        self.client.register_user(nome, email, senha, cpf)
                        st.success(f"Usuário {nome} criado com sucesso!")
                    except Exception as e:
                        st.error(f"Erro ao criar usuário: {e}")

if __name__ == "__main__":
    st.set_page_config(page_title="BaazarGo Admin", page_icon="🛒", layout="wide")
    API_URL = os.getenv("API_URL", "http://api.localhost")
    @st.cache_resource
    def get_api_client():
        return BazaarGoClient(base_url=API_URL)

    client = get_api_client()
    admin = Admin(client)
    admin.render()
