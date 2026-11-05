CREATE TABLE users(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  cpf VARCHAR(11) UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash BYTEA NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,

  CONSTRAINT users_email_format CHECK (email ~* '^[^@\s]+@[^@\s]+\.[^@\s]+$'),
  CONSTRAINT users_name_not_empty CHECK (LENGTH(TRIM(name)) > 0),
  CONSTRAINT users_cpf_format CHECK (cpf IS NULL OR cpf ~ '^\d{11}$')
)
