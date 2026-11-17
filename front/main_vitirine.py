
import streamlit as st

from api_call import BazaarGoClient
from vitrine import Vitrine


st.set_page_config(
    page_title="BaazarGo",
    page_icon="🛍️",
    layout="wide"
)


@st.cache_resource
def get_api_client():
    return BazaarGoClient(base_url="http://api.localhost")


client = get_api_client()

vitrine = Vitrine(client)
vitrine.renderizar_principal()
