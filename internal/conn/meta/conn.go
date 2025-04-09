package meta

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/risingwavelabs/wavekit/internal/utils"
	"golang.org/x/mod/semver"
)

type RisectlConnection struct {
	risectlPath string
	endpoint    string
	version     string
}

func (c *RisectlConnection) prepareCmd(ctx context.Context, args ...string) *exec.Cmd {

	cmd := exec.CommandContext(ctx, c.risectlPath, args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("RW_META_ADDR=%s", c.endpoint))
	return cmd
}

func (c *RisectlConnection) RunCombined(ctx context.Context, args ...string) (string, int, error) {
	log.Default().Printf("Running risectl (meta addr: %s) command: %s %v", c.endpoint, c.risectlPath, args)
	cmd := c.prepareCmd(ctx, args...)
	out, err := cmd.CombinedOutput()
	exitCode := -2
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	log.Default().Printf("risectl command output: %s, exit code: %d, error: %v", string(out), exitCode, err)
	return string(out), exitCode, err
}

func (c *RisectlConnection) Run(ctx context.Context, args ...string) (string, string, int, error) {
	log.Default().Printf("Running risectl (meta addr: %s) command: %s %v", c.endpoint, c.risectlPath, args)
	cmd := c.prepareCmd(ctx, args...)
	stdoutBuf := bytes.NewBuffer(nil)
	stderrBuf := bytes.NewBuffer(nil)
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	err := cmd.Run()

	exitCode := -2
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	log.Default().Printf("risectl command stdout: %s, stderr: %s, exit code: %d, error: %v", stdout, stderr, exitCode, err)
	return stdout, stderr, exitCode, err
}

// sample: backup job succeeded: job 1,
var regexExtractJobID = regexp.MustCompile(`backup job succeeded: job (\d+)`)

func (c *RisectlConnection) MetaBackup(ctx context.Context) (int64, error) {
	res, ec, err := c.RunCombined(ctx, "meta", "backup-meta")
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

func (c *RisectlConnection) DeleteSnapshot(ctx context.Context, snapshotID int64) error {
	res, ec, err := c.RunCombined(ctx,
		utils.IfElse(
			semver.Compare(c.version, "v2.0.1") >= 0,
			[]string{"meta", "delete-meta-snapshots", "--snapshot-ids", strconv.FormatInt(snapshotID, 10)},
			[]string{"meta", "delete-meta-snapshots", strconv.FormatInt(snapshotID, 10)},
		)...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w, output: %s, exit code: %d", err, res, ec)
	}
	return nil
}
