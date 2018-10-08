// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type listDataTmpl struct {
	List map[int]string
}

var listTmpl = template.Must(template.New("list.html").Parse(`<html>
  <head>
    <style>
/* make pre wrap */
pre {
 white-space: pre-wrap;       /* css-3 */
 white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
 white-space: -pre-wrap;      /* Opera 4-6 */
 white-space: -o-pre-wrap;    /* Opera 7 */
 word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
    </style>
  </head>
  <body>
	{{ range $i, $file := .List }}
    <p><li><a href="{{ $file }}">{{ $file}}</a></li></p>
	{{ end }}
  </body>
</html>
`))

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: pass it to direct/redirect download
	workDir, _ := os.Getwd()
	fi, err := os.Stat(workDir + r.URL.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	if r.URL.Path != "/download/" {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			redirectDownload(w, r, workDir+r.URL.Path+"/redirect")
			fmt.Println("request res " + r.URL.Path + " is a directory")
			break
		case mode.IsRegular():
			directDownload(w, r, workDir+r.URL.Path)
			fmt.Println("request res " + r.URL.Path + " is a file")
			break
		default:
			fmt.Println("unknown")
		}
	} else {
		fmt.Println("request page")
		pageHandler(w, r)
	}
}

func directDownload(w http.ResponseWriter, r *http.Request, path string) {
	//TODO: set up file server for static file
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// Don't use http.ServeFile due to directory traversal attack.
	http.ServeFile(w, r, path)
}

func redirectDownload(w http.ResponseWriter, r *http.Request, path string) {
	//TODO: redirect to the link in redirect file
	redirectUrl, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, string(redirectUrl), http.StatusSeeOther)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	files, err := ioutil.ReadDir(workDir + "/download/")
	if err != nil {
		fmt.Println(err)
		return
	}
	lists := make(map[int]string)
	i := int(0)
	for _, file := range files {
		s := file.Name()
		if file.Mode()&os.ModeSymlink != 0 && s != "..data" {
			lists[i] = s
			i++
		}
	}
	renderTemplate(w, listTmpl, listDataTmpl{
		List: lists,
	})

}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	err := tmpl.Execute(w, data)
	if err == nil {
		return
	}
	fmt.Println(err)

	switch err := err.(type) {
	case *template.Error:
		fmt.Println("Error rendering template %s: %s", tmpl.Name(), err)

		http.Error(w, "Internal server error. Template rendering failed", http.StatusInternalServerError)
		break
	default:
	}

}
func main() {
	port := flag.String("port", "80", "Port for the application")
	flag.Parse()

	http.HandleFunc("/download/", downloadHandler)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		panic(err)
	}
}
