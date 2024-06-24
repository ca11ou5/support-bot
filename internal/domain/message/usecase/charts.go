package usecase

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var wcData = map[string]interface{}{
	"Sam S Club":               10000,
	"Macys":                    6181,
	"Amy Schumer":              4386,
	"Jurassic World":           4055,
	"Charter Communications":   2467,
	"Chick Fil A":              2244,
	"Planet Fitness":           1898,
	"Pitch Perfect":            1484,
	"Express":                  1689,
	"Home":                     1112,
	"Johnny Depp":              985,
	"Lena Dunham":              847,
	"Lewis Hamilton":           582,
	"KXAN":                     555,
	"Mary Ellen Mark":          550,
	"Farrah Abraham":           462,
	"Rita Ora":                 366,
	"Serena Williams":          282,
	"NCAA baseball tournament": 273,
	"Point Break":              265,
}

func generateWCData(data map[string]interface{}) (items []opts.WordCloudData) {
	items = make([]opts.WordCloudData, 0)
	for k, v := range data {
		items = append(items, opts.WordCloudData{Name: k, Value: v})
	}
	return
}

func (uc *UseCase) GetCharts() (*charts.Line, *charts.Line, *charts.WordCloud) {
	stats := uc.messageRepo.GetStats()

	latestMessages := charts.NewLine()
	allMessages := charts.NewLine()

	wc := charts.NewWordCloud()
	wc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "basic WordCloud example",
		}))

	wc.AddSeries("wordcloud", generateWCData(wcData)).
		SetSeriesOptions(
			charts.WithWorldCloudChartOpts(
				opts.WordCloudChart{
					SizeRange: []float32{14, 80},
				}),
		)

	//wordCloud.Validate()
	//wordCloud.SetSeriesOptions()

	var xAllMsg []string
	var yAllMsg []opts.LineData

	var xLatMsg []string
	var yLatMsg []opts.LineData

	for _, v := range stats {
		xAllMsg = append(xAllMsg, v.Timestamp.Format("15:04"))
		yAllMsg = append(yAllMsg, opts.LineData{
			Value: v.AllMessagesCount,
		})

		xLatMsg = append(xLatMsg, v.Timestamp.Format("15:04"))
		yLatMsg = append(yLatMsg, opts.LineData{
			Value: v.LatestMessagesCount,
		})
	}

	allMessages.SetXAxis(xAllMsg).AddSeries("ЛиНиЯ", yAllMsg).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
	latestMessages.SetXAxis(xLatMsg).AddSeries("fgds", yLatMsg)

	allMessages.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", 400),
			Height: fmt.Sprintf("%dpx", 200),
		}),
	)

	return allMessages, latestMessages, wc
}
