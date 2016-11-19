package web

import (
	"fmt"
	"html/template"
	"net/http"
)

var testtemplate = template.Must(template.New("tests").Parse(`
<html>
	<body>
		<ul>
		{{range .Results}}
			<li>
				<h3>{{.Package}}</h3>
				<p>{{.Started.Format "Jan _2 15:04:05"}} ({{.Ended.Sub .Started}})</p>
				<p>{{.Err}}</p>
				<pre style="font-size: 10; overflow-wrap: break-word; ">
					{{.Output}}
				</pre>
			</li>
		{{end}}
		</ul>
	</body>
</html>
`))

func (ww *web) status(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")

	err := ww.track.WriteStatus(w)

	if err != nil {
		fmt.Fprintf(w, "Failed to marshal. %v", err)
	}
}
