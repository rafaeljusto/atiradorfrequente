INSERT INTO frequencia_atirador (id, controle, cr, calibre, arma_utilizada, numero_serie,
guia_de_trafego, quantidade_municao, data_inicio, data_termino, data_criacao, data_atualizacao,
data_confirmacao, imagem_numero_controle, imagem_confirmacao, revisao) VALUES

-- Frequência sem confirmação
(DEFAULT, 1234, 380308, '.380', 'Arma do Clube', 'HG72643653', 762556223, 100,
NOW() - interval '2 hour',
NOW() - interval '30 minutes',
NOW() - interval '29 minutes', NULL, NULL, '', NULL, 0),

-- Frequência sem confirmação expirada
(DEFAULT, 7344, 923714, '.45', 'Imbel 1911', 'SF9153921', 839201286, 150,
NOW() - interval '2 hour',
NOW() - interval '40 minutes',
NOW() - interval '31 minutes', NULL, NULL, '', NULL, 0),

-- Frequência com confirmação
(DEFAULT, 1246, 114239, '.40', 'Imbel MD2', 'DL28461184', 102483466, 50,
NOW() - interval '5 hours',
NOW() - interval '5 hours' - interval '30 minutes',
NOW() - interval '5 hours' - interval '10 minutes',
NOW() - interval '5 hours' - interval '10 minutes',
NOW() - interval '5 hours' - interval '5 minutes', '', '', 1);