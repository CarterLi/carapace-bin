package cmd

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fastfetch",
	Short: "A neofetch-like tool for fetching system information and displaying them in a pretty way",
	Long:  "https://github.com/fastfetch-cli/fastfetch",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() error {
	return rootCmd.Execute()
}

type FlagData struct {
	Short  string
	Long   string
	Desc   string
	Remark string
	Pseudo bool
	Arg    *struct {
		Type     string
		Optional bool
		Default  interface{}
		Enum     map[string]string
	}
}

func init() {
	carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
		carapace.Gen(rootCmd).Standalone()

		var groups map[string][]FlagData
		if out, err := exec.Command("fastfetch", "--help-raw").Output(); err == nil {
			json.Unmarshal(out, &groups)
			if len(groups) == 0 {
				return
			}
		} else {
			return
		}

		actionBools := carapace.ActionValues("true", "false")
		actionColors := carapace.ActionValues("black", "red", "green", "yellow", "blue", "magenta", "cyan", "white", "default")
		actionMap := carapace.ActionMap{}

		for _, flags := range groups {
			for _, flag := range flags {
				if flag.Pseudo {
					continue
				}

				if flag.Arg != nil {
					switch flag.Arg.Type {
					case "bool":
						rootCmd.Flags().BoolP(flag.Long, flag.Short, false, flag.Desc)
						actionMap[flag.Long] = actionBools

					case "color":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						actionMap[flag.Long] = actionColors

					case "command":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						actionMap[flag.Long] = carapace.ActionExecCommand("fastfetch", "--list-modules", "autocompletion")(func(output []byte) carapace.Action {
							options := []string{"color", "format"}
							for _, flags := range groups {
								for _, flag := range flags {
									options = append(options, flag.Long)
								}
							}
							for _, line := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
								x, _, _ := strings.Cut(line, ":")
								options = append(options, x)
							}

							return carapace.ActionValues(options...)
						})

					case "config":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						actionMap[flag.Long] = carapace.ActionExecCommand("fastfetch", "--list-presets", "autocompletion")(func(output []byte) carapace.Action {
							presets := strings.Split(strings.TrimRight(string(output), "\n"), "\n")
							return carapace.ActionValues(presets...)
						})

					case "enum":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						var valuesDescribed []string
						for key, desc := range flag.Arg.Enum {
							valuesDescribed = append(valuesDescribed, key, desc)
						}
						actionMap[flag.Long] = carapace.ActionValuesDescribed(valuesDescribed...)

					case "logo":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						actionMap[flag.Long] = carapace.ActionExecCommand("fastfetch", "--list-logos", "autocompletion")(func(output []byte) carapace.Action {
							res := []string{"none", "Disable logo", "small", "Show small logo if supported"}
							for _, logo := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
								res = append(res, logo, "Builtin logo")
							}
							return carapace.ActionValuesDescribed(res...)
						})

					case "num":
						rootCmd.Flags().IntP(flag.Long, flag.Short, 0, flag.Desc)

					case "path":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						actionMap[flag.Long] = carapace.ActionFiles()

					case "structure":
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
						actionMap[flag.Long] = carapace.ActionExecCommand("fastfetch", "--list-modules", "autocompletion")(func(output []byte) carapace.Action {
							var texts []string
							for _, line := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
								name, desc, _ := strings.Cut(line, ":")
								texts = append(texts, name, desc)
							}
							return carapace.ActionValuesDescribed(texts...).UniqueList(":")
						})

					default:
						rootCmd.Flags().StringP(flag.Long, flag.Short, "", flag.Desc)
					}
				} else {
					rootCmd.Flags().BoolP(flag.Long, flag.Short, false, flag.Desc)
				}
			}
		}

		carapace.Gen(rootCmd).FlagCompletion(actionMap)

		carapace.Gen(rootCmd).PositionalCompletion(
			carapace.ActionDirectories(),
		)
	})
}
