package terraform

import (
	"encoding/json"
	"fmt"
	"os"
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

func (c *ApplyConfig) getArgs() ([]string, error) {
	if len(c.PlanBinary) == 0 && c.PlanPath == "" {
		return nil, fmt.Errorf("plan binary or path must be provided")
	}

	if len(c.PlanBinary) > 0 && c.PlanPath != "" {
		return nil, fmt.Errorf("plan binary and path cannot both be provided")
	}

	planPath := c.PlanPath
	if len(c.PlanBinary) > 0 {
		planPath = tempDir + "/plan.tfplan"
		f, err := os.Create(planPath)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(c.PlanBinary); err != nil {
			return nil, err
		}
	}

	args := []string{"apply"}

	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}

	args = append(args, "-no-color", "-auto-approve", planPath)

	return args, nil
}

func (c *ApplyConfig) getShowArgs() []string {
	args := []string{"show", "-json", c.PlanPath}

	if c.ChDir != "" {
		args = append([]string{"-chdir=" + c.ChDir}, args...)
	}

	return args
}

func (c *ApplyConfig) Apply() (*ApplyResult, error) {
	args, err := c.getArgs()
	if err != nil {
		return nil, err
	}
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
