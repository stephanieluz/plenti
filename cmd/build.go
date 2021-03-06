package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plenti/cmd/build"
	"plenti/readers"
	"time"

	"github.com/spf13/cobra"
)

// BuildDirFlag allows users to override name of default build directory (public)
var BuildDirFlag string

func setBuildDir(siteConfig readers.SiteConfig) string {
	var buildDir string
	// Check if directory is overridden by flag.
	if BuildDirFlag != "" {
		// If dir flag exists, use it.
		buildDir = BuildDirFlag
	} else {
		buildDir = siteConfig.BuildDir
	}
	return buildDir
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Creates the static assets for your site",
	Long: `Build generates the actual HTML, JS, and CSS into a directory
of your choosing. The files that are created are all
you need to deploy for your website.`,
	Run: func(cmd *cobra.Command, args []string) {
		Build()
	},
}

// Build creates the compiled app that gets deployed.
func Build() {

	buildStart := time.Now()

	// Get settings from config file.
	siteConfig := readers.GetSiteConfig()

	// Check flags and config for directory to build to.
	buildDir := setBuildDir(siteConfig)

	// Get the full path for the build directory of the site.
	buildPath := filepath.Join(".", buildDir)

	// Clear out any previous build dir of the same name.
	if _, buildPathExistsErr := os.Stat(buildPath); buildPathExistsErr == nil {
		deleteBuildErr := os.RemoveAll(buildPath)
		fmt.Printf("\nRemoving old \"%v\" build directory\n", buildPath)
		if deleteBuildErr != nil {
			fmt.Println(deleteBuildErr)
			return
		}
	}
	// TODO: Should this automatically happen after stopping local server?

	// Create the buildPath directory.
	err := os.MkdirAll(buildPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Unable to create \"%v\" build directory: %s\n", buildDir, err)
	} else {
		fmt.Printf("\nCreating \"%v\" build directory\n", buildDir)
	}

	tempFiles := build.EjectTemp()

	build.EjectCopy(buildPath)

	start := time.Now()
	// Build JSON from "content/" directory.
	staticBuildStr, allNodesStr := build.DataSource(buildPath, siteConfig)
	elapsed := time.Since(start)
	fmt.Printf("Creating data_source took %s\n", elapsed)

	start = time.Now()
	// Prep the client SPA.
	clientBuildStr := build.Client(buildPath)
	elapsed = time.Since(start)
	fmt.Printf("Prepping client SPA data took %s\n", elapsed)

	start = time.Now()
	svelteBuild := exec.Command("node", "ejected/build.js", clientBuildStr, staticBuildStr, allNodesStr)
	svelteBuild.Stdout = os.Stdout
	svelteBuild.Stderr = os.Stderr
	svelteBuild.Run()
	elapsed = time.Since(start)
	fmt.Printf("\nCompiling components and creating static HTML took %s\n", elapsed)

	// Run Gopack (custom Snowpack alternative).
	build.Gopack(buildPath)

	build.EjectClean(tempFiles)

	elapsed = time.Since(buildStart)
	fmt.Printf("\nTotal build took %s\n", elapsed)

}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	buildCmd.Flags().StringVarP(&BuildDirFlag, "dir", "d", "", "change name of the build directory")
}
