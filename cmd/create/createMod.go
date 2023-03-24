package create

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"supreme-flamego/internal/mod/example"
	"supreme-flamego/pkg/colorful"
	"supreme-flamego/pkg/fs"
)

var (
	appName  string
	dir      string
	force    bool
	StartCmd = &cobra.Command{
		Use:     "create",
		Short:   "create a new mod",
		Example: "mod create -n users",
		Run: func(cmd *cobra.Command, args []string) {
			err := load()
			if err != nil {
				fmt.Println(colorful.Red(err.Error()))
				os.Exit(1)
			}
			fmt.Println(colorful.Green("Module " + appName + " generate success under " + dir))
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "create a new mod with provided name")
	StartCmd.PersistentFlags().StringVarP(&dir, "path", "p", "internal/mod", "new file will generate under provided path")
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the mod")
}

func load() error {
	if appName == "" {
		return errors.New("mod name should not be empty, use -n")
	}

	dirEntries, err := example.FS.ReadDir(".")
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			err = fs.IsNotExistMkDir(path.Join(dir, appName, dirEntry.Name()))
			if err != nil {
				return err
			}
			continue
		}
	}

	if !force {
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				continue
			}
			// check if the file is existed
			p := path.Join(dir, appName, dirEntry.Name())
			if fs.FileExist(p) {
				return errors.New("file " + p + " is existed, use -f to force generate")
			}
		}
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		if dirEntry.Name() == "embed.go" {
			continue
		}
		file, err := example.FS.ReadFile(dirEntry.Name())
		if err != nil {
			return err
		}
		file = bytes.ReplaceAll(file, []byte("example"), []byte(appName))
		fs.FileCreate(*bytes.NewBuffer(file), path.Join(dir, appName, dirEntry.Name()))
	}

	return nil
}
