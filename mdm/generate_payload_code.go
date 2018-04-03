// +build ignore

package main

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/as/micromdm/mdm"
	"golang.org/x/tools/imports"
)

const stub = `
// generated by go:generate, DO NOT EDIT

package mdm

func NewPayload(request *CommandRequest) (*Payload, error) {
	requestType := request.RequestType
	payload := newPayload(requestType)
	switch requestType {
	{{range .ManyFields}}
	case "{{.Name}}":
		payload.Command.{{.Name}} = request.{{.Name}}
	{{ end }}
	case "ProfileList",
		"ProvisioningProfileList",
		"CertificateList",
		"SecurityInfo",
		"StopMirroring",
		"ClearRestrictionsPassword",
		"UserList",
		"LogOutUser",
		"DisableLostMode",
		"DeviceLocation",
		"ManagedMediaList",
		"OSUpdateStatus",
		"DeviceConfigured",
		"AvailableOSUpdates",
		"Restrictions",
		"ShutDownDevice",
		"RestartDevice":
		return payload, nil
	default:
		return nil, fmt.Errorf("Unsupported MDM RequestType %v", requestType)
	}
	return payload, nil
}
`

func main() {
	out := flag.String("out", "new_payload.go", "path to output file")
	flag.Parse()
	type param struct {
		Name string
	}
	type params struct {
		ManyFields []param
		Simple     []param
	}

	var p params
	val := reflect.ValueOf(mdm.Command{})
	for i := 1; i < val.NumField(); i++ {
		name := val.Field(i).Type().Name()
		if val.Field(i).NumField() == 0 {
			p.Simple = append(p.Simple, param{Name: name})
		} else {
			p.ManyFields = append(p.ManyFields, param{Name: name})
		}
	}
	var tmpl = template.Must(template.New("test").Parse(stub))
	var buf bytes.Buffer
	tmpl.Execute(&buf, p)
	pretty, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(*out, pretty, 0644)

}
