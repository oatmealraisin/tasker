package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/oatmealraisin/tasker/pkg/tasker"

	"github.com/spf13/cobra"
)

const (
	// TODO: Finish these
	cliExplain = `Server for task management`
	cliLong    = `TODO`
	cliShort   = `TODO`
	cliUse     = `plant`
)

var (
	logsFile   string
	configFile string

	plantDataDir string

	postgreURL       string
	postgreUser      string
	postgrePassword  string
	postgreDirectory string
	postgreLogsFile  string
)

func main() {
	command := &cobra.Command{
		Use:   cliUse,
		Short: cliShort,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.RunE(cmd, args); err != nil {
				log.Fatal(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(); err != nil {
				fmt.Println(err.Error())
			}
			return nil
		},
	}

	command.Flags().StringVarP(&logsFile, "logs", "l", "", "File to output logs to.")
	command.Flags().StringVarP(&postgreURL, "database-url", "s", "", "URL of an already existing PostgreSQL database.")
	command.Flags().StringVarP(&postgreUser, "database-user", "u", "", "Username to use with PostgreSQL database.")
	command.Flags().StringVarP(&postgrePassword, "database-pass", "p", "", "Password to use with PostgreSQL database.")
	command.Flags().StringVarP(&plantDataDir, "data-dir", "d", "", "Working directory for plant.")
	command.Flags().StringVarP(&configFile, "config", "c", "", "Config directory to use. Defaults to '$XDG_CONFIG_HOME/plant/config'.")

	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var err error

	err = setup()
	if err != nil {
		return err
	}

	postgreLogsFile = ""
	if postgreURL == "" {
		postgres, err := startPostgre()
		if err != nil {
			return err
		}
		defer postgres.Process.Kill()
	}

	db, err := connectPostgre()
	if err != nil {
		return err
	} else {
		defer db.Close()
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err = tasker.Start(db); err != nil {
		return err
	}

	return nil
}

func setup() error {
	if plantDataDir == "" {
		if os.Getenv("XDG_DATA_HOME") != "" {
			plantDataDir = filepath.Join(os.Getenv("XDG_DATA_HOME"), "plant")
		} else if os.Getenv("HOME") != "" {
			plantDataDir = filepath.Join(os.Getenv("HOME"), ".local/plant")
		} else {
			plantDataDir = filepath.Join("/etc", "plant")
		}
	}

	return nil
}

func startPostgre() (*exec.Cmd, error) {
	if postgreUser == "" {
		postgreUser = "plant"
	}

	postgreDirectory = filepath.Join(plantDataDir, "postgres")

	if _, err := os.Stat(postgreDirectory); err != nil && os.IsNotExist(err) {
		fmt.Println("Creating data dir:", postgreDirectory)
		if err := os.MkdirAll(postgreDirectory, os.ModePerm); err != nil {
			return nil, fmt.Errorf("Could not create postgre data directory, %s", err.Error())
		}
	}

	fmt.Println("Initializing database..")
	initdb, err := exec.LookPath("initdb")
	if err != nil {
		return nil, err
	}

	output, err := exec.Command(initdb,
		"-D", filepath.Join(postgreDirectory, "database"),
		"-U", "plant",
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("initdb: %s\n%s", err.Error(), output)
	}

	fmt.Println("Starting PostgreSQL..")
	postgre, err := exec.LookPath("postgres")
	if err != nil {
		return nil, err
	}

	postgres := exec.Command(postgre,
		//"--single",
		//"-c max_connections=1",
		"-c", "listen_addresses=",
		//fmt.Sprintf("-r %s", postgreLogsFile),
		//fmt.Sprintf("-p %s", "6666"),
		"-D", filepath.Join(postgreDirectory, "database"),
		"-c", fmt.Sprintf("unix_socket_directories=%s", postgreDirectory),
	)

	err = postgres.Start()

	if err != nil {
		return nil, fmt.Errorf("postgres: %s", err.Error())
	}

	time.Sleep(3 * time.Second)

	fmt.Println("Creating database..")
	createdb, err := exec.LookPath("createdb")
	if err != nil {
		return nil, err
	}

	output, err = exec.Command(createdb, "plant",
		"--host", postgreDirectory,
		"-U", postgreUser,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("createdb: %s\n%s", err.Error(), output)
	}

	postgreURL = postgreDirectory

	return postgres, nil
}

func stopPostgre() error {
	return nil
}

func connectPostgre() (*sql.DB, error) {
	if postgreUser == "" {
		return nil, fmt.Errorf("When specifying a URL, must provide a Username")
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		postgreUser, postgrePassword, postgreUser, postgreURL)
	return sql.Open("postgres", dbinfo)
}
