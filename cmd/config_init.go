package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/filebrowser/filebrowser/v2/settings"
)

func init() {
	configCmd.AddCommand(configInitCmd)
	addConfigFlags(configInitCmd.Flags())
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new database",
	Long: `Initialize a new database to use with File Browser. All of
this options can be changed in the future with the command
'filebrowser config set'. The user related flags apply
to the defaults when creating new users and you don't
override the options.`,
	Args: cobra.NoArgs,
	Run: python(func(cmd *cobra.Command, args []string, d pythonData) {
		defaults := settings.UserDefaults{}
		flags := cmd.Flags()
		getUserDefaults(flags, &defaults, true)
		authMethod, auther := getAuthentication(flags)

		s := &settings.Settings{
			Key:        generateKey(),
			Signup:     mustGetBool(flags, "signup"),
			Shell:      convertCmdStrToCmdArray(mustGetString(flags, "shell")),
			AuthMethod: authMethod,
			Defaults:   defaults,
			Branding: settings.Branding{
				Name:            mustGetString(flags, "branding.name"),
				DisableExternal: mustGetBool(flags, "branding.disableExternal"),
				Files:           mustGetString(flags, "branding.files"),
			},
			OnlyOffice: "",
		}

		ser := &settings.Server{
			Address: mustGetString(flags, "address"),
			Socket:  mustGetString(flags, "socket"),
			Root:    mustGetString(flags, "root"),
			BaseURL: mustGetString(flags, "baseurl"),
			TLSKey:  mustGetString(flags, "key"),
			TLSCert: mustGetString(flags, "cert"),
			Port:    mustGetString(flags, "port"),
			Log:     mustGetString(flags, "log"),
		}

		err := d.store.Settings.Save(s)
		checkErr(err)
		err = d.store.Settings.SaveServer(ser)
		checkErr(err)
		err = d.store.Auth.Save(auther)
		checkErr(err)

		fmt.Printf(`
Congratulations! You've set up your database to use with File Browser.
Now add your first user via 'filebrowser users add' and then you just
need to call the main command to boot up the server.
`)
		printSettings(ser, s, auther)
	}, pythonConfig{noDB: true}),
}
