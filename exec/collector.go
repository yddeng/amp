package exec

import (
	"amp/util"
	"encoding/json"
	"fmt"
	psCpu "github.com/shirou/gopsutil/cpu"
	psDisk "github.com/shirou/gopsutil/disk"
	psHost "github.com/shirou/gopsutil/host"
	psMem "github.com/shirou/gopsutil/mem"
	psNet "github.com/shirou/gopsutil/net"
	"log"
	"sync"
	"time"
)

const (
	NET_INTERFACE_ALL = "all"
	NET_INTERFACE_VPN = "tun0"
)

type NetCollector struct {
	updateInterval    time.Duration
	totalBytesRecv    uint64
	totalBytesSent    uint64
	totalPacketsRecv  uint64
	totalPacketsSent  uint64
	recentBytesRecv   uint64
	recentBytesSent   uint64
	recentPacketsRecv uint64
	recentPacketsSent uint64
	tcpConnections    int
	udpConnections    int
	sync.RWMutex
}

func NewNetCollector() *NetCollector {
	self := &NetCollector{
		updateInterval: time.Second,
	}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

func (self *NetCollector) update() {
	interfaces, err := psNet.IOCounters(true)
	if err != nil {
		log.Printf("failed to get network activity from collector: %v", err)
		return
	}

	var totalBytesRecv, totalBytesSent uint64
	var totalPacketsRecv, totalPacketsSent uint64
	for _, _interface := range interfaces {
		// ignore VPN interface or filter interface by name
		if _interface.Name != NET_INTERFACE_VPN {
			totalBytesRecv += _interface.BytesRecv
			totalBytesSent += _interface.BytesSent
			totalPacketsRecv += _interface.PacketsRecv
			totalPacketsSent += _interface.PacketsSent
		}
	}

	var recentBytesRecv, recentBytesSent uint64
	var recentPacketsRecv, recentPacketsSent uint64

	if self.totalBytesRecv != 0 { // if this isn't the first update
		recentBytesRecv = totalBytesRecv - self.totalBytesRecv
		recentBytesSent = totalBytesSent - self.totalBytesSent
		recentPacketsRecv = totalPacketsRecv - self.totalPacketsRecv
		recentPacketsSent = totalPacketsSent - self.totalPacketsSent

		if int(recentBytesRecv) < 0 {
			log.Printf("error: negative value for recently received network data from collector. recentBytesRecv: %v", recentBytesRecv)
			recentBytesRecv = 0
		}
		if int(recentBytesSent) < 0 {
			log.Printf("error: negative value for recently sent network data from collector. recentBytesSent: %v", recentBytesSent)
			recentBytesSent = 0
		}
		if int(recentPacketsRecv) < 0 {
			log.Printf("error: negative value for recently received network data from collector. recentPacketsRecv: %v", recentPacketsRecv)
			recentPacketsRecv = 0
		}
		if int(recentPacketsSent) < 0 {
			log.Printf("error: negative value for recently sent network data from collector. recentPacketsSent: %v", recentPacketsSent)
			recentPacketsSent = 0
		}

		self.recentBytesRecv = recentBytesRecv
		self.recentBytesSent = recentBytesSent
		self.recentPacketsRecv = recentPacketsRecv
		self.recentPacketsSent = recentPacketsSent
	}

	// used in later calls to update
	self.totalBytesRecv = totalBytesRecv
	self.totalBytesSent = totalBytesSent
	self.totalPacketsRecv = totalPacketsRecv
	self.totalPacketsSent = totalPacketsSent

	if conns, err := psNet.ConnectionsPid("tcp", 0); err == nil {
		self.tcpConnections = len(conns)
	}
	if conns, err := psNet.ConnectionsPid("udp", 0); err == nil {
		self.udpConnections = len(conns)
	}
}

func (self *NetCollector) String() string {
	ret := self.Result()
	data, _ := json.Marshal(ret)
	return string(data)
}

func (self *NetCollector) Result() map[string]string {
	self.RLock()
	defer self.RUnlock()

	convertBytes := func(num uint64) string {
		converted, unit := util.ConvertBytes(num)
		return fmt.Sprintf("%.1f%s", converted, unit)
	}

	return map[string]string{
		"totalBytesRecv":    convertBytes(self.totalBytesRecv),
		"totalBytesSent":    convertBytes(self.totalBytesSent),
		"recentBytesRecv":   fmt.Sprintf("%s/s", convertBytes(self.recentBytesRecv)),
		"recentBytesSent":   fmt.Sprintf("%s/s", convertBytes(self.recentBytesSent)),
		"totalPacketsRecv":  fmt.Sprintf("%d", self.totalPacketsRecv),
		"totalPacketsSent":  fmt.Sprintf("%d", self.totalPacketsSent),
		"recentPacketsRecv": fmt.Sprintf("%d/s", self.recentPacketsRecv),
		"recentPacketsSent": fmt.Sprintf("%d/s", self.recentPacketsSent),
		"tcpConnections":    fmt.Sprintf("%d", self.tcpConnections),
		"udpConnections":    fmt.Sprintf("%d", self.udpConnections),
	}
}

type CPUCollector struct {
	updateInterval time.Duration
	cpuCores       int
	usedPercent    float64
	sync.RWMutex
}

func NewCPUCollector() *CPUCollector {
	self := &CPUCollector{
		updateInterval: time.Second,
	}

	var err error
	self.cpuCores, err = psCpu.Counts(false)
	if err != nil {
		log.Printf("failed to get CPU count from collector: %v", err)
	}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.update()
		}
	}()

	return self
}

