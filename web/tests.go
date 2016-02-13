package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/philpearl/rebuilder"
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

func (ww *web) tests(w http.ResponseWriter, r *http.Request) {
	failingTests := ww.testrunner.GetResults()

	err := testtemplate.Execute(w, struct {
		Results []*rebuilder.TestResult
	}{
		Results: failingTests,
	})

	if err != nil {
		fmt.Fprintf(w, "Failed to execute template. %v", err)
	}
}
