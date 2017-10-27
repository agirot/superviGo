package httpStatus

import (
	"time"
	ui "github.com/gizak/termui"
	"github.com/agirot/superviGo/config"
	"github.com/parnurzeal/gorequest"
	"github.com/agirot/superviGo/ressource"
	"fmt"
)

var (
	Tab        = [][]string{}
	nextUpdate = time.Now()
)

func getStatus(endPoint ressource.HttpEndpoint) (gorequest.Response, time.Duration, []error) {
	req := gorequest.New().Timeout(2*time.Second)
	req.Url = endPoint.Url
	req.Method = "GET"
	if endPoint.BasicAuthUser != "" || endPoint.BasicAuthPassword != "" {
		req.BasicAuth.Username = endPoint.BasicAuthUser
		req.BasicAuth.Password = endPoint.BasicAuthPassword
	}

	now := time.Now()
	resp, _, errs := req.End()
	timing := time.Since(now)
	if errs != nil {
		return nil, timing, errs
	}
	return resp, timing, nil
}

func Render() *ui.Table {

	for name, _ := range config.Config.HttpStatus {
		Tab = append(Tab, []string{name, "loading..."})
	}

	table2 := ui.NewTable()
	table2.BorderLabel = "HTTP STATUS"
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

func UpdateRender(tab *ui.Table) {
	now := time.Now()
	if now.Before(nextUpdate) {
		return
	}

	for index, endpoint := range Tab {
		endPoint := config.Config.HttpStatus[endpoint[0]]
		response, latency, errs := getStatus(endPoint)
		if errs != nil {
			endpoint[1] = fmt.Sprintf("%v", ressource.ErrorToSliceString(errs))
			continue
		}

		endpoint[1] = fmt.Sprintf("Return %v (%v)", response.StatusCode, latency)
		if endPoint.HttpCodeExpected != 0 && response.StatusCode != endPoint.HttpCodeExpected {
			tab.BgColors[index] = ressource.CriticalColor
			tab.FgColors[index] = ui.ColorWhite
		}
	}

	nextUpdate = nextUpdate.Add(10 * time.Second)
}