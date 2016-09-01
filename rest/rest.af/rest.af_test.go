package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
)

func Test_main(t *testing.T) {
	if os.Getenv("EXECUTAR_TESTE_REST_AF") == "1" {
		for i := len(os.Args) - 1; i >= 0; i-- {
			// remove o argumento utilizado para o ambiente de teste
			if os.Args[i] == "-test.run=Test_main" {
				os.Args = append(os.Args[:i], os.Args[i+1:]...)
				break
			}
		}

		main()
		return
	}

	loggerOriginal := log.LocalLogger
	defer func() {
		log.LocalLogger = loggerOriginal
	}()

	var servidorLog simulador.ServidorLog
	syslog, err := servidorLog.Executar("localhost:0")
	if err != nil {
		t.Fatalf("Erro ao inicializar o servidor de log. Detalhes: %s", err)
	}
	defer syslog.Close()

	cenários := []struct {
		descrição             string
		variáveisAmbiente     map[string]string
		arquivoConfiguração   string
		sucesso               bool
		mensagensEsperadas    *regexp.Regexp
		mensagensLogEsperadas *regexp.Regexp
	}{
		{
			descrição: "deve iniciar o servidor REST carregando as configurações de variáveis de ambiente",
			variáveisAmbiente: map[string]string{
				"AF_SERVIDOR_ENDERECO": "0.0.0.0:0",
				"AF_SYSLOG_ENDERECO":   syslog.Addr().String(),
			},
			sucesso:            true,
			mensagensEsperadas: regexp.MustCompile(`^$`),
			mensagensLogEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o banco de dados
.*Erro ao conectar o banco de dados. Detalhes: .*getsockopt: connection refused
.*Inicializando servidor
$`),
		},
	}

	for i, cenário := range cenários {
		cmd := exec.Command(os.Args[0], "-test.run=Test_main")
		cmd.Env = append(os.Environ(), "EXECUTAR_TESTE_REST_AF=1")

		for chave, valor := range cenário.variáveisAmbiente {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", chave, valor))
		}

		if cenário.arquivoConfiguração != "" {
			arquivoConfiguração, err := ioutil.TempFile("", "atirador-frequente-")
			if err != nil {
				t.Fatalf("Item %d, “%s”: erro ao criar o arquivo de configuração. Detalhes: %s",
					i, cenário.descrição, err)
			}

			arquivoConfiguração.WriteString(cenário.arquivoConfiguração)
			arquivoConfiguração.Close()

			cmd.Args = append(cmd.Args, []string{"--config", arquivoConfiguração.Name()}...)
		}

		var mensagens []byte

		go func() {
			mensagens, err = cmd.CombinedOutput()

			if erroSaída, ok := err.(*exec.ExitError); ok && erroSaída.Success() != cenário.sucesso {
				t.Errorf("Item %d, “%s”: resultado da execução inesperado. Resultado: %t",
					i, cenário.descrição, erroSaída.Success())
			} else if err != nil {
				t.Errorf("Item %d, “%s”: erro inesperado ao executar o servidor REST. Resultado: %s",
					i, cenário.descrição, err)
			}
		}()

		// aguarda os serviços serem executados
		time.Sleep(100 * time.Millisecond)

		if err := cmd.Process.Kill(); err != nil {
			// ignora o erro que informa que o processo já foi encerrado
			if err.Error() != "os: process already finished" {
				t.Fatalf("Item %d, “%s”: erro ao matar o servidor REST. Detalhes: %s",
					i, cenário.descrição, err)
			}
		}

		if !cenário.mensagensEsperadas.Match(mensagens) {
			t.Errorf("Item %d, “%s”: mensagem inesperada. Detalhes: %s",
				i, cenário.descrição, mensagens)
		}

		if !cenário.mensagensLogEsperadas.MatchString(servidorLog.Mensagens()) {
			t.Errorf("Item %d, “%s”: mensagens de log inesperadas. Detalhes: %s",
				i, cenário.descrição, servidorLog.Mensagens())
		}
	}
}
