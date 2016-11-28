package atirador

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/randômico"
	"golang.org/x/crypto/hkdf"
)

type frequência struct {
	ID                   int64
	Controle             int64
	CR                   int
	Calibre              string
	ArmaUtilizada        string
	NúmeroSérie          string
	GuiaDeTráfego        int
	QuantidadeMunição    int
	DataInício           time.Time
	DataTérmino          time.Time
	DataCriação          time.Time
	DataAtualização      time.Time
	DataConfirmação      time.Time
	ImagemNúmeroControle string
	ImagemConfirmação    string

	// revisão utilizado para o controle de versão do objeto na base de dados,
	// minimizando problemas de concorrência quando 2 transações alteram o mesmo
	// objeto.
	revisão int
}

func novaFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) frequência {
	return frequência{
		Controle:          randômico.FonteRandômica.Int63(),
		CR:                frequênciaPedidoCompleta.CR,
		Calibre:           frequênciaPedidoCompleta.Calibre,
		ArmaUtilizada:     frequênciaPedidoCompleta.ArmaUtilizada,
		NúmeroSérie:       frequênciaPedidoCompleta.NúmeroSérie,
		GuiaDeTráfego:     frequênciaPedidoCompleta.GuiaDeTráfego,
		QuantidadeMunição: frequênciaPedidoCompleta.QuantidadeMunição,
		DataInício:        frequênciaPedidoCompleta.DataInício,
		DataTérmino:       frequênciaPedidoCompleta.DataTérmino,
	}
}

func (f *frequência) confirmar(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) {
	f.DataConfirmação = time.Now().UTC()
	f.ImagemConfirmação = frequênciaConfirmaçãoPedidoCompleta.Imagem
}

func (f *frequência) gerarCódigoVerificação(chave string) string {
	buffer := new(bytes.Buffer)

	// o erro retornado nesta escrita é ignorado, pois o tipo bytes.Buffer não
	// gera erro no método Write
	binary.Write(buffer, binary.LittleEndian, int64(f.ID))

	derivaçãoChave := make([]byte, 32)
	funçãoDerivação := hkdf.New(sha256.New, []byte(chave), nil, buffer.Bytes())

	// o erro retornado nesta leitura é ignorado, pois como estamos lendo a
	// quantidade total de bytes o cenário de erro nunca será atingido
	io.ReadFull(funçãoDerivação, derivaçãoChave)

	mensagem := fmt.Sprintf("%010d %d %d", f.ID, f.CR, f.Controle)

	// o erro retornado é ignorado, pois o método Write do SHA256 não gera erro
	mac := hmac.New(sha256.New, derivaçãoChave)
	mac.Write([]byte(mensagem))
	mensagemCodificada := base58.Encode(mac.Sum(nil))

	if tamanho := 44 - len(mensagemCodificada); tamanho > 0 {
		mensagemCodificada += strings.Repeat("o", tamanho)
	}
	return mensagemCodificada
}

func (f frequência) protocoloPendente(códigoVerificação string) protocolo.FrequênciaPendenteResposta {
	return protocolo.FrequênciaPendenteResposta{
		NúmeroControle:    protocolo.NovoNúmeroControle(f.ID, f.Controle),
		CódigoVerificação: códigoVerificação,
		Imagem:            f.ImagemNúmeroControle,
	}
}
