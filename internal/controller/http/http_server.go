package http

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"net/http"
)

func (s *Server) RegisterHTTPServer() error {
	http.HandleFunc("/", s.chart)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) chart(w http.ResponseWriter, _ *http.Request) {
	chart1, chart2 := s.useCase.GetCharts()

	page := components.NewPage()
	page.AssetsHost = "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/"

	page.AddCharts(chart1)
	page.AddCharts(chart2)

	page.Render(w)

	//chartsData := []entity.ChartData{
	//	{ID: "chart1", Option: chart},
	//	{ID: "chart2", Option: chart},
	//}
	//
	//// Создание экземпляра PageData
	//pageData := entity.PageData{
	//	Charts: chartsData,
	//}
	//
	//tmpl, err := template.New("charts").Parse(htmlTemplate)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//if err := tmpl.Execute(w, pageData); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
}

const htmlTemplate = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Multiple Charts</title>
        <!-- Подключение стилей go-echarts -->
        <link href="https://cdn.jsdelivr.net/npm/go-echarts@5.1.1/dist/bundle.css" rel="stylesheet">
        <style>
            .chart-container {
                width: 600px;
                height: 400px;
                margin: 20px;
                float: left;
            }
        </style>
    </head>
    <body>
        {{range .Charts}}
        <div class="chart-container">
            <div id="{{.ID}}" style="width: 100%; height: 100%;"></div>
        </div>
        {{end}}

        <!-- Подключение библиотеки go-echarts -->
        <script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/"></script>
        <script>
            {{range .Charts}}
            var {{.ID}} = echarts.init(document.getElementById('{{.ID}}'));
            {{.Option.JSCode}}
            {{end}}
        </script>
    </body>
    </html>
    `
