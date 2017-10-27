package main

import (
	ui "github.com/gizak/termui"

	"github.com/agirot/superviGo/rabbitmq"
	"github.com/agirot/superviGo/config"
	"github.com/agirot/superviGo/sentry"
	"github.com/agirot/superviGo/redis"
	"github.com/agirot/superviGo/http_status"
)

func main() {
	config.HydrateConfiguration()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	//Redis supervision
	redisWindow := redis.Render()

	//RabbitMQ supervision
	rabbitWindow := rabbitmq.Render()

	//Sentry supervision
	sentryWindow := sentry.Render()

	//Http status
	httpWindow := httpStatus.Render()

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, sentryWindow),
			ui.NewCol(6, 0, rabbitWindow, redisWindow, httpWindow),
			),
		)

	// calculate layout
	ui.Body.Align()
	ui.Render(ui.Body)

	//Event list
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		rabbitmq.UpdateRender(rabbitWindow)
		sentry.UpdateRender(sentryWindow)
		redis.UpdateRender(redisWindow)
		httpStatus.UpdateRender(httpWindow)
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
