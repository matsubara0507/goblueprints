package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func makeCmdChain(filename string) []*exec.Cmd {
	return []*exec.Cmd{
		exec.Command("synonyms"),
		exec.Command("sprinkle", "-config="+filename),
		exec.Command("coolify"),
		exec.Command("domainify"),
		exec.Command("available"),
	}
}

func main() {
	var filename = flag.String("config", "config.yaml", "変換方法")
	flag.Parse()

	cmdChain := makeCmdChain(*filename)
	cmdChain[0].Stdin = os.Stdin
	cmdChain[len(cmdChain)-1].Stdout = os.Stdout

	for i := 0; i < len(cmdChain)-1; i++ {
		thisCmd := cmdChain[i]
		nextCmd := cmdChain[i+1]
		stdout, err := thisCmd.StdoutPipe()
		if err != nil {
			log.Panicln(err)
		}
		nextCmd.Stdin = stdout
	}

	for _, cmd := range cmdChain {
		err := cmd.Start()
		if err != nil {
			log.Panicln(err)
		} else {
			defer cmd.Process.Kill()
		}
	}

	for _, cmd := range cmdChain {
		err := cmd.Wait()
		if err != nil {
			log.Panicln(err)
		}
	}
}
