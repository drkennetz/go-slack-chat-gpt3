// Package gptslack handles slack appMention events and responds with chat-gpt response
package gptslack

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

// EventHandler handles slack events
func EventHandler(appToken string, botToken string, gptClient gpt3.Client, ctx context.Context, log *zap.SugaredLogger) error {

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(zap.NewStdLog(log.Desugar())),
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(zap.NewStdLog(log.Desugar())),
	)
	socketmodeHandler := socketmode.NewSocketmodeHandler(client)
	// should be a primary middleware handler, and these handle more granular events
	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)
	socketmodeHandler.Handle(socketmode.EventTypeHello, middlewareHello)
	socketmodeHandler.HandleEvents(slackevents.AppMention, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareAppMentionEvent(evt, client, gptClient, ctx)
	})
	socketmodeHandler.HandleEvents(slackevents.Message, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareMessageEvent(evt, client, gptClient, ctx)
	})
	return socketmodeHandler.RunEventLoop()
}
