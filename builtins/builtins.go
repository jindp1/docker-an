package builtins

import (
	"runtime"

	"example.com/m/v2/api"
	apiserver "example.com/m/v2/api/server"
	"example.com/m/v2/daemon/networkdriver/bridge"
	"example.com/m/v2/dockerversion"
	"example.com/m/v2/engine"
	"example.com/m/v2/events"
	"example.com/m/v2/pkg/parsers/kernel"
	"example.com/m/v2/registry"
)

func Register(eng *engine.Engine) error {
	if err := daemon(eng); err != nil {
		return err
	}
	if err := remote(eng); err != nil {
		return err
	}
	if err := events.New().Install(eng); err != nil {
		return err
	}
	if err := eng.Register("version", dockerVersion); err != nil {
		return err
	}
	return registry.NewService().Install(eng)
}

// remote: a RESTful api for cross-docker communication
func remote(eng *engine.Engine) error {
	if err := eng.Register("serveapi", apiserver.ServeApi); err != nil {
		return err
	}
	return eng.Register("acceptconnections", apiserver.AcceptConnections)
}

// daemon: a default execution and storage backend for Docker on Linux,
// with the following underlying components:
//
// * Pluggable storage drivers including aufs, vfs, lvm and btrfs.
// * Pluggable execution drivers including lxc and chroot.
//
// In practice `daemon` still includes most core Docker components, including:
//
// * The reference registry client implementation
// * Image management
// * The build facility
// * Logging
//
// These components should be broken off into plugins of their own.
//
func daemon(eng *engine.Engine) error {
	return eng.Register("init_networkdriver", bridge.InitDriver)
}

// builtins jobs independent of any subsystem
func dockerVersion(job *engine.Job) engine.Status {
	v := &engine.Env{}
	v.SetJson("Version", dockerversion.VERSION)
	v.SetJson("ApiVersion", api.APIVERSION)
	v.Set("GitCommit", dockerversion.GITCOMMIT)
	v.Set("GoVersion", runtime.Version())
	v.Set("Os", runtime.GOOS)
	v.Set("Arch", runtime.GOARCH)
	if kernelVersion, err := kernel.GetKernelVersion(); err == nil {
		v.Set("KernelVersion", kernelVersion.String())
	}
	if _, err := v.WriteTo(job.Stdout); err != nil {
		return job.Error(err)
	}
	return engine.StatusOK
}
