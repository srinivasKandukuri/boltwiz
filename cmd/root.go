package cmd

import (
	"fmt"
	"github.com/pkg/browser"
	"go.uber.org/zap"
	"time"

	"github.com/boltdbgui/common/logger"
	"github.com/boltdbgui/modules/database/repository"
	"github.com/boltdbgui/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Boltdb Server",
	Long:  `Start the boltdb browser server`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger("DEBUG")
		err := repository.Init(input.dbPath)
		if err != nil {
			panic(fmt.Sprintf("Unable to initialize db connection in dbpath %s : error %v", input.dbPath, err))
		}
		defer repository.Close() // nolint: errcheck

		if input.local {
			go func() {
				time.Sleep(2 * time.Second) // wait for server to start
				log.Info("Opening browser....")
				err := browser.OpenURL("http://localhost:8090")

				if err != nil {
					log.Error("Error while opening browser", zap.Error(err))
				}
			}()
		}

		server.StartServer()
	},
}

var input = new(struct {
	dbPath string
	local  bool
})

func init() {
	rootCmd.Flags().StringVarP(&input.dbPath, "db-path", "d", "", "path to the bolt db file")
	rootCmd.Flags().BoolVarP(&input.local, "local", "l", false, "open the browser automatically")

	err := rootCmd.MarkFlagRequired("db-path")
	if err != nil {
		panic(errors.Wrap(err, "error while setting required flag db-path"))
	}
}

func Execute() error {

	return rootCmd.Execute()
}
