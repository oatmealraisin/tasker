// Copyright © 2017 Ryan Murphy <murphy2902@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: Rework for SQLite
package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oatmealraisin/tasker/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	logsFile string
	workDir  string

	pgURL  string
	pgUser string
	pgPass string
	pgDir  string
	pgLogs string
)

// startCmd represents the init command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Initialize the a server daemon to interact with.",
	Long: `Sets up a Tasker Rest API with a Postgre database. Will read in config from a
user defined file (See man page) and restart the Postgre server from saved
data, if possible.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.RunE(cmd, args); err != nil {
			log.Fatal(err)
		}
	},
	RunE: start,
}

// TODO: Comment
func init() {
	startCmd.Flags().StringVarP(&logsFile, "logs", "l", "", "File to output logs to.")
	startCmd.Flags().StringVarP(&pgURL, "database-url", "s", "", "URL of an already existing PostgreSQL database.")
	startCmd.Flags().StringVarP(&pgUser, "database-user", "u", "", "Username to use with PostgreSQL database.")
	startCmd.Flags().StringVarP(&pgPass, "database-pass", "p", "", "Password to use with PostgreSQL database.")
	startCmd.Flags().StringVarP(&workDir, "data-dir", "d", "", "Working directory for tasker.")

	TaskerCmd.AddCommand(startCmd)
}

// Main function for the start command, starts up the Postgre server and Rest
// API
func start(cmd *cobra.Command, args []string) error {
	if err := setup(); err != nil {
		return err
	}

	store, _ := storage.NewSQLiteStorage()

	if err := store.Connect(); err != nil {
		return err
	}

	if err := server.Start(db); err != nil {
		return err
	}

	return nil
}

func connectPostgre() (*sql.DB, error) {
	if pgUser == "" {
		return nil, fmt.Errorf("When specifying a URL, must provide a Username")
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		pgUser, pgPass, pgUser, pgURL)
	return sql.Open("postgres", dbinfo)
}
