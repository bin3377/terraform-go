package terraform

import (
	"encoding/json"
	"strings"
)

type ApplyConfig struct {
	ChDir      string
	PlanBinary []byte
	PlanPath   string
}

type ApplyResult struct {
	State State
	Logs  []string
}

func (c *ApplyConfig) getArgs() []string {
	args := []string{"apply"}

	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}

	args = append(args, "-json", "-auto-approve", c.PlanPath)

	return args
}

func (c *ApplyConfig) getShowArgs() []string {
	args := []string{"show", "-json", c.PlanPath}

	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}

	return args
}

func (c *ApplyConfig) Apply() (*ApplyResult, error) {
	args := c.getArgs()
	out, err := runTFWithEnv(args)
	if err != nil {
		return nil, err
	}

	r := ApplyResult{
		Logs: strings.Split(out, "\n"),
	}

	out, err = runTFWithEnv(c.getShowArgs())
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal([]byte(out), &state); err != nil {
		return nil, err
	}
	r.State = state
	return &r, nil
}
