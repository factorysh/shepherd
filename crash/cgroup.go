package crash

import (
	"path"

	"io/ioutil"
	"runtime"

	"github.com/docker/docker/api/types"

	log "github.com/sirupsen/logrus"
)

// /sys/fs/cgroup/memory/system.slice/containerd.service

type Cgroup struct {
	mount string // /sys/fs/cgroup
	root  string // system.slice/containerd.service
}

func NewCgroup() *Cgroup {
	return &Cgroup{
		mount: "/sys/fs/cgroup",
		root:  "system.slice/containerd.service",
	}
}

func (c *Cgroup) Path(group, key string, container *types.ContainerJSON) string {
	return path.Join(
		c.mount,
		group,
		c.root,
		container.HostConfig.CgroupParent,
		container.HostConfig.Cgroup.Container(),
		key)
}

func (c *Cgroup) readPath(group, key string, container *types.ContainerJSON) string {
	v, err := ioutil.ReadFile(c.Path(group, key, container))
	if err != nil {
		log.WithFields(log.Fields{
			"group":     group,
			"key":       key,
			"container": container.Name,
		}).WithError(err).Error("Reading cgroup")
		return ""
	}
	return string(v)
}

func (c *Cgroup) fetchCgroupStates(container *types.ContainerJSON) map[string]string {
	states := make(map[string]string)
	if runtime.GOOS != "linux" {
		return states
	}
	states["memory.max_usage_in_bytes"] = c.readPath("memory", "memory.max_usage_in_bytes", container)
	states["memory.memsw.max_usage_in_bytes"] = c.readPath("memory", "memory.memsw.max_usage_in_bytes", container)
	states["memory.kmem.max_usage_in_bytes"] = c.readPath("memory", "memory.kmem.max_usage_in_bytes", container)
	states["memory.kmem.tcp.max_usage_in_bytes"] = c.readPath("memory", "memory.kmem.tcp.max_usage_in_bytes", container)

	states["cpuacct.usage_all"] = c.readPath("cpu,cpuacct", "cpuacct.usage_all", container)

	states["blkio.io_service_time"] = c.readPath("blkio", "blkio.io_service_time", container)
	states["blkio.io_wait_time"] = c.readPath("blkio", "blkio.io_wait_time", container)
	states["blkio.io_serviced"] = c.readPath("blkio", "blkio.io_serviced", container)

	return states
}
