package main

import (
	"fmt"
	"time"
)

type Token struct {
	Data      string
	Recipient int
	TTL       int
}

func main() {
	var N, ttl, recipient int
	fmt.Print("Введите число узлов: ")
	fmt.Scan(&N)

	channels := make([]chan Token, N)
	for i := range channels {
		channels[i] = make(chan Token)
	}

	fmt.Print("Введите время жизни токена: ")
	fmt.Scan(&ttl)

	for i := 0; i < N; i++ {
		go node(i, channels[i], channels[(i+1)%N])
	}

	for {
		fmt.Print("Введите адресата сообщения (c 0 по ", N-1, "): ")
		fmt.Scan(&recipient)

		if recipient < 0 || recipient >= N {
			fmt.Println("Адресат сообщения не попадает в доступный интервал. Пожалуйста, введите корректное значение.")
		} else {
			break
		}
	}

	var message string
	fmt.Print("Введите сообщение: ")
	fmt.Scan(&message)

	channels[0] <- Token{Data: message, Recipient: recipient, TTL: ttl}

	time.Sleep(time.Second)

	fmt.Println("Программа завершена")
}

func node(id int, in chan Token, out chan Token) {
	for {
		token := <-in
		//fmt.Printf("Узел %d получил токен, %v\n", id, token)
		fmt.Printf("Узел %d получил токен, id: %d. Остаток времени жизни: %d\n", id, token.Recipient, token.TTL)

		if token.Recipient == id {
			fmt.Printf("Ура! Узел %d достиг адресата и доставил сообщение: %s\n", id, token.Data)
			return
		} else if token.TTL > 0 {
			token.TTL--
			out <- token
			time.Sleep(100)
		} else {
			fmt.Printf("Эх :( На узле %d истекло время жизни токена\n", id)
			return
		}
	}
}
