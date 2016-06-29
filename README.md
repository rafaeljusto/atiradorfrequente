![Atirador Frequente](https://raw.githubusercontent.com/rafaeljusto/atiradorfrequente/master/logo.png)

Sistema de controle de frequência para Atiradores com CR (Certificado de
Registro) em Clubes de Tiro proposto ao Exército Brasileiro.

## English Introduction

This project was developed for the Brazillian Army to make it easier to control
shooter frequencies in shooting ranges. The whole content is written in
Portuguese as there's no objective to use it outside from Brazil.

## Proposta

A ideia central é desenvolver um sistema em que os Clubes de Tiro devem se
reportar sempre que um Atirador frequente o estande. Isso acabaria com os atuais
livros de registro, detectaria automaticamente Atiradores que não frequentam e
reduziria a possibilidade de fraude.

O funcionamento seria da seguinte forma:

1. Após o término do treino do Atirador, o Clube de Tiro acessa o Sistema do
Exército e informa o CR do Atirador, o calibre utilizado, quantidade de
munições, horários de inicio/fim e números de controle da arma de fogo
utilizada;
2. O Sistema do Exército cadastra esta informação e gera um número de controle;
3. O Clube de Tiro imprime o número de controle e uma folha sulfite;
4. O Clube de Tiro tira uma foto digital do Atirador exibindo o número de
controle;
5. A foto digital é enviada para o Sistema do Exército em até 30 minutos, que é
associada aos dados inicialmente enviados.

Obrigando o Atirador a tirar uma foto com um número gerado pelo sistema do
Exército, garantimos que ele realmente esta presente no Clube de Tiro naquele
momento, o que diminui os problemas de fraude. Este Sistema do Exército poderia
automaticamente identificar em um período quais Atiradores não tiveram a
frequência mínima necessária para manter o CR, não exigindo mais a análise
manual do livro de registro do Clube de Tiro. Uma auditoria poderia pegar por
amostragem fotos enviadas por um Clube de Tiro e analisar se são realmente do
Atirador reportado com o número de controle gerado.

Este Sistema poderia ser isolado dos demais Sistemas do Exército, já que não
precisa de nenhuma informação externa para funcionar. O lado negativo desta
solução, é que exige que o Clube de Tiro tenha acesso a Internet, uma impressora
e uma máquina fotográfica digital (ou celular com foto). Mas acredito que poucos
Clubes de Tiro não possuem estes requisitos.

## Serviços

Abaixo a lista de serviços a serem implementadas neste projeto. Conforme
surgirem novos serviços ou os listados forem concluídos esta tabela será
alterada.

| Descrição                            | REST                  | WEB                   | URI                                         |
| ------------------------------------ | :-------------------: | :-------------------: | ------------------------------------------- |
| Criar uma freqência (clube)          | :white_check_mark:    | :white_medium_square: | /frequencia/{cr} **[POST]**                 |
| Confirmar uma frequência (clube)     | :white_check_mark:    | :white_medium_square: | /frequencia/{cr}/{numeroControle} **[PUT]** |
| Cadastrar um clube (administrativo?) | :white_medium_square: | :white_medium_square: | /clube **[POST]**                           |
| Login (clube e administrativo)       | :white_medium_square: | :white_medium_square: | /login **[POST]**                           |
| Listar frequências (administrativo)  | :white_medium_square: | :white_medium_square: | /frequencia **[GET]**                       |
