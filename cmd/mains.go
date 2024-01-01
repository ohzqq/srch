package main

import (
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type cmds []*cobra.Command

func main() {
	//c := readCfg()
	//fmt.Printf("%+v\n", c)
	err := cmd.rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func readCfg() cmds {
	d, err := os.ReadFile("../testdata/cmds.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var dc Cfg
	err = yaml.Unmarshal(d, &dc)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("cmd/commands.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = tmpl.ExecuteTemplate(f, "cobra", dc)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range dc.Commands {
		f, err := os.Create("cmd/" + c.Name + ".go")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = tmpl.ExecuteTemplate(f, "run", c.Name)
		if err != nil {
			log.Fatal(err)
		}
	}

	var c cmds
	return c
}

type Cfg struct {
	Viper    bool  `yaml:"Viper"`
	Commands []Cmd `yaml:"Commands"`
}

type Cmd struct {
	Use        string            `yaml:"Use"`
	Name       string            `yaml:"Name"`
	Short      string            `yaml:"Short"`
	Long       string            `yaml:"Long"`
	Aliases    string            `yaml:"Aliases"`
	Run        string            `yaml:"Run"`
	FlagStruct map[string]string `yaml:"FlagStruct"`
	Flags      []Flag            `yaml:"Flags"`
	Parent     string            `yaml:"Parent"`
}

type Flag struct {
	Name       string `yaml:"Name"`
	Shorthand  string `yaml:"Shorthand"`
	Usage      string `yaml:"Usage"`
	Type       string `yaml:"Type"`
	Var        string `yaml:"Var"`
	Value      string `yaml:"Value"`
	Persistent bool   `yaml:"Persistent"`
	Viper      bool   `yaml:"Viper"`
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

var tmpl = template.Must(template.New("cobra").ParseFiles("../testdata/cobra.tmpl"))

const tmplStr = `
package main

import (
	"os"

	"github.com/spf13/cobra"
{{- if .Viper }}
	"github.com/spf13/viper"{{ end }}
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

{{range .Commands}}
	{{- $cmd := .Name}}

var {{$cmd}} = &cobra.Command{
	{{with .Use}}Use: "{{.}}",{{end}}
	{{with .Aliases}}Aliases: {{.}},{{end}}
	{{with .Short}}Short: "{{.}}",{{end}}
	{{with .Long}}Long: "{{.}}",{{end}}
	{{with .Run}}Run: {{.}},{{end}}
}

func init() {
	{{- range .Children}}
		{{$cmd}}.AddCommand({{.}})
	{{- end}}

	{{range .Flags}}
	{{$cmd}}.
		{{- with .Persistent}}Persistent{{end -}}
	Flags().{{.Type}}(
		"{{.Name}}",
		{{with .Shorthand}}"{{.}}",{{end}}
		{{.Value}},
		"{{.Usage}}",
	)
		{{if .Viper}}
	viper.BindPFlag("{{.Name}}", {{$cmd}}.Flags().Lookup("{{.Name}}"))
		{{end}}
		{{end}}
}
		{{end}}
`
