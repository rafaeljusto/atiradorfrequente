// Package servidor inicializa o servidor REST e se conecta com os componentes
// necessários.
package servidor

import (
	"crypto/tls"
	"net"
	"net/http"
	"runtime"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/handler"
	"github.com/registrobr/gostk/db"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy"
)

// Iniciar realiza todas as inicializações iniciais e sobe o servidor REST.
// Supõe que a configuração já foi carregada. Está função fica bloqueada
// enquanto o servidor estiver executando. Recebe como argumento a conexão TCP
// que esta escutando, podendo ser promovida a conexão TLS por está função.
func Iniciar(escuta net.Listener) error {
	if err := iniciarConexãoSyslog(); err != nil {
		log.Critf("Erro ao conectar servidor de log. Detalhes: %s", erros.Novo(err))
		return erros.Novo(err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			log.Errorf("Erro ao fechar a conexão do log. Detalhes: %s", erros.Novo(err))
		}
	}()

	// o sistema não é interrompido caso ocorra um problema de conexão com o banco
	// de dados. Novas tentativas serão feitas a cada tratamento de requisição.
	if err := iniciarConexãoBancoDados(); err != nil {
		log.Critf("Erro ao conectar o banco de dados. Detalhes: %s", erros.Novo(err))
	}
	defer func() {
		// TODO(rafaeljusto): mover esta verificação para o próprio objeto
		if bd.Conexão == nil {
			return
		}

		if err := bd.Conexão.Close(); err != nil {
			log.Errorf("Erro ao fechar a conexão do banco de dados. Detalhes: %s", erros.Novo(err))
		}
	}()

	if err := iniciarServidor(escuta); err != nil {
		log.Critf("Erro ao iniciar o servidor. Detalhes: %s", erros.Novo(err))
		return erros.Novo(err)
	}

	return nil
}

func iniciarConexãoSyslog() error {
	log.Info("Inicializando conexão com o servidor de log")

	return erros.Novo(log.Dial("tcp", config.Atual().Syslog.Endereço, "atirador-frequente", config.Atual().Syslog.TempoEsgotadoConexão))
}

func iniciarConexãoBancoDados() error {
	log.Info("Inicializando conexão com o banco de dados")

	err := bd.IniciarConexão(db.ConnParams{
		Username:           config.Atual().BancoDados.Usuário,
		Password:           config.Atual().BancoDados.Senha,
		DatabaseName:       config.Atual().BancoDados.Nome,
		Host:               config.Atual().BancoDados.Endereço,
		ConnectTimeout:     config.Atual().BancoDados.TempoEsgotadoConexão,
		StatementTimeout:   config.Atual().BancoDados.TempoEsgotadoComando,
		MaxIdleConnections: config.Atual().BancoDados.MáximoNúmeroConexõesInativas,
		MaxOpenConnections: config.Atual().BancoDados.MáximoNúmeroConexõesAbertas,
	}, config.Atual().BancoDados.TempoEsgotadoTransação)

	return erros.Novo(err)
}

func iniciarServidor(escuta net.Listener) error {
	log.Info("Inicializando servidor")

	handy.ErrorFunc = log.Error

	h := handy.NewHandy()
	h.Recover = func(r interface{}) {
		const tamanho = 1 << 16
		buffer := make([]byte, tamanho)
		buffer = buffer[:runtime.Stack(buffer, false)]
		log.Critf("Erro grave detectado. Detalhes: %v\n%s", r, buffer)
	}

	for rota, handler := range handler.Rotas {
		h.Handle(rota, handler)
	}

	servidor := http.Server{
		Handler:     h,
		ReadTimeout: config.Atual().Servidor.TempoEsgotadoLeitura,
	}

	if config.Atual().Servidor.TLS.Habilitado {
		certificado, err := tls.LoadX509KeyPair(config.Atual().Servidor.TLS.ArquivoCertificado, config.Atual().Servidor.TLS.ArquivoChave)
		if err != nil {
			return erros.Novo(err)
		}

		cifras := []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_FALLBACK_SCSV,
		}

		configuraçãoTLS := tls.Config{
			MinVersion:               tls.VersionTLS10,
			PreferServerCipherSuites: true,
			CipherSuites:             cifras,
			Certificates:             []tls.Certificate{certificado},
		}

		escuta = tls.NewListener(escuta, &configuraçãoTLS)
	}

	return erros.Novo(servidor.Serve(escuta))
}
