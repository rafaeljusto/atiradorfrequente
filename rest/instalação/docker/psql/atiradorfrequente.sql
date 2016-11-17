CREATE TYPE LogAcao AS ENUM ('CRIACAO', 'ATUALIZACAO');

CREATE TABLE log (
  id SERIAL PRIMARY KEY,
  data_criacao TIMESTAMP NOT NULL CONSTRAINT data_criacao_mandatorio CHECK (data_criacao > '2016-01-01'::TIMESTAMP),
  endereco_remoto INET
);

CREATE TABLE frequencia_atirador (
  id SERIAL PRIMARY KEY,
  controle VARCHAR NOT NULL CONSTRAINT controle_mandatorio CHECK (controle != ''),
  cr INT NOT NULL CONSTRAINT cr_mandatorio CHECK (cr > 0),
  calibre VARCHAR NOT NULL CONSTRAINT calibre_mandatorio CHECK (calibre != ''),
  arma_utilizada VARCHAR NOT NULL CONSTRAINT arma_utilizada_mandatorio CHECK (arma_utilizada != ''),
  numero_serie VARCHAR NOT NULL DEFAULT '',
  guia_de_trafego INT NOT NULL DEFAULT 0,
  quantidade_municao INT NOT NULL CONSTRAINT quantidade_municao_mandatorio CHECK (quantidade_municao > 0),
  data_inicio TIMESTAMP NOT NULL CONSTRAINT data_inicio_mandatorio CHECK (data_inicio > '2016-01-01'::TIMESTAMP),
  data_termino TIMESTAMP NOT NULL CONSTRAINT data_termino_mandatorio CHECK (data_termino > '2016-01-01'::TIMESTAMP),
  data_criacao TIMESTAMP NOT NULL CONSTRAINT data_criacao_mandatorio CHECK (data_criacao > '2016-01-01'::TIMESTAMP),
  data_atualizacao TIMESTAMP,
  data_confirmacao TIMESTAMP,
  imagem_numero_controle VARCHAR,
  imagem_confirmacao VARCHAR,
  revisao INT NOT NULL DEFAULT 0
);

CREATE TABLE frequencia_atirador_log (
  id SERIAL PRIMARY KEY,
  id_log INT REFERENCES log(id),
  acao LogAcao,
  id_frequencia_atirador INT NOT NULL CONSTRAINT id_frequencia_atirador_mandatorio CHECK (id_frequencia_atirador > 0),
  controle VARCHAR NOT NULL CONSTRAINT controle_mandatorio CHECK (controle != ''),
  cr INT NOT NULL CONSTRAINT cr_mandatorio CHECK (cr > 0),
  calibre VARCHAR NOT NULL CONSTRAINT calibre_mandatorio CHECK (calibre != ''),
  arma_utilizada VARCHAR NOT NULL CONSTRAINT arma_utilizada_mandatorio CHECK (arma_utilizada != ''),
  numero_serie VARCHAR NOT NULL DEFAULT '',
  guia_de_trafego INT NOT NULL DEFAULT 0,
  quantidade_municao INT NOT NULL CONSTRAINT quantidade_municao_mandatorio CHECK (quantidade_municao > 0),
  data_inicio TIMESTAMP NOT NULL CONSTRAINT data_inicio_mandatorio CHECK (data_inicio > '2016-01-01'::TIMESTAMP),
  data_termino TIMESTAMP NOT NULL CONSTRAINT data_termino_mandatorio CHECK (data_termino > '2016-01-01'::TIMESTAMP),
  data_criacao TIMESTAMP NOT NULL CONSTRAINT data_criacao_mandatorio CHECK (data_criacao > '2016-01-01'::TIMESTAMP),
  data_atualizacao TIMESTAMP,
  data_confirmacao TIMESTAMP,
  imagem_numero_controle VARCHAR,
  imagem_confirmacao VARCHAR,
  revisao INT NOT NULL DEFAULT 0
);