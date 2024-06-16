package http

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"html/template"
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

	data := struct {
		AssetsHost string
		Chart1     interface{}
		Chart2     interface{}
	}{
		AssetsHost: page.AssetsHost,
		Chart1:     chart1,
		Chart2:     chart2,
	}

	// Создание шаблона из текста и его исполнение
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const htmlTemplate = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <title>Charts Demo</title>
		    <!-- Подключение стилей go-echarts -->
		    <link rel="stylesheet" href="{{ .AssetsHost }}/echarts.min.css">
		    <style>
		        .chart-container {
		            width: 300px;
		            height: 200px;
		            margin: 2f0px;
		            float: left;
		        }
		    </style>
		</head>
		<body>
		    <!-- Блок для графика 1 -->
		    <div class="chart-container" height="100" id="chart1"></div>
		    
		    <!-- Блок для графика 2 -->
		    <div class="chart-container" width="300" id="chart2"></div>
		    
		    <!-- Подключение библиотеки go-echarts -->
		    <script src="{{ .AssetsHost }}/echarts.min.js"></script>
		    <script>
		        // Инициализация и настройка графиков
		        var chart1 = echarts.init(document.getElementById('chart1'));
		        var chart2 = echarts.init(document.getElementById('chart2'));
		        
		        // Данные для графиков
		        var option1 = {{ .Chart1 }};
		        var option2 = {{ .Chart2 }};
		        
		        // Установка опций и отображение графиков
		        chart1.setOption(option1);
		        chart2.setOption(option2);
		    </script>
		</body>
		</html>
		`
