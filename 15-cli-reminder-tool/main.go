package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "GOLANG_CLI_REMINDER"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage:%s <hh:mm> <text_message>\n", os.Args[0])
		os.Exit(1)
	}

	now := time.Now()
	when := when.New(nil)
	when.Add(en.All...)
	when.Add(common.All...)

	t, err := when.Parse(os.Args[1], now)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if t == nil {
		fmt.Println("Unable to parse time")
		os.Exit(2)
	}
	if now.After(t.Time) {
		fmt.Println("Exit time is in the future")
		os.Exit(3)
	}

	diff := t.Time.Sub(now)
	if os.Getenv(markName) == markValue {
		time.Sleep(diff)
		err := beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	} else {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s\n", markName, markValue))
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
		fmt.Printf("Reminder will be displayed after: %v\n", diff.Round(time.Second))
	}
}
