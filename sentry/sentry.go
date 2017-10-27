package sentry

import (
	ui "github.com/gizak/termui"
	"github.com/parnurzeal/gorequest"
	"fmt"

	"github.com/agirot/superviGo/config"
	"time"
	"github.com/agirot/superviGo/ressource"
	"strconv"
)

var (
	nextUpdate    = time.Now()
	previousAlert = make(map[string]int, config.Config.Sentry.MaxItem)
)

type SentryList struct {
	SentryItem []SentryItem
}

type SentryItem struct {
	Count     string `json:"count"`
	Title     string `json:"title"`
	Level     string `json:"level"`
	Culprit   string `json:"culprit"`
	FirstSeen string `json:"firstSeen"`
	LastSeen  string `json:"lastSeen"`
	ShareId   string `json:"shareId"`
}

func Render() *ui.List {
	if config.Config.Sentry.MaxItem < 0 {
		config.Config.Sentry.MaxItem = 2
	} else {
		config.Config.Sentry.MaxItem = config.Config.Sentry.MaxItem + 1
	}
	items := []string{"Loading..."}
	list := ui.NewList()
	list.Items = items
	list.BorderLabel = "Sentry"
	list.Height = config.Config.Sentry.MaxItem * 2
	list.Width = 25
	list.Y = 0

	return list
}

func getStats() (SentryList, []error) {
	request := gorequest.New()
	sentryConfig := config.Config.Sentry
	data := []SentryItem{}
	url := fmt.Sprintf("http://%v/api/0/projects/%v/%v/groups/", sentryConfig.Host, sentryConfig.OrganisationSlug, sentryConfig.ProjectSlug)
	_, _, errs := request.Get(url).SetBasicAuth(sentryConfig.ApiKey, "").EndStruct(&data)
	if errs != nil {
		return SentryList{}, errs
	}

	sentryList := SentryList{SentryItem: data}
	return sentryList, nil
}

func formatTitleString(item SentryItem) string {
	startStr := fmt.Sprintf("[%v] %v", item.Count, item.Title)

	return startStr
}

func formatMetaString(item SentryItem) string {
	f, _ := time.Parse("2006-01-02T15:04:05Z07:00", item.LastSeen)
	endStr := fmt.Sprintf("    %v %v", ressource.FormatDate(f), item.Culprit)

	return endStr
}

func UpdateRender(tab *ui.List) {
	now := time.Now()
	if now.Before(nextUpdate) {
		return
	}

	previousAlertCopy := previousAlert
	previousAlert = make(map[string]int, config.Config.Sentry.MaxItem)

	items := []string{}

	sentryData, errs := getStats()
	if errs != nil {
		tab.Items = ressource.ErrorToSliceString(errs)
		return
	}

	for i, item := range sentryData.SentryItem {
		var color string
		prevCount, isSet := previousAlertCopy[item.ShareId]
		if isSet {
			count, _ := strconv.Atoi(item.Count)
			if (count > prevCount) {
				color = "fg-black,bg-yellow"
			}
		} else {
			color = "bg-magenta"
		}

		countInt, _ := strconv.Atoi(item.Count)
		previousAlert[item.ShareId] = countInt

		var sentryTitle string
		var sentryInfo string
		if color != "" {
			sentryTitle = "[" + formatTitleString(item) + "](" + color + ")"
			sentryInfo = "[" + formatMetaString(item) + "](" + color + ")"
		} else {
			sentryTitle = formatTitleString(item)
			sentryInfo = formatMetaString(item)
		}
		items = append(items, sentryTitle)
		items = append(items, sentryInfo)

		if i >= config.Config.Sentry.MaxItem {
			break
		}
	}

	tab.Items = items
	nextUpdate = nextUpdate.Add(10 * time.Second)
}
