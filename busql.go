package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"
)

// type option struct {
// 	user string
// 	pwd  string
// }
var containerID string

func main() {
	// backup sooon db
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "backup",
				Aliases: []string{"b"},
				Usage:   string([]byte("mysql container ID")),
			},
		},
		Action: func(c *cli.Context) error {
			if c.IsSet("backup") {
				containerID = c.String("backup")
				DoBackup(scannerFileName, scannerPwd)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type getFileNameUI func() string
type getPwdUI func() string

func scannerFileName() string {

	var outputName string

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Input output file name(press return key for the default name):")
	scanner.Scan()
	outputName = scanner.Text()

	if scanner.Text() == "" {
		outputName = time.Now().Format("2006-01-02.sql")
	}

	return outputName + ".sql"
}
func scannerPwd() string {
	var pwd string

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Input mysql password:")
	scanner.Scan()
	pwd = scanner.Text()
	return "--password=" + pwd
}

func DoBackup(uiname getFileNameUI, uipwd getPwdUI) {
	filename := uiname()
	cmd := exec.Command("kubectl", "exec", containerID, "--", "/usr/bin/mysqldump", "-B", "-u", "root", uipwd(), "sooon_db")
	cmdOut, _ := cmd.StdoutPipe()
	cmdErr, _ := cmd.StderrPipe()
	cmd.Start()

	output, _ := ioutil.ReadAll(cmdOut)
	err, _ := ioutil.ReadAll(cmdErr)
	cmd.Wait()
	fmt.Println("mysqldump output is : " + string(output))
	if string(err) != "" {
		fmt.Println("mysqldump error is: " + string(err))
		os.Exit(4)
	}
	backfile, _ := os.Create(filename)
	defer backfile.Close()
	backfile.WriteString(string(output))
}
