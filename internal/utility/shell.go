package utility

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type ShellResult struct {
	Success  bool
	ExitCode int
	Args     []string
	Stdout   string
}

func (m *ShellResult) AsError() error {
	return fmt.Errorf(
		"command[%s] failed with exitcode=[%d]",
		strings.Join(m.Args, " "),
		m.ExitCode,
	)
}

func Shell(args ...string) *ShellResult {
	logger := slog.Default()
	if len(args) < 1 {
		panic("Method called in an invalid manner, expected at least one argument")
	}

	first, rest := args[0], args[1:]

	logger.Debug(
		"Running shell command",
		"command", first,
		"args", strings.Join(rest, " "),
	)

	cmd := exec.Command(
		first,
		rest...,
	)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	stdoutBuf := bytes.NewBufferString("")
	stderrBuf := bytes.NewBufferString("")

	stderrDone := logStreamAsync(stderr, logger, "stderr", stderrBuf)
	stdoutDone := logStreamAsync(stdout, logger, "stdout", stdoutBuf)

	if err := cmd.Start(); err != nil {
		logger.Error(err.Error())
		return &ShellResult{
			Success:  false,
			ExitCode: -1,
			Args:     args,
		}
	}

	cmd.Wait()

	// wait for logging the output to complete
	_, _ = <-stderrDone, <-stdoutDone

	exitCode := cmd.ProcessState.ExitCode()
	return &ShellResult{
		Success:  exitCode == 0,
		ExitCode: exitCode,
		Args:     args,
		Stdout:   stdoutBuf.String(),
	}
}

func logStreamAsync(readStream io.Reader, logger *slog.Logger, name string, outBuf *bytes.Buffer) chan int {
	c := make(chan int)

	go func() {
		scanner := bufio.NewScanner(readStream)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			text := scanner.Text()
			fmt.Fprintln(outBuf, text)
			logger.Debug(fmt.Sprintf(
				"[%s] %s",
				name,
				text,
			))
		}

		// tell caller we're done, may take a bit to actually finish logging the stream
		// after the stream is technically closed
		c <- 0
	}()

	return c
}
