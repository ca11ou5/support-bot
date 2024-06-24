package usecase

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateWCData(data map[string]interface{}) (items []opts.WordCloudData) {
	items = make([]opts.WordCloudData, 0)
	for k, v := range data {
		items = append(items, opts.WordCloudData{Name: k, Value: v})
	}
	return
}

func (uc *UseCase) GetCharts() (*charts.Line, *charts.Line, *charts.WordCloud, *charts.Pie) {
	stats := uc.messageRepo.GetStats()

	latestMessages := charts.NewLine()
	allMessages := charts.NewLine()

	wc := charts.NewWordCloud()
	circle := charts.NewPie()
	circle.AddSeries("Национальность обращающихся пользователей", []opts.PieData{{
		Name:  "RU",
		Value: 100,
	}})

	words := uc.messageRepo.GetWords()
	wc.AddSeries("Самые популярные слова", generateWCData(words)).
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

	allMessages.SetXAxis(xAllMsg).AddSeries("Общее количество входящих запросов", yAllMsg).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
	latestMessages.SetXAxis(xLatMsg).AddSeries("Количество последних входящих запросов", yLatMsg)

	allMessages.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", 550),
			Height: fmt.Sprintf("%dpx", 275),
		}),
	)
	latestMessages.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", 550),
			Height: fmt.Sprintf("%dpx", 275),
		}))
	wc.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", 550),
			Height: fmt.Sprintf("%dpx", 275),
		}))
	circle.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprintf("%dpx", 550),
			Height: fmt.Sprintf("%dpx", 275),
		}))

	return allMessages, latestMessages, wc, circle
}
