\set frequencia_atirador_campos 'id, controle, cr, calibre, arma_utilizada, numero_serie, guia_de_trafego, quantidade_municao, data_inicio, data_termino, data_criacao, data_atualizacao, data_confirmacao, imagem_numero_controle, imagem_confirmacao, revisao'
\set frequencia_atirador_log_campos 'id, id_log, acao, id_frequencia_atirador, controle, cr, calibre, arma_utilizada, numero_serie, guia_de_trafego, quantidade_municao, data_inicio, data_termino, data_criacao, data_atualizacao, data_confirmacao, imagem_numero_controle, imagem_confirmacao, revisao'
\set log_campos 'id, data_criacao, endereco_remoto'

WITH

--
-- Frequência sem confirmação
--

idFrq1 AS (
  INSERT INTO frequencia_atirador (:frequencia_atirador_campos)
  VALUES (DEFAULT, 1234, 380308, '.380', 'Arma do Clube', 'HG72643653', 762556223, 100,
  NOW() - interval '2 hour',
  NOW() - interval '30 minutes',
  NOW() - interval '29 minutes', NULL, NULL,
  'iVBORw0KGgoAAAANSUhEUgAAAA0AAAANAQAAAABakNnRAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAgY0hSTQAAeiYAAICEAAD6AAAAgOgAAHUwAADqYAAAOpgAABdwnLpRPAAAAAJiS0dEAAHdihOkAAAACXBIWXMAAABIAAAASABGyWs+AAAAG0lEQVQI12P4/4Oht4Ph7g4wgjFWrWAIDUAmAWzDEKvudB5PAAAAAElFTkSuQmCC',
  NULL, 0) RETURNING id
),

idLog1 AS (
  INSERT INTO log (:log_campos)
  VALUES (DEFAULT, NOW() - interval '29 minutes', '198.51.100.1')
  RETURNING id
),

--
-- Frequência sem confirmação expirada
--

idFrq2 AS (
  INSERT INTO frequencia_atirador (:frequencia_atirador_campos)
  VALUES (DEFAULT, 7344, 923714, '.45', 'Imbel 1911', 'SF9153921', 839201286, 150,
  NOW() - interval '2 hour',
  NOW() - interval '40 minutes',
  NOW() - interval '31 minutes', NULL, NULL,
  'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAi0lEQVR4Xo2PIQ4CMRBF24QDIZEILHJXcQFSsxssGg2KC2A5QGVlT9FrzNZ8mIT87HQhYZKX5v+Zdvo9APdPrVrjMNx583E9ejYAkD7cUEohqj89+2KtVTH662qR6Y0Y3a4m2+6MnDP0nPl2cLM/IaVEVC8G17uAGOMC9U0Ymez/Zj7DMOF4ebpf9QJNxJnbJOlfVQAAAABJRU5ErkJggg==',
  NULL, 0) RETURNING id
),

idLog2 AS (
  INSERT INTO log (:log_campos)
  VALUES (DEFAULT, NOW() - interval '31 minutes', '198.51.100.2')
  RETURNING id
),

--
-- Frequência com confirmação
--

idFrq3 AS (
  INSERT INTO frequencia_atirador (:frequencia_atirador_campos)
  VALUES (DEFAULT, 1246, 114239, '.40', 'Imbel MD2', 'DL28461184', 102483466, 50,
  NOW() - interval '5 hours',
  NOW() - interval '4 hours' - interval '30 minutes',
  NOW() - interval '4 hours' - interval '10 minutes',
  NOW() - interval '4 hours' - interval '5 minutes',
  NOW() - interval '4 hours' - interval '5 minutes',
  'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAUklEQVR4XqWQ0QmAMAxEL8FhdIp2Hvep87hFxzm5D6GE0ip9kI/A5XHESCJfBzHgPqtZKjvxAe9cawbBCVtrEmFX///G9zKaFjtGc8T1TExQ5gHN9xsWe3/FugAAAABJRU5ErkJggg==',
  'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAATElEQVR4XpWPCwoAIAhDswPrWTyxMcKEBn2EwYPG0yQi2st0gJllO5kHRlUNBIwsrkzjbnN3AdPqc5lXV/gMNqYt7cn9Vvr+NT0k7xkVBY1RndW3lwAAAABJRU5ErkJggg==',
  1) RETURNING id
),

