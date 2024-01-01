package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type cmds []*cobra.Command

func main() {
	c := readCfg()
	fmt.Printf("%+v\n", c)
	for _, name := range []string{"toot", "poot"} {
		rootCmd.AddCommand(&cobra.Command{
			Use: name,
			Run: printName,
		})
	}
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func readCfg() cmds {
	d, err := os.ReadFile("testdata/cmds.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var dc []map[string]any
	err = yaml.Unmarshal(d, &dc)
	if err != nil {
		log.Fatal(err)
	}
	var c cmds
	for _, cc := range dc {
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, cc)
		if err != nil {
			log.Fatal(err)
		}
		println(buf.String())
	}
	return c
}

func printName(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var tmpl = template.Must(template.New("name").Parse(tmplStr))

const tmplStr = `
&cobra.Command{
{{range $key, $val := . }}
	{{$key}}: {{$val}},
{{end -}}
}
`
