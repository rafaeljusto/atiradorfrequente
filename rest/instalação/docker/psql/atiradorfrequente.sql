CREATE TYPE LogAcao AS ENUM ('CRIACAO', 'ATUALIZACAO');

CREATE TABLE log (
  id SERIAL PRIMARY KEY,
  data_criacao TIMESTAMP,
  endereco_remoto INET
);

CREATE TABLE frequencia_atirador (
  id SERIAL PRIMARY KEY,
  controle VARCHAR,
  cr INT,
  calibre VARCHAR,
  arma_utilizada VARCHAR,
  numero_serie VARCHAR,
  guia_de_trafego INT,
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

CREATE TABLE frequencia_atirador_log (
  id SERIAL PRIMARY KEY,
  id_log INT REFERENCES log(id),
  acao LogAcao,
  id_frequencia_atirador INT,
  controle VARCHAR,
  cr INT,
  calibre VARCHAR,
  arma_utilizada VARCHAR,
  numero_serie VARCHAR,
  guia_de_trafego INT,
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