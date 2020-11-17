package builder

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuildConfig struct {
	// The command to execute.
	Command []string `hcl:"command,optional"`
}

type Builder struct {
	config BuildConfig
}

// Implement Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Implement ConfigurableNotify
func (b *Builder) ConfigSet(config interface{}) error {
	c, ok := config.(*BuildConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *BuildConfig as parameter")
	}

	// validate the config
	if len(c.Command) == 0 {
		return fmt.Errorf("A command must be provided")
	}

	return nil
}

// Implement Builder
func (b *Builder) BuildFunc() interface{} {
	// return a function which will be called by Waypoint
	return b.build
}

// A BuildFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - *datadir.Project
// - *datadir.App
// - *datadir.Component
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet
//
// The output parameters for BuildFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (b *Builder) build(ctx context.Context, ui terminal.UI, src *component.Source) (*Binary, error) {

	// Code Sample From platform 'exec' plugin
	// Getting Command From Config
	args := b.config.Command

	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()

	// If we have a step set, abort it on exit
	var s terminal.Step
	defer func() {
		if s != nil {
			s.Abort()
		}
	}()

	// Render templates if set
	s = sg.Add("Executing command: %s", strings.Join(args, " "))

	// Ensure we're executing a binary
	if !filepath.IsAbs(args[0]) {
		log.Debug("command is not absolute, will look up on PATH", "command", args[0])
		path, err := exec.LookPath(args[0])
		if err != nil {
			log.Info("failed to find command on PATH", "command", args[0])
			return nil, err
		}

		log.Info("command is not absolute, replaced with value on PATH",
			"old_command", args[0],
			"new_command", path,
		)
		args[0] = path
	}

	// Run our command
	var cmd exec.Cmd
	cmd.Path = args[0]
	cmd.Args = args
	// cmd.Dir = src.Path
	cmd.Stdout = s.TermOutput()
	cmd.Stderr = cmd.Stdout

	// Run it
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	s.Done()

	return &Binary{}, nil
}
