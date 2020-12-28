package ethereum

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	prysmLogger       = "prysm"
	prysmStdErrLogger = "prysm err"
)

// logPipe prints out logs from prysm. We don't end when context
// is canceled beacause there are often logs printed after this.
func logPipe(pipe io.ReadCloser, identifier string) error {
	reader := bufio.NewReader(pipe)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println("closing", identifier, err)
			return err
		}

		message := strings.ReplaceAll(str, "\n", "")
		log.Println(identifier, message)
	}
}

// StartPrysm starts a prysm daemon in another goroutine
// and logs the results to the console.
func StartPrysm(ctx context.Context, arguments string, g *errgroup.Group) error {
	parsedArgs := strings.Split(arguments, " ")
	cmd := exec.Command(
		"/app/beacon-chain",
		parsedArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	g.Go(func() error {
		return logPipe(stdout, prysmLogger)
	})

	g.Go(func() error {
		return logPipe(stderr, prysmStdErrLogger)
	})

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%w: unable to start prysm", err)
	}

	g.Go(func() error {
		<-ctx.Done()

		log.Println("sending interrupt to prysm")
		return cmd.Process.Signal(os.Interrupt)
	})

	return cmd.Wait()
}
