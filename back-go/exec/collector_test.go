package exec

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	netC := NewNetCollector()
	cpuC := NewCPUCollector()
	memC := NewMemCollector()
	hostC := NewHostCollector()
	diskC := NewDiskCollector()

	time.Sleep(time.Second * 2)
	fmt.Println(netC.String())
	fmt.Println()

	fmt.Println(cpuC.String())
	fmt.Println()

	fmt.Println(memC.String())
	fmt.Println()

	fmt.Println(hostC.String())
	fmt.Println()

	fmt.Println(diskC.String())
	fmt.Println()

}
