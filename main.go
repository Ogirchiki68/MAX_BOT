package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	// Правильное создание клиента - функция New возвращает два значения
	api, err := maxbot.New("f9LHodD0cOLarUT0WyQcxJ-McEx2oLM8Q2EFbWisYXQLZlzc2MQuJJSbUnh3MzNFP0Ign9HbNJMlBUUefpv8")
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}

	// Получение информации о боте
	info, err := api.Bots.GetBot(context.Background())
	if err != nil {
		log.Printf("Ошибка получения информации о боте: %v", err)
	} else {
		fmt.Printf("Информация о боте: %#v\n", info)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Обработка сигналов завершения
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt)
		<-exit
		cancel()
	}()

	// Чтение обновлений
	for upd := range api.GetUpdates(ctx) {
		switch upd := upd.(type) {
		case *schemes.MessageCreatedUpdate:
			fmt.Printf("Получено сообщение: %s\n", upd.Message.Body.Text)

			// Создаем сообщение с помощью NewMessage
			message := maxbot.NewMessage().
				SetChat(upd.Message.Recipient.ChatId).
				SetText("Привет! Я получил ваше сообщение: " + upd.Message.Body.Text)

			// Отправка сообщения
			messageID, err := api.Messages.Send(context.Background(), message)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			} else {
				fmt.Printf("Сообщение отправлено с ID: %s\n", messageID)
			}
		}
	}
}
