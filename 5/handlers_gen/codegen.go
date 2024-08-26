package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type CommentVars struct {
	Url    string
	Auth   bool
	Method string
}

type MethodData struct {
	FuncName     string
	Vars         CommentVars
	RecvName     string
	InStructName string
	Validation   []MethodDataValidation
}

type MethodDataValidation struct {
	Name               string
	Type               string
	ValidatorRequired  bool
	ValidatorParamName string
	ValidatorMin       int
	ValidatorMax       int
	ValidatorDefault   string
	ValidatorEnum      []string
}

type tpl struct {
	Reciver string
	Methods []MethodData
}

var reciverTpl = template.Must(template.New("reciverTpl").Parse(`
// {{.Reciver}}
func (h *{{.Reciver}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
	{{ range .Methods }}
	case r.URL.Path == "{{ .Vars.Url }}":
		{{ if ne .Vars.Method "ALL" }}
		if r.Method != "{{ .Vars.Method }}" {
			w.WriteHeader(http.StatusNotAcceptable)
			json, _ := json.Marshal(JsonResponse{Error: "bad method"})
			w.Write(json)
			return 
		}
		{{ end }}
        h.{{ .FuncName }}{{ .Vars.Method }}(w, r)
	{{ end }}
    default:
        w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(JsonResponse{Error: "unknown method"})
		w.Write(json)
    }
}

{{ range .Methods }}
func (h *{{$.Reciver}}) {{ .FuncName }}{{ .Vars.Method }}(w http.ResponseWriter, r *http.Request) {
	params := {{ .InStructName }}{}
	ctx := r.Context()
	{{ if eq .Vars.Auth true }}
	authToken := r.Header.Get("X-Auth")
	if authToken != "100500" {
		w.WriteHeader(http.StatusForbidden)
		json, _ := json.Marshal(JsonResponse{Error: "unauthorized"})
		w.Write(json)
		return
	}
	{{ end }}
	{{ range .Validation }} {{ if eq .Type "int" }}
	{{ .Name }}Value, err := strconv.Atoi(r.FormValue("{{ .ValidatorParamName }}"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} must be int"})
		w.Write(json)
		return
	}
	params.{{ .Name }} = {{ .Name }}Value
	{{ if eq .ValidatorRequired true }}
	if params.{{ .Name }} == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} must me not empty"})
		w.Write(json)
		return
	}
	{{ end}}
	{{ if gt (len .ValidatorDefault) 0 }}
	if len(params.{{ .Name }}) == 0 {
		params.{{ .Name }} = {{ .ValidatorDefault }}
	}
	{{ end }}
	{{ if gt .ValidatorMin -1 }}
	if params.{{ .Name }} < {{ .ValidatorMin }} {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} must be >= {{ .ValidatorMin }}"})
		w.Write(json)
		return
	}
	{{ end }}
	{{ if gt .ValidatorMax -1 }}
	if params.{{ .Name }} > {{ .ValidatorMax }} {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} must be <= {{ .ValidatorMax }}"})
		w.Write(json)
		return
	}
	{{ end }}{{ else if eq .Type "string" }}
	params.{{ .Name }} = r.FormValue("{{ .ValidatorParamName }}")
	{{ if eq .ValidatorRequired true }}
	if len(params.{{ .Name }}) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} must me not empty"})
		w.Write(json)
		return
	}
	{{end}}
	{{ if gt (len .ValidatorDefault) 0 }}
	if len(params.{{ .Name }}) == 0 {
		params.{{ .Name }} = "{{ .ValidatorDefault }}"
	}
	{{ end }}
	{{ if gt (len .ValidatorEnum) 0 }}
	{{$sep := ""}}{{ .Name }}Enums := []string{ {{ range .ValidatorEnum }}{{$sep}}"{{.}}"{{$sep = ", "}}{{ end }} }
	if !slices.Contains({{ .Name }}Enums, params.{{ .Name }}) {
		w.WriteHeader(http.StatusBadRequest)
		{{$sep := ""}}json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} must be one of [{{ range .ValidatorEnum }}{{$sep}}{{.}}{{$sep = ", "}}{{ end }}]"})
		w.Write(json)
		return
	}
	{{ end }}
	{{ if gt .ValidatorMin -1 }}
	if utf8.RuneCountInString(params.{{ .Name }}) < {{ .ValidatorMin }} {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "{{ .ValidatorParamName }} len must be >= {{ .ValidatorMin }}"})
		w.Write(json)
		return
	}
	{{ end }}
	{{ end }} {{ end }}
	res, err := h.{{ .FuncName }}(ctx, params)
	if err != nil {
		statuseCode := http.StatusInternalServerError
		if apiErr, ok := err.(ApiError); ok {
			statuseCode = apiErr.HTTPStatus
		}
		w.WriteHeader(statuseCode)
		json, _ := json.Marshal(JsonResponse{Error: err.Error()})
		w.Write(json)
		return
	}
	json, err := json.Marshal(JsonResponse{Response: res})
	if err != nil {
		panic(err)
	}
	w.Write(json)
}
{{ end }}
`))

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])
	defer out.Close()

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import (`)
	fmt.Fprintln(out, "\t\"strconv\"")
	//fmt.Fprintln(out, "\t\"fmt\"")
	fmt.Fprintln(out, "\t\"net/http\"")
	fmt.Fprintln(out, "\t\"encoding/json\"")
	fmt.Fprintln(out, "\t\"slices\"")
	fmt.Fprintln(out, "\t\"unicode/utf8\"")
	fmt.Fprintln(out, `)`)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "type JsonResponse struct {")
	fmt.Fprintln(out, "\tError string `json:\"error\"`")
	fmt.Fprintln(out, "\tResponse interface{} `json:\"response,omitempty\"`")
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out)

	recivers := make(map[string][]MethodData, 0)

