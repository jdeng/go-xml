package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"aqwari.net/xml/xsdgen"
)

func main() {
	log.SetFlags(0)
	var cfg xsdgen.Config
	cfg.NamespaceMap = map[string]string{
		"http://www.mismo.org/residential/2009/schemas": "",
		"http://www.w3.org/1999/xlink":                  "xlink",
		"http://www.datamodelextension.org/Schema/ULAD": "ULAD",
		"http://www.datamodelextension.org/Schema/DU":   "DU",
		"http://www.datamodelextension.org/Schema/LPA":  "LPA",
		"http://www.datamodelextension.org/Schema/ULDD": "ULDD",
		"http://www.transunion.com/namespace": "tuxml",
	}

	cfg.Option(xsdgen.DefaultOptions...)
	cfg.Option(xsdgen.LogOutput(log.New(os.Stderr, "", 0)))

	if err := cfg.GenCLI(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}

	for _, t := range cfg.SkipTypes {
		fmt.Printf("type %s = DEFAULT_EXTENSION\n", strings.ReplaceAll(t, "_", ""))
	}

	for _, t := range cfg.MixedTypes {
		switch t.BaseType {
		case "decimal":
			t.BaseType = "float64"
		case "gYear", "gMonth", "gDay", "dateTime", "time":
			t.BaseType = "time.Time"
		case "anyURI", "anySimpleType":
			t.BaseType = "string"
		case "boolean":
			t.BaseType = "bool"
		}

		t.Name = strings.ReplaceAll(t.Name, "_", "")
		t.Type = strings.ReplaceAll(t.Type, "_", "")
		tt := t.Name
		if t.Name == "Value" {
			tt = t.BaseType
		}
		fmt.Printf(`
func X%s(v %s) %s {
	return %s{
		%s: %s(v),
	}
}
func (x *%s) V() %s {
	return %s(x.%s)
}
		`, t.Type, t.BaseType, t.Type, t.Type, t.Name, tt, t.Type, t.BaseType, t.BaseType, t.Name)
	}

}
