package redis

import (
	"time"
	"github.com/go-redis/redis"
	ui "github.com/gizak/termui"

	"github.com/agirot/superviGo/config"
	"strings"
)

var (
	Tab        = [][]string{}
	nextUpdate = time.Now()
	client     *redis.Client
)

func getStats() map[string]string {
	client = redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Host,
	})

	infoRaw := client.Info()
	infosSplit := strings.Split(infoRaw.String(), "\r\n")
	infoMap := make(map[string]string)
	for _, info := range infosSplit {
		split := strings.Split(info, ":")
		if len(split) > 1 {
			infoMap[split[0]] = split[1]
		}
	}

	return infoMap
}

func Render() *ui.Table {

	for _, info := range config.Config.Redis.Info {
		Tab = append(Tab, []string{info, "loading..."})
	}

	table2 := ui.NewTable()
	table2.BorderLabel = "Redis"
	table2.Rows = Tab
	table2.FgColor = ui.ColorWhite
	table2.BgColor = ui.ColorDefault
	table2.TextAlign = ui.AlignLeft
	table2.Separator = false
	table2.Analysis()
	table2.SetSize()
	table2.Y = 10
	table2.X = 0
	table2.Border = true

	return table2
}

/*func formatStringStatsCell(req rabbitMqRequest) string {
	return fmt.Sprintf("%9.f messages/s", req.MessageStats.DeliverStats.Rate)
}*/

func UpdateRender(tab *ui.Table) {
	now := time.Now()
	if now.Before(nextUpdate) {
		return
	}

	stats := getStats()
	for _, req := range Tab {
		value, ok := stats[req[0]]
		if ok {
			req[1] = value
		}
	}

	nextUpdate = nextUpdate.Add(10 * time.Second)
}
