package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	configCmd.AddCommand(configSetCmd)
	addConfigFlags(configSetCmd.Flags())
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Updates the configuration",
	Long: `Updates the configuration. Set the flags for the options
you want to change. Other options will remain unchanged.`,
	Args: cobra.NoArgs,
	Run: python(func(cmd *cobra.Command, args []string, d pythonData) {
		flags := cmd.Flags()
		set, err := d.store.Settings.Get()
		checkErr(err)

		ser, err := d.store.Settings.GetServer()
		checkErr(err)

		hasAuth := false
		flags.Visit(func(flag *pflag.Flag) {
			switch flag.Name {
			case "baseurl":
				ser.BaseURL = mustGetString(flags, flag.Name)
			case "root":
				ser.Root = mustGetString(flags, flag.Name)
			case "socket":
				ser.Socket = mustGetString(flags, flag.Name)
			case "cert":
				ser.TLSCert = mustGetString(flags, flag.Name)
			case "key":
				ser.TLSKey = mustGetString(flags, flag.Name)
			case "address":
				ser.Address = mustGetString(flags, flag.Name)
			case "port":
				ser.Port = mustGetString(flags, flag.Name)
			case "log":
				ser.Log = mustGetString(flags, flag.Name)
			case "signup":
				set.Signup = mustGetBool(flags, flag.Name)
			case "auth.method":
				hasAuth = true
			case "shell":
				set.Shell = convertCmdStrToCmdArray(mustGetString(flags, flag.Name))
			case "branding.name":
				set.Branding.Name = mustGetString(flags, flag.Name)
			case "branding.disableExternal":
				set.Branding.DisableExternal = mustGetBool(flags, flag.Name)
			case "branding.files":
				set.Branding.Files = mustGetString(flags, flag.Name)
			case "onlyoffice":
				set.OnlyOffice = mustGetString(flags, flag.Name)
			}
		})

		getUserDefaults(flags, &set.Defaults, false)

		// read the defaults
		auther, err := d.store.Auth.Get(set.AuthMethod)
		checkErr(err)

		// check if there are new flags for existing auth method
		set.AuthMethod, auther = getAuthentication(flags, hasAuth, set, auther)

		err = d.store.Auth.Save(auther)
		checkErr(err)
		err = d.store.Settings.Save(set)
		checkErr(err)
		err = d.store.Settings.SaveServer(ser)
		checkErr(err)
		printSettings(ser, set, auther)
	}, pythonConfig{}),
}
