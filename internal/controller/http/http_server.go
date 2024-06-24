package http

import (
	"bufio"
	"bytes"
	"fmt"
	chartrender "github.com/go-echarts/go-echarts/v2/render"

	"html/template"
	"io"
	"net/http"
	"os"
)

func (s *Server) RegisterHTTPServer() error {
	http.HandleFunc("/", s.chart)
	http.HandleFunc("/stats", s.statsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) chart(w http.ResponseWriter, _ *http.Request) {
	chart1, chart2, chart3 := s.useCase.GetCharts()

	f1, _ := os.Create("./images/chart1.txt")
	defer f1.Close()
	f2, _ := os.Create("./images/chart2.txt")
	defer f2.Close()
	f3, _ := os.Create("./images/chart3.txt")
	defer f3.Close()

	// Сохранение графика как изображения
	chart1.Renderer = newSnippetRenderer(chart1, chart1.Validate)
	chart2.Renderer = newSnippetRenderer(chart2, chart3.Validate)
	chart3.Renderer = newSnippetRenderer(chart3, chart3.Validate)

	chart1.Render(f1)
	chart1.Render(f2)
	chart1.Render(f3)

	var buffer bytes.Buffer
	file, _ := os.Open("./images/chart.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n") // Добавляем перевод строки, чтобы сохранить оригинальный формат файла
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println(buffer.String())

	templatePath := "./web/2.html"

	// Чтение содержимого шаблона из файла
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}

	data := []string{buffer.String()}

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}

	//page := components.NewPage()
	//
	//page.AddCharts(chart1, chart2, chart3)
	//
	//page.Render(w)
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	// Определяем путь к файлу index.html
	htmlFilePath := "./web/2.html"

	// Открываем файл
	file, err := os.Open(htmlFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось открыть файл HTML: %s", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Устанавливаем правильный Content-Type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Копируем содержимое файла в ответ
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Не удалось прочитать файл HTML: %s", err), http.StatusInternalServerError)
		return
	}
}

var baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}
</script>
`

type snippetRenderer struct {
	c      interface{}
	before []func()
}

func newSnippetRenderer(c interface{}, before ...func()) chartrender.Renderer {
	return &snippetRenderer{c: c, before: before}
}

func (r *snippetRenderer) Render(w io.Writer) error {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	tpl := template.
		Must(template.New(tplName).
			Funcs(template.FuncMap{
				"safeJS": func(s interface{}) template.JS {
					return template.JS(fmt.Sprint(s))
				},
			}).
			Parse(baseTpl),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}
