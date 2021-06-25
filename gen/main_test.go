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
	tmpl := template.Must(template.ParseFS(embeddedTemplates, "tmpl/v1.tmpl"))

	_ = tmpl.Execute(buf, map[string]interface{}{
		"dbHelperPkg": dbHelperPkg,
		"ecmLogPkg":   ecmLogPkg,
		"outPkg":      "testpkg",
		"table":       `"t_role"`,
		"db":          "dao.EcmAdminDbClient.DB",
		"fullName":    "admindao.RoleInfo",
		"structName":  "RoleInfo",
		"ormStruct":   "modelRole",
		"ormName":     "SRole",
		"converters":  map[string]string{"Name": "string"},
		"groupers":    map[string]string{"Id": "int64"},
		"extFile":     filepath.Base(extFilePath),
	})
	fmt.Println(buf)
}
