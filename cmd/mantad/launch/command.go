package launch

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	launcher := &Launcher{}

	cmd := &cobra.Command{
		Use:          "mantad",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// global id generator need a real random number to init
			rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return launcher.run()
		},
	}

	viper.SetEnvPrefix("MANTA")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configPath := viper.GetString("CONFIG_PATH"); configPath != "" {
		switch path.Ext(configPath) {
		case ".json", ".yml", ".yaml":
			viper.SetConfigFile(configPath)
		case "":
			viper.AddConfigPath(configPath)
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	bindOptions(cmd, launcher.Options())

	return cmd
}

type Option struct {
	DestP interface{} // pointer to the destination

	EnvVar     string
	Flag       string
	Hidden     bool
	Persistent bool
	Required   bool
	Short      rune // using rune b/c it guarantees correctness. a short must always be a string of length 1

	Default interface{}
	Desc    string
}

func bindOptions(cmd *cobra.Command, opts []Option) {
	for _, o := range opts {
		flagset := cmd.Flags()
		if o.Persistent {
			flagset = cmd.PersistentFlags()
		}

		if o.Required {
			err := cmd.MarkFlagRequired(o.Flag)
			if err != nil {
				panic(err)
			}
		}

		envVar := o.Flag
		if o.EnvVar != "" {
			envVar = o.EnvVar
		}

		hasShort := o.Short != 0

		switch destP := o.DestP.(type) {
		case *string:
			var d string
			if o.Default != nil {
				d = o.Default.(string)
			}
			if hasShort {
				flagset.StringVarP(destP, o.Flag, string(o.Short), d, o.Desc)
			} else {
				flagset.StringVar(destP, o.Flag, d, o.Desc)
			}
			mustBindPFlag(o.Flag, flagset)
			*destP = viper.GetString(envVar)
		case *int:
			var d int
			if o.Default != nil {
				d = o.Default.(int)
			}
			if hasShort {
				flagset.IntVarP(destP, o.Flag, string(o.Short), d, o.Desc)
			} else {
				flagset.IntVar(destP, o.Flag, d, o.Desc)
			}
			mustBindPFlag(o.Flag, flagset)
			*destP = viper.GetInt(envVar)
		case *bool:
			var d bool
			if o.Default != nil {
				d = o.Default.(bool)
			}
			if hasShort {
				flagset.BoolVarP(destP, o.Flag, string(o.Short), d, o.Desc)
			} else {
				flagset.BoolVar(destP, o.Flag, d, o.Desc)
			}
			mustBindPFlag(o.Flag, flagset)
			*destP = viper.GetBool(envVar)
		case *time.Duration:
			var d time.Duration
			if o.Default != nil {
				d = o.Default.(time.Duration)
			}
			if hasShort {
				flagset.DurationVarP(destP, o.Flag, string(o.Short), d, o.Desc)
			} else {
				flagset.DurationVar(destP, o.Flag, d, o.Desc)
			}
			mustBindPFlag(o.Flag, flagset)
			*destP = viper.GetDuration(envVar)
		case *[]string:
			var d []string
			if o.Default != nil {
				d = o.Default.([]string)
			}
			if hasShort {
				flagset.StringSliceVarP(destP, o.Flag, string(o.Short), d, o.Desc)
			} else {
				flagset.StringSliceVar(destP, o.Flag, d, o.Desc)
			}
			mustBindPFlag(o.Flag, flagset)
			*destP = viper.GetStringSlice(envVar)
		case *map[string]string:
			var d map[string]string
			if o.Default != nil {
				d = o.Default.(map[string]string)
			}
			if hasShort {
				flagset.StringToStringVarP(destP, o.Flag, string(o.Short), d, o.Desc)
			} else {
				flagset.StringToStringVar(destP, o.Flag, d, o.Desc)
			}
			mustBindPFlag(o.Flag, flagset)
			*destP = viper.GetStringMapString(envVar)
		case pflag.Value:
			if hasShort {
				flagset.VarP(destP, o.Flag, string(o.Short), o.Desc)
			} else {
				flagset.Var(destP, o.Flag, o.Desc)
			}
			if o.Default != nil {
				_ = destP.Set(o.Default.(string))
			}
			mustBindPFlag(o.Flag, flagset)
			_ = destP.Set(viper.GetString(envVar))
		default:
			// if you get a panic here, sorry about that!
			// anyway, go ahead and make a PR and add another type.
			panic(fmt.Errorf("unknown destination type %t", o.DestP))
		}

		// so weirdness with the flagset her, the flag must be set before marking it
		// hidden. This is in contrast to the MarkRequired, which can be set before...
		if o.Hidden {
			if err := flagset.MarkHidden(o.Flag); err != nil {
				panic(err)
			}
		}
	}
}

func mustBindPFlag(key string, flagset *pflag.FlagSet) {
	if err := viper.BindPFlag(key, flagset.Lookup(key)); err != nil {
		panic(err)
	}
}
