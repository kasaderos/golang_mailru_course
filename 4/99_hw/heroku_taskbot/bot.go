package main

// сюда писать код

import (
	"fmt"
	"net/http"
	"os"
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

func caseBackslashTasks(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if len(tasks) == 0 {
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Нет задач",
		))
	} else {
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
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			s,
		))
	}
}

func addNewTask(bot *tgbotapi.BotAPI, update tgbotapi.Update, text string) {
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
}

func assignUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, cmd string) {
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
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"can't find task with task id "+strconv.Itoa(id),
		))
	} else {
		if tasks[ind].eUser != nil {
			bot.Send(tgbotapi.NewMessage(
				tasks[ind].eUser.chatId,
				"Задача \""+tasks[ind].text+"\" назначена на @"+update.Message.From.UserName,
			))
		} else {
			bot.Send(tgbotapi.NewMessage(
				tasks[ind].vendor.chatId,
				"Задача \""+tasks[ind].text+"\" назначена на @"+update.Message.From.UserName,
			))
		}
		tasks[ind].eUser = &User{
			chatId: update.Message.Chat.ID,
			name:   update.Message.From.UserName,
		}
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Задача \""+tasks[ind].text+"\" назначена на вас",
		))
	}
}

func unassignUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, cmd string) {
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
}

func resolveTask(bot *tgbotapi.BotAPI, update tgbotapi.Update, cmd string) {
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
	if L > 0 {
		tasks[L-1], tasks[ind] = tasks[ind], tasks[L-1]
		tasks = tasks[:L-1]
	}
}
func backslashMyCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	for _, t := range tasks {
		if t.eUser.name == update.Message.From.UserName {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				strconv.Itoa(t.id)+". "+t.text+" by @"+
					t.vendor.name+"\n/unassign_"+strconv.Itoa(t.id)+
					" /resolve_"+strconv.Itoa(t.id),
			))
		}
	}
}
func backslashOwnerCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	for _, t := range tasks {
		if t.vendor.name == update.Message.From.UserName {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				strconv.Itoa(t.id)+". "+t.text+" by @"+
					t.vendor.name+"\n/assign_"+strconv.Itoa(t.id),
			))
		}
	}
}
func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	go http.ListenAndServe(":"+port, nil)
	fmt.Println("start listen :8080")

	updates := bot.ListenForWebhook("/")
	for update := range updates {
		cmd, text := getCmdAndText(update.Message.Text)
		if cmd == "" {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Bad command",
			))
		}
		if cmd == "/tasks" {
			caseBackslashTasks(bot, update)
		} else if cmd == "/new" {
			addNewTask(bot, update, text)
		} else if strings.HasPrefix(cmd, "/assign_") {
			assignUser(bot, update, cmd)
		} else if strings.HasPrefix(cmd, "/unassign_") {
			unassignUser(bot, update, cmd)
		} else if strings.HasPrefix(cmd, "/resolve_") {
			resolveTask(bot, update, cmd)
		} else if cmd == "/my" {
			backslashMyCommand(bot, update)
		} else if cmd == "/owner" {
			backslashOwnerCommand(bot, update)
		} else {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"/tasks - выводит список всех активных задач\n"+
					"/new <название> - создаёт новую задачу\n"+
					"/assign_* - назначает задачу на себя\n"+
					"/unassign_* - снимает задачу с себя\n"+
					"/resolve_* завершает задачу, удаляет её из хранилища\n"+
					"/my показывает задачи которые назначены на меня\n"+
					"/owner - показывает задачи, которы я создал\n",
			))
		}
	}
}
