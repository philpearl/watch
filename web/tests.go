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
		<table>
		{{range .Results}}
			<tr>
				<td>{{.Package}}</td>
				<td>{{.Err}}</td>
				<td>{{.Output}}</td>
				<td>{{.Started}}</td>
				<td>{{.Ended}}</td>
			</tr>
		{{end}}
		</table>
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
