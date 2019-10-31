package main

// сюда писать код

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "1038691213:AAHpsEnXhBEW0QDqQQ_vD5sgz52TCA4XBg8"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://cool-taskbot.herokuapp.com"
	tasks      []Task
	autoinc    = 1
)

type User struct {
	name   string
	chatId int64
	vendor bool
}

type Task struct {
	id         int
	text       string
	eUser      *User
	isResolved bool
	vendor     *User
}

func findById(id int) int {
	for i, task := range tasks {
		if task.id == id {
			return i
		}
	}
	return -1
}

func isAssignedMe(m string, users []User) bool {
	for _, u := range users {
		if m == u.name && !u.vendor {
			return true
		}
	}
	return false
}

func removeFromTaskUsername(v string, users []User) ([]User, bool) {
	for i, u := range users {
		if u.name == v {
			L := len(users)
			users[L-1], users[i] = users[i], users[L-1]
			return users[:L-1], false
		}
	}
	return users, true
}

func getCmdAndText(s string) (string, string) {
	v := strings.Split(s, " ")
	if v == nil {
		return "", ""
	} else if len(v) > 1 {
		return v[0], strings.Join(v[1:], " ")
	} else if len(v) == 1 {
		return v[0], ""
	}
	return "", ""
}

func sendMessage(bot *tgbotapi.BotAPI, chatId int64, message string) {
	bot.Send(tgbotapi.NewMessage(
		chatId,
		message,
	))
}

func caseBackslashTasks(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if len(tasks) == 0 {
		sendMessage(bot, update.Message.Chat.ID, "Нет задач")
	} else {
		fmt.Println(tasks)
		s := ""
		for i, task := range tasks {
			id := strconv.Itoa(task.id)
			if update.Message.From.UserName == task.eUser.name && !task.eUser.vendor {
				s += id + ". " + task.text + " by " + "@" + task.vendor.name +
					"\nassignee: я\n/unassign_" + id + " /resolve_" + id
			} else if task.eUser.name != task.vendor.name {
				s += strconv.Itoa(task.id) + ". " + task.text + " by " + "@" + task.vendor.name +
					"\nassignee: @" + task.eUser.name
			} else {
				s += strconv.Itoa(task.id) + ". " + task.text + " by " + "@" + task.vendor.name +
					"\n/assign_" + strconv.Itoa(task.id)
			}
			if i < len(tasks)-1 {
				s += "\n\n"
			}
		}
		sendMessage(bot, update.Message.Chat.ID, s)
	}
}

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	go http.ListenAndServe(":8081", nil)

	updates := bot.ListenForWebhook("/")
	for update := range updates {
		cmd, text := getCmdAndText(update.Message.Text)
		if cmd == "" {
			sendMessage(bot, update.Message.Chat.ID, "Bad command")

		}
		if cmd == "/tasks" {
			caseBackslashTasks(bot, update)
		} else if cmd == "/new" {
			t := Task{
				id:   autoinc,
				text: text,
			}
			autoinc++
			t.vendor = &User{
				name:   update.Message.From.UserName,
				chatId: update.Message.Chat.ID,
				vendor: true,
			}
			t.eUser = t.vendor
			tasks = append(tasks, t)
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Задача \""+t.text+"\" создана, id="+strconv.Itoa(t.id),
			))
		} else if strings.HasPrefix(cmd, "/assign_") {
			cmds := strings.Split(cmd, "_")
			id, err := strconv.Atoi(cmds[1])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"assign with unknown task id",
				))
			}
			ind := findById(id)
			if ind == -1 {
				sendMessage(bot, update.Message.Chat.ID,
					"can't find task with task id "+strconv.Itoa(id))
			} else {
				// уведомляем
				if tasks[ind].eUser != nil {
					s := "Задача \"" + tasks[ind].text + "\" назначена на @" + update.Message.From.UserName
					sendMessage(bot, tasks[ind].eUser.chatId, s)
				} else {
					s := "Задача \"" + tasks[ind].text + "\" назначена на @" + update.Message.From.UserName
					sendMessage(bot, tasks[ind].vendor.chatId, s)
				}
				tasks[ind].eUser = &User{
					chatId: update.Message.Chat.ID,
					name:   update.Message.From.UserName,
				}
				s := "Задача \"" + tasks[ind].text + "\" назначена на вас"
				sendMessage(bot, update.Message.Chat.ID, s)

			}
		} else if strings.HasPrefix(cmd, "/unassign_") {
			cmds := strings.Split(cmd, "_")
			id, err := strconv.Atoi(cmds[1])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"unassign with unknown task id",
				))
			}
			ind := findById(id)
			if tasks[ind].eUser.name != update.Message.From.UserName {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Задача не на вас",
				))
			} else {
				tasks[ind].eUser = nil
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Принято",
				))
				bot.Send(tgbotapi.NewMessage(
					tasks[ind].vendor.chatId,
					"Задача \""+tasks[ind].text+"\" осталась без исполнителя",
				))
			}
		} else if strings.HasPrefix(cmd, "/resolve_") {
			cmds := strings.Split(cmd, "_")
			id, err := strconv.Atoi(cmds[1])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"resolve with unknown task id",
				))
			}
			ind := findById(id)
			bot.Send(tgbotapi.NewMessage(
				tasks[ind].vendor.chatId,
				"Задача \""+tasks[ind].text+"\" выполнена @"+update.Message.From.UserName,
			))

			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Задача \""+tasks[ind].text+"\" выполнена",
			))
			L := len(tasks)
			if L >= 1 {
				tasks[L-1], tasks[ind] = tasks[ind], tasks[L-1]
				tasks = tasks[:L-1]
			}
		} else if cmd == "/my" {
			for _, t := range tasks {
				if t.eUser.name == update.Message.From.UserName {
					s := strconv.Itoa(t.id) + ". " + t.text + " by @" +
						t.vendor.name + "\n/unassign_" + strconv.Itoa(t.id) +
						" /resolve_" + strconv.Itoa(t.id)
					sendMessage(bot, update.Message.Chat.ID, s)
				}
			}
		} else if cmd == "/owner" {
			for _, t := range tasks {
				if t.vendor.name == update.Message.From.UserName {
					s := strconv.Itoa(t.id) + ". " + t.text + " by @" +
						t.vendor.name + "\n/assign_" + strconv.Itoa(t.id)
					sendMessage(bot, update.Message.Chat.ID, s)
				}
			}
		}
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
