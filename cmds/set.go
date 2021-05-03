package cmds

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var setCmds = &cobra.Command{
	Use:     "set [options]",
	Short:   "Configures CLI settings",
	Long:    "Configures text colors and notification properties.\nValid colors are black, blue, cyan, green, magenta, red, white, yellow, default, or none.\nValid audios are default, fancy, or none.\nValid durations are short or long.\n\"cali set default\" sets everything to their default.",
	Example: "cali set --subtitle=magenta --success=blue --error=yellow --audio=fancy --duration=long",
	Run: func(cmd *cobra.Command, args []string) {
		configs := db.GetConfigurations()

		if cmd.Flags().NFlag() == 0 && len(args) == 0 {
			if len(configs) == 0 {
				cmd.Flags().VisitAll(func(flag *pflag.Flag) {
					configs[flag.Name] = flag.DefValue
				})
			}
			printSubtitle("\nCurrent Settings\n")
			for k, v := range configs {
				if k == "audio" || k == "duration" || k == "subtitle" || k == "success" || k == "error" {
					fmt.Printf("%s: %s\n", k, v)
				}
			}
			return
		}

		setDefault := len(args) == 1 && args[0] == "default"
		if len(args) > 0 && !setDefault {
			printError("Oops, invalid argument. Make sure to use flags to set configurations")
			return
		}

		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			// fmt.Printf("%+v\n", flag)
			if flag.Changed || setDefault {
				var err error
				var dataType string
				switch flag.Name {
				case "audio":
					err = validAudio(flag.Value.String())
					dataType = "audio"
				case "duration":
					err = validDuration(flag.Value.String())
					dataType = "duration"
				case "subtitle":
					_, err = strToColor(flag.Value.String(), color.New(color.FgCyan))
					dataType = "color"
				case "success":
					_, err = strToColor(flag.Value.String(), color.New(color.FgGreen))
					dataType = "color"
				case "error":
					_, err = strToColor(flag.Value.String(), color.New(color.FgRed))
					dataType = "color"
				}
				if err != nil {
					printError(fmt.Sprintf("Oops, '%s' is not a valid %s", flag.Value, dataType))
				} else {
					printSuccess(fmt.Sprintf("Set %s to %s", flag.Name, flag.Value))
					configs[flag.Name] = flag.Value.String()
				}
			}

			if _, ok := configs[flag.Name]; !ok {
				configs[flag.Name] = flag.DefValue
			}
		})

		err := db.SetConfigurations(configs)
		if err != nil {
			printError("Oops, unable to set configurations")
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(setCmds)
	setCmds.Flags().String("subtitle", "default", "Set subtitle text color")
	setCmds.Flags().String("error", "default", "Set error text color")
	setCmds.Flags().String("success", "default", "Set success text color")
	setCmds.Flags().String("audio", "default", "Set notification sound")
	setCmds.Flags().String("duration", "short", "Set notification duration")
}

func strToColor(colorStr string, defaultColor *color.Color) (*color.Color, error) {
	switch colorStr {
	case "default":
		return defaultColor, nil
	case "none":
		return color.New(), nil
	case "red":
		return color.New(color.FgRed), nil
	case "green":
		return color.New(color.FgGreen), nil
	case "blue":
		return color.New(color.FgBlue), nil
	case "black":
		return color.New(color.FgBlack), nil
	case "cyan":
		return color.New(color.FgCyan), nil
	case "magenta":
		return color.New(color.FgMagenta), nil
	case "white":
		return color.New(color.FgWhite), nil
	case "yellow":
		return color.New(color.FgYellow), nil
	default:
		return defaultColor, errors.New("Oops, invalid color")
	}
}

func validAudio(audioStr string) error {
	switch audioStr {
	case "default":
		return nil
	case "none":
		return nil
	case "fancy":
		return nil
	default:
		return errors.New("Oops, invalid audio")
	}
}

func validDuration(durationStr string) error {
	switch durationStr {
	case "short":
		return nil
	case "long":
		return nil
	default:
		return errors.New("Oops, invalid duration")
	}
}
