package system

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type SystemStat struct {
	Cpu *CpuStat
	Mem *MemStat
}

type CpuStat struct {
	Sys  float64
	User float64
	Idle float64
}

type MemStat struct {
	Used   int64
	Cached int64
	Free   int64
}

func GetSystemStat() *SystemStat {
	stat := new(SystemStat)
	top_output := top()
	for _, line := range strings.Split(top_output, "\n") {
		if strings.HasPrefix(line, "CPU") {
			// ex) CPU usage: 9.37% user, 53.12% sys, 37.50% idle
			re := regexp.MustCompile("[0-9]+.[0-9]+")
			usage := re.FindAllString(line, -1)

			sys, _ := strconv.ParseFloat(usage[0], 64)
			user, _ := strconv.ParseFloat(usage[1], 64)
			idle, _ := strconv.ParseFloat(usage[2], 64)

			stat.Cpu = &CpuStat{
				Sys:  sys,
				User: user,
				Idle: idle,
			}
		} else if strings.HasPrefix(line, "PhysMem") {
			// ex) PhysMem: 1293M wired, 3782M active, 1281M inactive, 6357M used, 1834M free.
			re := regexp.MustCompile("[0-9]+")
			mems := re.FindAllString(line, -1)

			wired, _ := strconv.ParseInt(mems[0], 10, 64)
			active, _ := strconv.ParseInt(mems[1], 10, 64)
			inactive, _ := strconv.ParseInt(mems[2], 10, 64)
			free, _ := strconv.ParseInt(mems[4], 10, 64)

			stat.Mem = &MemStat{
				Used:   wired + active,
				Cached: inactive,
				Free:   free,
			}
		}
	}

	return stat
}

func top() string {
	out, err := execute("top", "-l 1", "-n 0", "-s 0")
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return out
}

func execute(cmd string, opts ...string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s not found.", cmd))
		return "", err
	}

	output, err := exec.Command(path, opts...).Output()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(output), nil

}