DECLS_LOOP:
	for _, f := range node.Decls {

		g, ok := f.(*ast.FuncDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.FuncDecl\n", f)
			continue
		}

		if g.Doc == nil {
			fmt.Printf("SKIP %v, doesnt have comments\n", g.Name)
			continue
		}

		needCodegen := false
		vars := CommentVars{Url: "/", Auth: false, Method: "ALL"}
		for _, comment := range g.Doc.List {
			needCodegen = strings.HasPrefix(comment.Text, "// apigen:api")
			if needCodegen {
				_, jsonStr, found := strings.Cut(comment.Text, "apigen:api")
				if found {
					err := json.Unmarshal([]byte(jsonStr), &vars)
					if err != nil {
						fmt.Printf("SKIP %v, cant parse json vars\n", g.Name)
						continue DECLS_LOOP
					}
				}
				break
			}
		}

		if !needCodegen {
			fmt.Printf("SKIP %#v, doesnt have apigen:api mark\n", g.Name)
			continue
		}

		method := MethodData{
			FuncName: g.Name.String(),
			Vars:     vars,
		}

		recvStarExpr, ok := g.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		recvIdent, ok := recvStarExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		recvTypespec, ok := recvIdent.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		method.RecvName = recvTypespec.Name.String()

		inIdent, ok := g.Type.Params.List[1].Type.(*ast.Ident)
		if !ok {
			continue
		}

		inTypespec, ok := inIdent.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		inStruct, ok := inTypespec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		method.InStructName = inIdent.Name

		for _, inStructField := range inStruct.Fields.List {

			validatorField := MethodDataValidation{
				Name:               inStructField.Names[0].Name,
				ValidatorParamName: strings.ToLower(inStructField.Names[0].Name),
				ValidatorMin:       -1,
				ValidatorMax:       -1,
			}

			inFieldIdent, ok := inStructField.Type.(*ast.Ident)
			if !ok {
				continue
			}

			if inFieldIdent.Name != "string" && inFieldIdent.Name != "int" {
				continue
			}

			validatorField.Type = inFieldIdent.Name

			if inStructField.Tag != nil {
				tagStr := reflect.StructTag(inStructField.Tag.Value[1 : len(inStructField.Tag.Value)-1])
				tags := strings.Split(tagStr.Get("apivalidator"), ",")
				for _, tag := range tags {
					tagSplited := strings.SplitN(tag, "=", 2)
					paramName := tagSplited[0]
					paramValue := ""
					if len(tagSplited) == 2 {
						paramValue = tagSplited[1]
					}

					switch paramName {
					case "required":
						validatorField.ValidatorRequired = true
					case "min":
						value, err := strconv.Atoi(paramValue)
						if err != nil {
							continue
						}
						validatorField.ValidatorMin = value

					case "max":
						value, err := strconv.Atoi(paramValue)
						if err != nil {
							continue
						}
						validatorField.ValidatorMax = value

					case "paramname":
						if len(paramValue) == 0 {
							continue
						}
						validatorField.ValidatorParamName = paramValue

					case "default":
						validatorField.ValidatorDefault = paramValue

					case "enum":
						validatorField.ValidatorEnum = strings.Split(paramValue, "|")
					}
				}
			}

			method.Validation = append(method.Validation, validatorField)
		}

		if _, ok := recivers[method.RecvName]; !ok {
			recivers[method.RecvName] = make([]MethodData, 0)
		}

		recivers[method.RecvName] = append(recivers[method.RecvName], method)
	}

	for reciver, methods := range recivers {
		reciverTpl.Execute(out, tpl{Reciver: reciver, Methods: methods})
	}
}
