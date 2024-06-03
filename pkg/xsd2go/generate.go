package xsd2go

import (
	"fmt"

	"github.com/harald-mue/xsd2go/pkg/template"
	"github.com/harald-mue/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string, xmlnsOverrides []string) error {
	fmt.Printf("Processing '%s'\n", xsdPath)
	ws, err := xsd.NewWorkspace(fmt.Sprintf("%s/%s", goModule, outputDir), xsdPath, xmlnsOverrides)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}
		if err := template.GenerateTypes(sch, outputDir); err != nil {
			return err
		}
	}

	return nil
}
