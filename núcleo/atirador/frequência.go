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
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
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

func (f *frequência) gerarCódigoVerificação(chave string) (string, error) {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.LittleEndian, int64(f.ID)); err != nil {
		return "", erros.Novo(err)
	}

	derivaçãoChave := make([]byte, 32)
	funçãoDerivação := hkdf.New(sha256.New, []byte(chave), nil, buffer.Bytes())
	if _, err := io.ReadFull(funçãoDerivação, derivaçãoChave); err != nil {
		return "", erros.Novo(err)
	}

	mensagem := fmt.Sprintf("%010d %d %d", f.ID, f.CR, f.Controle)
	mac := hmac.New(sha256.New, derivaçãoChave)
	if _, err := mac.Write([]byte(mensagem)); err != nil {
		return "", erros.Novo(err)
	}
	mensagemCodificada := base58.Encode(mac.Sum(nil))

	if tamanho := 44 - len(mensagemCodificada); tamanho > 0 {
		mensagemCodificada += strings.Repeat("o", tamanho)
	}
	return mensagemCodificada, nil
}

func (f frequência) protocoloPendente(códigoVerificação string) protocolo.FrequênciaPendenteResposta {
	return protocolo.FrequênciaPendenteResposta{
		NúmeroControle:    protocolo.NovoNúmeroControle(f.ID, f.Controle),
		CódigoVerificação: códigoVerificação,
		Imagem:            f.ImagemNúmeroControle,
	}
}
