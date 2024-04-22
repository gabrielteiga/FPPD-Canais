// Desenvolvido por: Gabriel Teiga
// Ultima modificação: 22/04/2024
//
// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NJ = 5 // numero de jogadores
const M = 4  // numero de cartas

type carta string // carta é um strirng

var ch [NJ]chan carta          // NJ canais de itens tipo carta
var bateu = make(chan int, NJ) // canal para bater

func jogador(id int, in chan carta, out chan carta, bateu chan int, cartasIniciais []carta) {
	mao := cartasIniciais // estado local - as cartas na mao do jogador
	nroDeCartas := M      // quantas cartas ele tem

	if id == 0 {
		mao = append(mao, carta("Joker"))
		nroDeCartas++
	}

	fmt.Println("Jogador", id, "começa com", mao)

	for {
		select {
		case <-bateu:
			fmt.Printf("Jogador %d bateu\n", id)
			bateu <- id
			return
		default:
			if readyToWin(mao) {
				fmt.Printf("Jogador %d ganhou\n", id)
				fmt.Println(mao)
				bateu <- id
				return
			} else if nroDeCartas > M {
				randomIndex := rand.Intn(len(mao))
				cartaParaSair := mao[randomIndex]
				mao = append(mao[:randomIndex], mao[randomIndex+1:]...)

				fmt.Printf("%d passou %s\n", id, cartaParaSair)
				out <- cartaParaSair
				nroDeCartas--
			} else {
				select {
				case cartaRecebida := <-in:
					fmt.Printf("%d recebeu %s\n", id, cartaRecebida)
					mao = append(mao, cartaRecebida)
					nroDeCartas++
				case <-bateu:
					fmt.Printf("Jogador %d bateu\n", id)
					bateu <- id
					return
				}
			}
		}
	}
}

func readyToWin(mao []carta) bool {
	if len(mao) < 4 {
		return false
	}
	return mao[0] == mao[1] && mao[1] == mao[2] && mao[2] == mao[3]
}

func main() {
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	// cria um baralho com NJ*M cartas
	baralho := make([]carta, NJ*M)
	for i := 0; i < M; i++ {
		baralho[i] = carta("A")
		baralho[i+4] = carta("B")
		baralho[i+8] = carta("C")
		baralho[i+12] = carta("D")
		baralho[i+16] = carta("E")
	}

	// embaralha o baralho
	rand.Shuffle(len(baralho), func(i, j int) {
		baralho[i], baralho[j] = baralho[j], baralho[i]
	})

	// distribui as cartas para os jogadores
	for i := 0; i < NJ; i++ {
		cartasEscolhidas := baralho[i*M : (i+1)*M]
		go jogador(i, ch[i], ch[(i+1)%NJ], bateu, cartasEscolhidas)
	}

	// <-make(chan struct{}) // bloqueia
	time.Sleep(2 * time.Second)
}
