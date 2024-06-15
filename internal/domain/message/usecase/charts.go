package usecase

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (uc *UseCase) GetCharts() (*charts.Line, *charts.Line) {
	stats := uc.messageRepo.GetStats()

	latestMessages := charts.NewLine()
	allMessages := charts.NewLine()

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

	return allMessages, latestMessages
}
