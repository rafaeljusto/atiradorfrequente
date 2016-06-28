package interceptador

import (
	"net"
	"net/http"
	"strings"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

type endereçoRemoto interface {
	EndereçoRemoto() net.IP
	DefineEndereçoRemoto(net.IP)
	EndereçoProxy() net.IP
	DefineEndereçoProxy(net.IP)
	Req() *http.Request
}

// EndereçoRemoto disponibiliza ao handler o endereço do cliente. Possuí a
// capacidade de tratar endereços enviados via proxy tanto com o cabeçalhos HTTP
// X-Forwarded-For quanto com o X-Real-IP.
type EndereçoRemoto struct {
	interceptor.NopInterceptor
	handler endereçoRemoto
}

// NovoEndereçoRemoto cria um novo interceptador EndereçoRemoto.
func NovoEndereçoRemoto(e endereçoRemoto) *EndereçoRemoto {
	return &EndereçoRemoto{handler: e}
}

// Before interpreta a conexão e os cabeçalhos HTTP para identificar o endereço
// IP do cliente. A prioridade é dada na seguinte ordem: IP da conexão, cabeçalhos HTTP
// X-Forwarded-For e X-Real-IP. Os cabeçalhos HTTP só serão analisados se o IP
// da conexão estiver na lista de proxies liberados no arquivo de configuração.
func (r *EndereçoRemoto) Before() int {
	endereçoCliente, _, err := net.SplitHostPort(r.handler.Req().RemoteAddr)

	if err != nil {
		log.Warning(erros.Novo(err))
		return http.StatusInternalServerError
	}

	r.handler.DefineEndereçoRemoto(net.ParseIP(strings.TrimSpace(endereçoCliente)))

	// se não houver configuração, ignora os proxies
	if config.Atual() == nil {
		return 0
	}

	éProxy := false
	for _, proxy := range config.Atual().Proxies {
		if r.handler.EndereçoRemoto().Equal(proxy) {
			éProxy = true
			break
		}
	}

	if !éProxy {
		return 0
	}

	endereçosProxies := r.handler.Req().Header.Get("X-Forwarded-For")
	endereçoReal := r.handler.Req().Header.Get("X-Real-IP")

	r.handler.DefineEndereçoProxy(r.handler.EndereçoRemoto())
	if ip := obtémEndereçoCliente(endereçosProxies, endereçoReal); ip != nil {
		r.handler.DefineEndereçoRemoto(ip)
	}

	return 0
}

func obtémEndereçoCliente(endereçosProxies, endereçoReal string) net.IP {
	endereçosProxies = strings.TrimSpace(endereçosProxies)
	for _, cliente := range strings.Split(endereçosProxies, ",") {
		cliente = strings.TrimSpace(cliente)
		if ip := net.ParseIP(cliente); ip != nil {
			return ip
		}
	}

	return net.ParseIP(strings.TrimSpace(endereçoReal))
}

// EndereçoRemotoCompatível implementa os métodos que serão utilizados pelo
// handler para acessar o endereço IP remoto armazenado por este interceptador.
type EndereçoRemotoCompatível struct {
	endereçoRemoto net.IP
	endereçoProxy  net.IP
}

// DefineEndereçoRemoto armazena o endereço IP remoto.
func (r *EndereçoRemotoCompatível) DefineEndereçoRemoto(ip net.IP) {
	r.endereçoRemoto = ip
}

// EndereçoRemoto obtém o endereço IP remoto.
func (r EndereçoRemotoCompatível) EndereçoRemoto() net.IP {
	return r.endereçoRemoto
}

// DefineEndereçoProxy armazena o endereço IP do último proxy. Somente será
// considerado proxy quando o proxy estiver liberado no arquivo de configuração
// do servidor REST.
func (r *EndereçoRemotoCompatível) DefineEndereçoProxy(ip net.IP) {
	r.endereçoProxy = ip
}

// EndereçoProxy obtém o endereço IP do último proxy. Somente será
// considerado proxy quando o proxy estiver liberado no arquivo de configuração
// do servidor REST.
func (r EndereçoRemotoCompatível) EndereçoProxy() net.IP {
	return r.endereçoProxy
}
