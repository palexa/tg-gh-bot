package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	ctx context.Context
	ts  oauth2.TokenSource
	tc  *http.Client
}

func NewBot(token string, githubToken string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	t, err := ts.Token()

	t.AccessToken = githubToken

	tc := oauth2.NewClient(context.Background(), ts)

	return &Bot{
		bot: bot,
		ctx: context.Background(),
		ts:  ts,
		tc:  tc,
	}, nil
}

func (b *Bot) Start() error {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	updates := b.bot.GetUpdatesChan(config)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				b.handleHelp(update.Message)
			case "issues":
				b.handleIssues(update.Message)
			case "pullrequests":
				b.handlePullRequests(update.Message)
			default:
				b.handleUnknown(update.Message)
			}
		}
	}

	return nil
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	helpText := `Usage:
/issues - List open issues in the repository
/pullrequests - List open pull requests in the repository`

	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, helpText))
}

func (b *Bot) handleIssues(msg *tgbotapi.Message) {
	client := github.NewClient(b.tc)

	issues, _, err := client.Issues.ListByRepo(b.ctx, "owner", "repo", &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		log.Printf("Failed to list issues: %v", err)
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Failed to list issues"))
		return
	}

	if len(issues) == 0 {
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "No open issues found"))
		return
	}

	responseText := "Open issues:\n"
	for _, issue := range issues {
		responseText += fmt.Sprintf("#%d %s\n", *issue.Number, *issue.Title)
	}

	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, responseText))
}

func (b *Bot) handlePullRequests(msg *tgbotapi.Message) {
	client := github.NewClient(b.tc)

	pullRequests, _, err := client.PullRequests.List(b.ctx, "owner", "repo", &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		log.Printf("Failed to list pull requests: %v", err)
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Failed to list pull requests"))
		return
	}

	if len(pullRequests) == 0 {
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "No open pull requests found"))
		return
	}

	responseText := "Open pull requests:\n"
	for _, pr := range pullRequests {
		responseText += fmt.Sprintf("#%d %s\n", *pr.Number, *pr.Title)
	}

	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, responseText))
}

func (b *Bot) handleUnknown(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unknown command"))
}
