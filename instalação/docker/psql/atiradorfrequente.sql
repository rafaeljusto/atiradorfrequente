CREATE ROLE atiradorfrequente;
ALTER ROLE atiradorfrequente WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN NOREPLICATION NOBYPASSRLS PASSWORD 'md578e0745d7f14ffd47c0a6bff808da2a4';

CREATE TYPE LogAcao AS ENUM ('CRIACAO', 'ATUALIZACAO');

CREATE TABLE log (
  id SERIAL PRIMARY KEY,
  data_criacao TIMESTAMP,
  endereco_remoto INET
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE log TO atiradorfrequente;

CREATE TABLE frequencia_atirador (
  id SERIAL PRIMARY KEY,
  controle VARCHAR,
  cr VARCHAR,
  calibre VARCHAR,
  arma_utilizada VARCHAR,
  numero_serie VARCHAR,
  guia_de_trafego VARCHAR,
  quantidade_municao INT,
  data_inicio TIMESTAMP,
  data_termino TIMESTAMP,
  data_criacao TIMESTAMP,
  data_atualizacao TIMESTAMP,
  data_confirmacao TIMESTAMP,
  imagem_numero_controle VARCHAR,
  imagem_confirmacao VARCHAR,
  revisao INT
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE frequencia_atirador TO atiradorfrequente;

CREATE TABLE frequencia_atirador_log (
  id SERIAL PRIMARY KEY,
  id_log INT REFERENCES log(id),
  acao LogAcao,
  id_frequencia_atirador INT,
  controle VARCHAR,
  cr VARCHAR,
  calibre VARCHAR,
  arma_utilizada VARCHAR,
  numero_serie VARCHAR,
  guia_de_trafego VARCHAR,
  quantidade_municao INT,
  data_inicio TIMESTAMP,
  data_termino TIMESTAMP,
  data_criacao TIMESTAMP,
  data_atualizacao TIMESTAMP,
  data_confirmacao TIMESTAMP,
  imagem_numero_controle VARCHAR,
  imagem_confirmacao VARCHAR,
  revisao INT
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE frequencia_atirador_log TO atiradorfrequente;