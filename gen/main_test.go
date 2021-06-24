package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"
	"text/template"
)

func TestGen(t *testing.T) {
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.ParseFS(embeddedTemplates, "**/*.tmpl"))

	_ = tmpl.Execute(buf, map[string]interface{}{
		"dbHelperPkg": dbHelperPkg,
		"ecmLogPkg":   ecmLogPkg,
		"structName":  "Device",
		"ormName":     "SDevice",
		"converters":  map[string]string{"InstanceID": "string"},
		"groupers":    map[string]string{"AppID": "int64"},
		"outPkg":      "testpkg",
		"extFile":     filepath.Base(extFilePath),
		"table":       `"testTable"`,
		"db":          *dbVar,
		"fullName":    "dao.Device",
		"ormStruct":   "MDevice",
	})
	fmt.Println(buf)
}
