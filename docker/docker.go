// Package docker is a simple wrapper around the Docker client binary for sending Aenthill events.
package docker

import (
	"os/exec"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/log"
)

// Docker is our working struct.
type Docker struct {
	ctx *context.Context
}

// New creates a Docker instance.
func New(ctx *context.Context) (*Docker, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil, errors.Wrap("docker", err)
	}
	return &Docker{ctx}, nil
}

// Run calls the run command from docker client binary.
func (d *Docker) Run(image, key, event, payload string) error {
	b := &builder{}
	b.run(image, event, payload)
	b.withEnv(context.KeyEnvVar, key)
	b.withEnv(context.ImageEnvVar, image)
	b.withEnv(context.FromContainerIDEnvVar, d.ctx.FromContainerID)
	b.withEnv(context.HostProjectDirEnvVar, d.ctx.HostProjectDir)
	b.withEnv(context.ContainerProjectDirEnvVar, d.ctx.ProjectDir)
	b.withEnv(context.LogLevelEnvVar, d.ctx.LogLevel)
	b.withMount("/var/run/docker.sock", "/var/run/docker.sock")
	b.withMount(d.ctx.HostProjectDir, d.ctx.ProjectDir)
	cmd := b.build()
	if d.ctx.IsContainer() {
		log.Infof(`"%s" (container ID = "%s") is sending event "%s" with payload "%s" to "%s" (manifest key = "%s")`, d.ctx.Image, d.ctx.Hostname, event, payload, image, key)
	} else {
		log.Infof(`sending event "%s" with payload "%s" to "%s" (manifest key = "%s")`, event, payload, image, key)
	}
	return errors.Wrapf("docker", cmd.Run(), "%s", cmd.Args)
}

// Reply calls the exec command from docker client binary.
func (d *Docker) Reply(event, payload string) error {
	b := &builder{}
	b.exec(d.ctx.FromContainerID, event, payload)
	cmd := b.build()
	log.Infof(`"%s" (container ID = "%s") is replying to "%s" with event "%s" and payload "%s"`, d.ctx.Image, d.ctx.Hostname, d.ctx.FromContainerID, event, payload)
	return errors.Wrapf("docker", cmd.Run(), "%s", cmd.Args)
}