func (self *CPUCollector) update() {
	go func() {
		percent, err := psCpu.Percent(self.updateInterval, false)
		if err != nil {
			log.Printf("failed to get average CPU usage percent from collector: %v. self.updateInterval: %v. percpu: %v", err, self.updateInterval, false)
		} else {
			self.Lock()
			defer self.Unlock()
			self.usedPercent = percent[0]
		}
	}()

}

func (self *CPUCollector) String() string {
	ret := self.Result()
	data, _ := json.Marshal(ret)
	return string(data)
}

func (self *CPUCollector) Result() map[string]string {
	self.RLock()
	defer self.RUnlock()
	return map[string]string{
		"cpuCores":    fmt.Sprintf("%d", self.cpuCores),
		"usedPercent": fmt.Sprintf("%.1f%%", self.usedPercent),
	}
}

type MemCollector struct {
	updateInterval     time.Duration
	virtualTotal       uint64
	virtualUsed        uint64
	virtualUsedPercent float64
	swapTotal          uint64
	swapUsed           uint64
	swapUsedPercent    float64
	sync.RWMutex
}

func NewMemCollector() *MemCollector {
	self := &MemCollector{
		updateInterval: time.Second,
	}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

func (self *MemCollector) update() {
	mainMemory, err := psMem.VirtualMemory()
	if err != nil {
		log.Printf("failed to get main memory info from collector: %v", err)
	} else {
		self.virtualTotal = mainMemory.Total
		self.virtualUsed = mainMemory.Used
		self.virtualUsedPercent = mainMemory.UsedPercent
	}

	swapMemory, err := psMem.SwapMemory()
	if err != nil {
		log.Printf("failed to get swap memory info from collector: %v", err)
	} else {
		self.swapTotal = swapMemory.Total
		self.swapUsed = swapMemory.Used
		self.swapUsedPercent = swapMemory.UsedPercent
	}
}

func (self *MemCollector) String() string {
	ret := self.Result()
	data, _ := json.Marshal(ret)
	return string(data)
}

func (self *MemCollector) Result() map[string]string {
	self.RLock()
	defer self.RUnlock()
	convertBytes := func(num uint64) string {
		converted, unit := util.ConvertBytes(num)
		return fmt.Sprintf("%.1f%s", converted, unit)
	}
	return map[string]string{
		"virtualTotal":       convertBytes(self.virtualTotal),
		"virtualUsed":        convertBytes(self.virtualUsed),
		"virtualUsedPercent": fmt.Sprintf("%.1f%%", self.virtualUsedPercent),
		"swapTotal":          convertBytes(self.swapTotal),
		"swapUsed":           convertBytes(self.swapUsed),
		"swapUsedPercent":    fmt.Sprintf("%.1f%%", self.swapUsedPercent),
	}
}

type HostCollector struct {
	hostname string
	os       string
	arch     string
}

func NewHostCollector() *HostCollector {
	self := &HostCollector{}
	info, err := psHost.Info()
	if err != nil {
		log.Printf("failed to get host info from collector: %v", err)
	} else {
		self.hostname = info.Hostname
		self.os = info.OS
		self.arch = info.KernelArch
		//self.bootTime = info.BootTime
	}
	return self
}

func (self *HostCollector) String() string {
	ret := self.Result()
	data, _ := json.Marshal(ret)
	return string(data)
}

func (self *HostCollector) Result() map[string]string {
	return map[string]string{
		"hostname": self.hostname,
		"os":       self.os,
		"arch":     self.arch,
	}
}

type DiskCollector struct {
	updateInterval   time.Duration
	mounted          map[string]string
	total            uint64
	used             uint64
	avail            uint64
	free             uint64
	usedPercent      float64
	recentBytesRead  uint64
	recentBytesWrite uint64
	sync.RWMutex
}

func NewDiskCollector() *DiskCollector {
	self := &DiskCollector{
		updateInterval: time.Second,
		mounted: map[string]string{
			"/": "",
		},
	}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

func (self *DiskCollector) update() {
	stat, err := psDisk.Usage("/")
	if err != nil {
		log.Printf("failed to get disk usage from collector: %v", err)
	} else {
		// 依赖库写法与系统不一致
		self.total = stat.Total
		self.free = stat.Total - stat.Used
		self.avail = stat.Free
		self.used = self.total - self.free
		self.usedPercent = float64(self.used) * 100 / float64(self.avail+self.used)

		psDisk.Partitions()
	}
}

func (self *DiskCollector) String() string {
	ret := self.Result()
	data, _ := json.Marshal(ret)
	return string(data)
}

func (self *DiskCollector) Result() map[string]string {
	self.RLock()
	defer self.RUnlock()
	convertBytes := func(num uint64) string {
		converted, unit := util.ConvertBytes(num)
		return fmt.Sprintf("%.1f%s", converted, unit)
	}
	return map[string]string{
		"total":       convertBytes(self.total),
		"used":        convertBytes(self.used),
		"avail":       convertBytes(self.avail),
		"free":        convertBytes(self.free),
		"usedPercent": fmt.Sprintf("%.1f%%", self.usedPercent),
	}
}
