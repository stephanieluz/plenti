package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"plenti/generated"

	"github.com/spf13/cobra"
)

// siteCmd represents the site command
var siteCmd = &cobra.Command{
	Use:   "site [name]",
	Short: "Creates default folders and files for a new site",
	Long: `The project scaffolding follows this convention:
- plenti.json = sitewide configuration
- content/ = json files that hold site content
- content/pages/ = regular site pages in json format
- content/pages/_blueprint.json = template for the structure of a typical page
- content/pages/_index.json = the aggregate, or landing page
- content/pages/about.json = an example page
- content/pages/contact.json = another example page
- layout/ =  the html structure of the site
- layout/content/ = node level structure that has a route and correspond to content
- layout/components/ = smaller reusable structures that can be used within larger ones
- layout/global/ = base level html wrappers
- layout/static/ = holds assets like images or videos
- node_modules/ = frontend libraries managed by npm
- package.json = npm configuration file
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}
		if len(args) > 1 {
			return errors.New("names cannot have spaces")
		}
		if len(args) == 1 {
			return nil
		}
		return fmt.Errorf("invalid name specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Create base directory for site
		newpath := filepath.Join(".", args[0])
		os.MkdirAll(newpath, os.ModePerm)

		// Loop through generated file defaults to create site scaffolding
		for file, content := range generated.Defaults {
			// Create the directories needed for the current file
			os.MkdirAll(newpath+filepath.Dir(file), os.ModePerm)
			// Create the current default file
			err := ioutil.WriteFile(newpath+file, content, 0755)
			if err != nil {
				fmt.Printf("Unable to write file: %v", err)
			}
		}

		fmt.Printf("Created plenti site scaffolding in \"%v\" folder\n", newpath)

		fmt.Printf("Installing NPM dependencies...\n")
		command := exec.Command("npm", "install")
		command.Dir = newpath
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Run()
		fmt.Printf("NPM install complete!")

	},
}

func init() {
	newCmd.AddCommand(siteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// siteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// siteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
