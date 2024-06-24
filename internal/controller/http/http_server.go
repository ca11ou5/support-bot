package http

import (
	"bytes"
	"fmt"
	chartrender "github.com/go-echarts/go-echarts/v2/render"
	"log"

	"html/template"
	"io"
	"net/http"
)

func (s *Server) RegisterHTTPServer() error {
	http.HandleFunc("/", s.statsHandler)
	http.HandleFunc("/stats", s.index)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) index(w http.ResponseWriter, _ *http.Request) {
	chart1, chart2, chart3, chart4 := s.useCase.GetCharts()

	//f1, _ := os.Create("./images/chart1.txt")
	//defer f1.Close()
	//f2, _ := os.Create("./images/chart2.txt")
	//defer f2.Close()
	//f3, _ := os.Create("./images/chart3.txt")
	//defer f3.Close()

	// Сохранение графика как изображения
	//chart1.Renderer = newSnippetRenderer(chart1, chart1.Validate)
	//chart2.Renderer = newSnippetRenderer(chart2, chart3.Validate)
	//chart3.Renderer = newSnippetRenderer(chart3, chart3.Validate)
	//
	//chart1.Render(f1)
	//chart2.Render(f2)
	//chart3.Render(f3)

	var htmlSnippet1 = renderToHtml(chart1)
	var htmlSnippet2 = renderToHtml(chart2)
	var htmlSnippet3 = renderToHtml(chart3)
	var htmlSnippet4 = renderToHtml(chart4)

	templatePath := "./web/2.html"

	// Чтение содержимого шаблона из файла
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}

	data := []template.HTML{htmlSnippet1, htmlSnippet2, htmlSnippet3, htmlSnippet4}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./web/1.html"

	// Загрузка файла и отправка его клиенту
	http.ServeFile(w, r, filePath)
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

func renderToHtml(c interface{}) template.HTML {
	var buf bytes.Buffer
	r := c.(chartrender.Renderer)
	err := r.Render(&buf)
	if err != nil {
		log.Printf("Failed to render chart: %s", err)
		return ""
	}

	return template.HTML(buf.String())
}