idLog3 AS (
  INSERT INTO log (:log_campos)
  VALUES (DEFAULT, NOW() - interval '4 hours' - interval '10 minutes', '198.51.100.3')
  RETURNING id
),

idLog4 AS (
  INSERT INTO log (:log_campos)
  VALUES (DEFAULT, NOW() - interval '4 hours' - interval '5 minutes', '198.51.100.3')
  RETURNING id
)

--
-- Insere informações na tabela log correspondente
--

INSERT INTO frequencia_atirador_log (:frequencia_atirador_log_campos) VALUES

--
-- Frequência sem confirmação (log)
--

(DEFAULT, (SELECT id FROM idLog1), 'CRIACAO', (SELECT id FROM idFrq1), 1234, 380308, '.380',
'Arma do Clube', 'HG72643653', 762556223, 100,
NOW() - interval '2 hour',
NOW() - interval '30 minutes',
NOW() - interval '29 minutes', NULL, NULL,
'iVBORw0KGgoAAAANSUhEUgAAAA0AAAANAQAAAABakNnRAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAgY0hSTQAAeiYAAICEAAD6AAAAgOgAAHUwAADqYAAAOpgAABdwnLpRPAAAAAJiS0dEAAHdihOkAAAACXBIWXMAAABIAAAASABGyWs+AAAAG0lEQVQI12P4/4Oht4Ph7g4wgjFWrWAIDUAmAWzDEKvudB5PAAAAAElFTkSuQmCC',
NULL, 0),

--
-- Frequência sem confirmação expirada (log)
--

(DEFAULT, (SELECT id FROM idLog2), 'CRIACAO', (SELECT id FROM idFrq2), 7344, 923714, '.45',
'Imbel 1911', 'SF9153921', 839201286, 150,
NOW() - interval '2 hour',
NOW() - interval '40 minutes',
NOW() - interval '31 minutes', NULL, NULL,
'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAi0lEQVR4Xo2PIQ4CMRBF24QDIZEILHJXcQFSsxssGg2KC2A5QGVlT9FrzNZ8mIT87HQhYZKX5v+Zdvo9APdPrVrjMNx583E9ejYAkD7cUEohqj89+2KtVTH662qR6Y0Y3a4m2+6MnDP0nPl2cLM/IaVEVC8G17uAGOMC9U0Ymez/Zj7DMOF4ebpf9QJNxJnbJOlfVQAAAABJRU5ErkJggg==',
NULL, 0),

--
-- Frequência com confirmação (log)
--

(DEFAULT, (SELECT id FROM idLog3), 'CRIACAO', (SELECT id FROM idFrq3), 1246, 114239, '.40',
'Imbel MD2', 'DL28461184', 102483466, 50,
NOW() - interval '5 hours',
NOW() - interval '4 hours' - interval '30 minutes',
NOW() - interval '4 hours' - interval '10 minutes', NULL, NULL,
'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAUklEQVR4XqWQ0QmAMAxEL8FhdIp2Hvep87hFxzm5D6GE0ip9kI/A5XHESCJfBzHgPqtZKjvxAe9cawbBCVtrEmFX///G9zKaFjtGc8T1TExQ5gHN9xsWe3/FugAAAABJRU5ErkJggg==',
NULL, 0),

(DEFAULT, (SELECT id FROM idLog4), 'ATUALIZACAO', (SELECT id FROM idFrq3), 1246, 114239, '.40',
'Imbel MD2', 'DL28461184', 102483466, 50,
NOW() - interval '5 hours',
NOW() - interval '4 hours' - interval '30 minutes',
NOW() - interval '4 hours' - interval '10 minutes',
NOW() - interval '4 hours' - interval '5 minutes',
NOW() - interval '4 hours' - interval '5 minutes',
'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAUklEQVR4XqWQ0QmAMAxEL8FhdIp2Hvep87hFxzm5D6GE0ip9kI/A5XHESCJfBzHgPqtZKjvxAe9cawbBCVtrEmFX///G9zKaFjtGc8T1TExQ5gHN9xsWe3/FugAAAABJRU5ErkJggg==',
'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAATElEQVR4XpWPCwoAIAhDswPrWTyxMcKEBn2EwYPG0yQi2st0gJllO5kHRlUNBIwsrkzjbnN3AdPqc5lXV/gMNqYt7cn9Vvr+NT0k7xkVBY1RndW3lwAAAABJRU5ErkJggg==',
1);