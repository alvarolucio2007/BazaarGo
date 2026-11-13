import tkinter as tk
from tkinter import messagebox

# -------------------- Configuração da janela --------------------
BG = "#050B16"          # preto azulado
CARD = "#0B1628"        # painel
BLUE = "#087CFF"        # azul principal
BLUE_HOVER = "#0066D6"
WHITE = "#F5F8FF"
MUTED = "#9DB2D1"
FIELD = "#101F35"
BORDER = "#243A59"

root = tk.Tk()
root.title("Admin - BazaarGo")
root.geometry("1000x620")
root.minsize(850, 540)
root.configure(bg=BG)

# -------------------- Layout principal --------------------
main = tk.Frame(root, bg=BG)
main.pack(fill="both", expand=True, padx=28, pady=25)

# Lado esquerdo: marca e apresentação
left = tk.Frame(main, bg=BG)
left.pack(side="left", fill="both", expand=True, padx=(15, 35))

logo = tk.Label(left, text="🛒", font=("Segoe UI Emoji", 54), bg=BG, fg=WHITE)
logo.pack(pady=(90, 5))

brand = tk.Label(left, text="Bazaar", font=("Segoe UI", 38, "bold"),
                 bg=BG, fg=WHITE)
brand.pack()
brand_go = tk.Label(left, text="GO", font=("Segoe UI", 38, "bold"),
                    bg=BG, fg=BLUE)
brand_go.place(relx=0.70, rely=0.35)

subtitle = tk.Label(left, text="PAINEL ADMINISTRATIVO",
                    font=("Segoe UI", 13, "bold"), bg=BG, fg=MUTED)
subtitle.pack(pady=(8, 25))

line = tk.Frame(left, bg=BLUE, height=2, width=300)
line.pack()

description = tk.Label(
    left,
    text="Gerencie sua loja de forma\\nrápida e segura.",
    font=("Segoe UI", 15),
    bg=BG,
    fg=MUTED,
    justify="center"
)
description.pack(pady=28)

# Separador
separator = tk.Frame(main, bg=BORDER, width=1)
separator.pack(side="left", fill="y", padx=(0, 35))

# Lado direito: formulário
right = tk.Frame(main, bg=CARD, highlightbackground=BLUE,
                 highlightthickness=1)
right.pack(side="right", fill="both", expand=True, padx=10, pady=25)

form = tk.Frame(right, bg=CARD)
form.pack(fill="both", expand=True, padx=42, pady=35)

tk.Label(form, text="🔐", font=("Segoe UI Emoji", 35),
         bg=CARD, fg=BLUE).pack(pady=(0, 8))

tk.Label(form, text="Área do Administrador",
         font=("Segoe UI", 22, "bold"), bg=CARD, fg=WHITE).pack()

tk.Label(form, text="Faça login para acessar o painel.",
         font=("Segoe UI", 11), bg=CARD, fg=MUTED).pack(pady=(8, 25))

# Campo de usuário
tk.Label(form, text="Usuário", font=("Segoe UI", 10, "bold"),
         bg=CARD, fg=MUTED, anchor="w").pack(fill="x", pady=(0, 7))

user_entry = tk.Entry(
    form, font=("Segoe UI", 13), bg=FIELD, fg=WHITE,
    insertbackground=WHITE, relief="flat", bd=0
)
user_entry.pack(fill="x", ipady=13, pady=(0, 18))
user_entry.insert(0, "")

# Campo de senha
tk.Label(form, text="Senha", font=("Segoe UI", 10, "bold"),
         bg=CARD, fg=MUTED, anchor="w").pack(fill="x", pady=(0, 7))

password_entry = tk.Entry(
    form, font=("Segoe UI", 13), bg=FIELD, fg=WHITE,
    insertbackground=WHITE, relief="flat", bd=0, show="•"
)
password_entry.pack(fill="x", ipady=13, pady=(0, 10))

show_password = tk.BooleanVar(value=False)

def toggle_password():
    password_entry.configure(show="" if show_password.get() else "•")

tk.Checkbutton(
    form, text="Mostrar senha", variable=show_password,
    command=toggle_password, bg=CARD, fg=MUTED,
    activebackground=CARD, activeforeground=WHITE,
    selectcolor=FIELD, font=("Segoe UI", 9), bd=0
).pack(anchor="w", pady=(0, 20))

# Exemplo demonstrativo de autenticação.
# Em um sistema real, valide no backend e nunca deixe senha fixa no código.
def login():
    username = user_entry.get().strip()
    password = password_entry.get()

    if not username or not password:
        messagebox.showwarning("Campos obrigatórios",
                               "Preencha o usuário e a senha.")
        return

    # Substitua este trecho pela autenticação real da sua API/backend.
    messagebox.showinfo("Login", f"Dados preenchidos para: {username}")

login_button = tk.Button(
    form, text="Entrar  →", command=login,
    font=("Segoe UI", 13, "bold"), bg=BLUE, fg=WHITE,
    activebackground=BLUE_HOVER, activeforeground=WHITE,
    relief="flat", cursor="hand2"
)
login_button.pack(fill="x", ipady=12)

def on_enter(event):
    login_button.configure(bg=BLUE_HOVER)

def on_leave(event):
    login_button.configure(bg=BLUE)

login_button.bind("<Enter>", on_enter)
login_button.bind("<Leave>", on_leave)
password_entry.bind("<Return>", lambda event: login())

tk.Label(
    form, text="BazaarGo © 2026 | Painel Administrativo",
    font=("Segoe UI", 9), bg=CARD, fg=MUTED
).pack(side="bottom", pady=(25, 0))

user_entry.focus()
root.mainloop()
