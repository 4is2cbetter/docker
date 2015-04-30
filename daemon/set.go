package daemon

import (
	"github.com/docker/docker/engine"
	"github.com/docker/docker/runconfig"

	log "github.com/Sirupsen/logrus"
)

func (daemon *Daemon) ContainerSet(job *engine.Job) engine.Status {
	name := job.Args[0]
	hostConfig := runconfig.ContainerHostConfigFromJob(job)

	warnings, err := daemon.verifyHostConfig(hostConfig)
	if err != nil {
		return job.Error(err)
	}

	for w := range warnings {
		log.Warn(w)
	}

	container, err := daemon.Get(name)
	if err != nil {
		return job.Error(err)
	}

	daemon.updateResources(container, hostConfig)

	if err := daemon.Set(container.command); err != nil {
		return job.Error(err)
	}

	container.LogEvent("set")
	return engine.StatusOK
}

func (daemon *Daemon) updateResources(container *Container, hostConfig *runconfig.HostConfig) {
	if hostConfig.CpuShares != 0 {
		container.hostConfig.CpuShares = hostConfig.CpuShares
	}
	if hostConfig.CpusetCpus != "" {
		container.hostConfig.CpusetCpus = hostConfig.CpusetCpus
	}
	if hostConfig.Memory != 0 {
		container.hostConfig.Memory = hostConfig.Memory
	}
	if hostConfig.MemorySwap != 0 {
		container.hostConfig.MemorySwap = hostConfig.MemorySwap
	}

	container.command.Resources.CpuShares = container.hostConfig.CpuShares
	container.command.Resources.CpusetCpus = container.hostConfig.CpusetCpus
	container.command.Resources.Memory = container.hostConfig.Memory
	container.command.Resources.MemorySwap = container.hostConfig.MemorySwap
}
