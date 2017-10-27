package rabbitmq

import (
	"github.com/parnurzeal/gorequest"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/agirot/superviGo/config"
	"time"

	"github.com/agirot/superviGo/ressource"
)

type rabbitMqRequest struct {
	MessageReady int          `json:"messages"`
	Name         string       `json:"name"`
	MessageStats messageStats `json:"message_stats"`
	AlertLevel   int
}

type messageStats struct {
	DeliverStats deliverStats `json:"deliver_get_details"`
}

type deliverStats struct {
	Rate float32 `json:"rate"`
}

var (
	Tab        = [][]string{}
	nextUpdate = time.Now()
)

func getStats(queueName string) (rabbitMqRequest, []error) {
	queueConfig, ok := config.Config.RabbitMq.Queue[queueName]
	if !ok {
		panic("Queue name undefined in configuration")
	}

	request := gorequest.New()
	currentQueue := rabbitMqRequest{}
	url := fmt.Sprintf("http://%v:%v/api/queues/%%2F/%v", config.Config.RabbitMq.Host, config.Config.RabbitMq.WebPort, queueName)
	_, _, errs := request.Get(url).SetBasicAuth(config.Config.RabbitMq.User, config.Config.RabbitMq.Password).EndStruct(&currentQueue)
	if errs != nil {
		return currentQueue, errs
	}
	if currentQueue.MessageReady >= queueConfig.AlertCriticMessage {
		currentQueue.AlertLevel = ressource.CriticalLevel
	} else if (currentQueue.MessageReady >= queueConfig.AlertWarningMessage) {
		currentQueue.AlertLevel = ressource.WarningLevel
	}
	return currentQueue, nil
}

func Render() *ui.Table {
	for name, _ := range config.Config.RabbitMq.Queue {
		Tab = append(Tab, []string{name, "loading...", "loading..."})
	}

	table2 := ui.NewTable()
	table2.BorderLabel = "RabbitMQ"
	table2.Rows = Tab
	table2.FgColor = ui.ColorWhite
	table2.BgColor = ui.ColorDefault
	table2.TextAlign = ui.AlignLeft
	table2.Separator = false
	table2.Analysis()
	table2.SetSize()
	table2.Y = 0
	table2.X = 0
	table2.Border = true

	return table2
}

func formatStringStatsCell(req rabbitMqRequest) string {
	return fmt.Sprintf("%9.f messages/s", req.MessageStats.DeliverStats.Rate)
}

func formatStringCountMessageCell(req rabbitMqRequest) string {
	return fmt.Sprintf("%v messages", req.MessageReady)
}

func UpdateRender(tab *ui.Table) {
	now := time.Now()
	if now.Before(nextUpdate) {
		return
	}

	for index, req := range Tab {
		stats, errs := getStats(req[0])
		if errs != nil {
			req[1] = fmt.Sprintf("%v", ressource.ErrorToSliceString(errs))
			break
		}

		if stats.Name == "" {
			req[1] = "Queue not found"
			req[2] = ""
			continue
		}

		req[1] = formatStringCountMessageCell(stats)
		req[2] = formatStringStatsCell(stats)
		if stats.AlertLevel == ressource.WarningLevel && tab.BgColor != ressource.CriticalColor {
			tab.BgColors[index] = ressource.WarningColor
			tab.FgColors[index] = ui.ColorBlack
		} else if stats.AlertLevel == ressource.CriticalLevel {
			tab.BgColors[index] = ressource.CriticalColor
			tab.FgColors[index] = ui.ColorWhite
		} else {
			tab.BgColors[index] = ui.ColorDefault
			tab.FgColors[index] = ui.ColorWhite
		}
	}

	nextUpdate = nextUpdate.Add(10 * time.Second)
}
