package risectl

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

type RisectlConn interface {
	Run(ctx context.Context, args ...string) (string, int, error)
	MetaBackup(ctx context.Context) (int64, error)
}

type RisectlConnection struct {
	risectlPath string
	endpoint    string
	version     string
}

func (c *RisectlConnection) Run(ctx context.Context, args ...string) (string, int, error) {
	log.Default().Printf("Running risectl (meta addr: %s) command: %s %v", c.endpoint, c.risectlPath, args)
	cmd := exec.CommandContext(ctx, c.risectlPath, args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("RW_META_ADDR=%s", c.endpoint))
	out, err := cmd.CombinedOutput()
	exitCode := -2
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	log.Default().Printf("risectl command output: %s, exit code: %d, error: %v", string(out), exitCode, err)
	return string(out), exitCode, err
}

// sample: backup job succeeded: job 1,
var regexExtractJobID = regexp.MustCompile(`backup job succeeded: job (\d+)`)

func (c *RisectlConnection) MetaBackup(ctx context.Context) (int64, error) {
	res, ec, err := c.Run(ctx, "meta", "backup-meta")
	if err != nil {
		return 0, fmt.Errorf("failed to backup meta: %w, output: %s, exit code: %d", err, res, ec)
	}

	matches := regexExtractJobID.FindStringSubmatch(res)
	if len(matches) != 2 {
		return 0, fmt.Errorf("failed to extract job ID from output: %s, exit code: %d", res, ec)
	}

	jobID := matches[1]
	return strconv.ParseInt(jobID, 10, 64)
}
