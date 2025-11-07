package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type ProductiveBot struct {
	api *maxbot.Api
}

// todo вынести в конфиг
func NewProductiveBot() (*ProductiveBot, error) {
	api, err := maxbot.New("f9LHodD0cOLarUT0WyQcxJ-McEx2oLM8Q2EFbWisYXQLZlzc2MQuJJSbUnh3MzNFP0Ign9HbNJMlBUUefpv8")
	if err != nil {
		return nil, err
	}
	return &ProductiveBot{api: api}, nil
}

func (b *ProductiveBot) Start() error {

	info, err := b.api.Bots.GetBot(context.Background())
	if err != nil {
		log.Printf("Ошибка получения информации о боте: %v", err)
	} else {
		fmt.Printf("Бот запущен: %s\n", info.Name)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt)
		<-exit
		fmt.Println("Завершение работы...")
		cancel()
	}()

	fmt.Println("Ожидание сообщений...")
	for upd := range b.api.GetUpdates(ctx) {
		switch upd := upd.(type) {
		case *schemes.MessageCreatedUpdate:
			b.handleMessage(upd)
		}
	}

	return nil
}

// todo сделать через enum комнады, убрать стрингу из параметров
func (b *ProductiveBot) handleMessage(update *schemes.MessageCreatedUpdate) {
	userMessage := strings.ToLower(strings.TrimSpace(update.Message.Body.Text))
	chatID := update.Message.Recipient.ChatId
	userID := update.Message.Sender.UserId

	fmt.Printf("Сообщение от пользователя %s: %s\n", userID, userMessage)

	response := b.generateResponse(userMessage)

	message := maxbot.NewMessage().
		SetChat(chatID).
		SetText(response)

	messageID, err := b.api.Messages.Send(context.Background(), message)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	} else {
		fmt.Printf("Ответ отправлен с ID: %s\n", messageID)
	}
}

// todo вынести в хенделры
func (b *ProductiveBot) generateResponse(message string) string {
	switch {
	case strings.Contains(message, "привет") || strings.Contains(message, "старт") || strings.Contains(message, "начать"):
		return `👋 Привет! Я бот для повышения продуктивности молодого поколения!

Что вас интересует?
• Напишите "продуктивность" для советов
• Напишите "тайм-менеджмент" для методов управления временем
• Напишите "фокус" для советов по концентрации`

	case strings.Contains(message, "продуктивность"):
		return `🎯 **Советы для повышения продуктивности:**

1. **Метод Pomodoro** 🍅 - 25 минут работы, 5 минут отдыха
2. **Список задач** 📝 - планируйте день с утра
3. **Цели SMART** 🎯 - Specific, Measurable, Achievable, Relevant, Time-bound
4. **Digital Detox** 📵 - ограничьте соцсети во время работы
5. **Здоровый сон** 💤 - 7-9 часов для восстановления

Какой метод хотите попробовать первым?`

	case strings.Contains(message, "тайм-менеджмент") || strings.Contains(message, "время"):
		return `⏰ **Методы тайм-менеджмента:**

• **Матрица Эйзенхауэра** - разделение на срочное/важное
• **Eat the Frog** 🐸 - начинайте с самой сложной задачи
• **Time Blocking** 🗓️ - планирование времени по блокам
• **Правило 2 минут** ⏱️ - если дело занимает <2 мин, делайте сразу

Что из этого вам ближе?`

	case strings.Contains(message, "фокус") || strings.Contains(message, "концентрация"):
		return `🎧 **Советы для улучшения фокуса:**

1. **Рабочая атмосфера** 🏢 - найдите тихое место
2. **Убрать отвлекающие факторы** 📵 - отключите уведомления
3. **Музыка для фокуса** 🎵 - инструментальная или ambient
4. **Техника 90/30** ⏰ - 90 минут работы, 30 минут отдыха
5. **Медитация** 🧘‍♂️ - 5 минут для очистки мыслей

Попробуйте начать с одного совета!`

	case strings.Contains(message, "помощь") || strings.Contains(message, "команды"):
		return `📚 **Доступные команды:**

• "продуктивность" - советы по продуктивности
• "тайм-менеджмент" - методы управления временем  
• "фокус" - улучшение концентрации
• "помощь" - показать это сообщение

Или просто напишите ваш вопрос!`

	default:
		return `🤔 Я не совсем понял ваш вопрос. Я специализируюсь на продуктивности молодого поколения!

Попробуйте написать:
• "продуктивность" - для общих советов
• "тайм-менеджмент" - для управления временем
• "фокус" - для улучшения концентрации
• "помощь" - для списка команд`
	}
}

func main() {
	bot, err := NewProductiveBot()
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Ошибка работы бота: %v", err)
	}
}
