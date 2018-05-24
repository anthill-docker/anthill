// Package docker helps sending Aenthill events with the docker client binary.
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/apex/log"
	isatty "github.com/mattn/go-isatty"
)

const (
	// HostProjectDirEnvVariable is the name of the environment variable which contains the host project directory.
	HostProjectDirEnvVariable = "AENTHILL_HOST_PROJECT_DIR"
	// SenderEnvVariable is the name of the environment variable which contains the image name of the sender.
	SenderEnvVariable = "AENTHILL_SENDER"
	// LogLevelEnvVariable is the name of the environment variable which contains the log level.
	LogLevelEnvVariable = "AENTHILL_LOG_LEVEL"
	// WhoAmIEnvVariable is the name of the environment variable which contains the image name.
	WhoAmIEnvVariable = "AENTHILL_WHO_AM_I"
	// InsideContainerProjectDir is the location the the container where the host project directory is mounted.
	InsideContainerProjectDir = "/aenthill"
	// DefaultBinary is the default binary to call in an image.
	DefaultBinary = "aent"
)

// EventContext gathers all required data of an event.
type EventContext struct {
	// WhoAmI is the image/binary which is sending the event.
	WhoAmI string
	// Image is the image which receives the event.
	Image string
	// Binary is the command which handles the event in the targeted image.
	Binary string
	// HostProjectDir is the project directory on the host.
	HostProjectDir string
	// LogLevel is the log level which should be used by the targeted image.
	// Accepted values for log level: DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
	LogLevel string
}

/*
Send use the docker client binary to send an event.

It will in fact run a command in the targeted image, using the following template:

 docker run [-ti] --rm -v "/var/run/docker.sock:/var/run/docker.sock" -v "HostProjectDir:/aenthill" -e "AENTHILL_SENDER=WhoAmI" -e "AENTHILL_HOST_PROJECT_DIR=HostProjectDir"  -e "AENTHILL_LOG_LEVEL=LogLevel" Image Binary even payload

Important: it relies on COMSPEC environment variable on Windows and SHELL
on posix system to know which interpreter to use for calling the docker client
binary.
*/
func Send(event string, payload string, context *EventContext) error {
	log.WithFields(log.Fields{
		"from":     context.WhoAmI,
		"to":       context.Image,
		"event":    event,
		"paypload": payload,
	}).Info("sending event")

	var dockerOpts []string

	// attaches Stdin if TTY is a terminal.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		dockerOpts = append(dockerOpts, "-ti")
	}

	dockerOpts = append(dockerOpts, "--rm")
	dockerOpts = append(dockerOpts, fmt.Sprintf("-v \"%s:%s\"", "/var/run/docker.sock", "/var/run/docker.sock"))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-v \"%s:%s\"", context.HostProjectDir, InsideContainerProjectDir))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", SenderEnvVariable, context.WhoAmI))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", HostProjectDirEnvVariable, context.HostProjectDir))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", LogLevelEnvVariable, context.LogLevel))

	var args []string
	args = append(args, []string{"docker", "run"}...)
	args = append(args, dockerOpts...)
	args = append(args, []string{context.Image, context.Binary, event, payload}...)

	var e *exec.Cmd
	if runtime.GOOS == "windows" {
		e = exec.Command(os.Getenv("COMSPEC"), "/c", strings.Join(args, " "))
	} else {
		e = exec.Command(os.Getenv("SHELL"), "-c", strings.Join(args, " "))
	}

	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	e.Stdin = os.Stdin

	log.WithField("from", context.WhoAmI).Debugf("executing command %s", e.Args)
	return e.Run()
}