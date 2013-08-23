package system

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type SystemStatHistory struct {
	Time string
	Stat *SystemStat
}

type SystemStat struct {
	Load *LoadStat
	Mem  *MemStat
}

type LoadStat struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

type MemStat struct {
	Used   float64
	Cached float64
	Free   float64
}

func GetSystemStat() *SystemStat {
	stat := new(SystemStat)
	top_output := top()
	for _, line := range strings.Split(top_output, "\n") {
		if strings.HasPrefix(line, "Load Avg:") {
			// ex) Load Avg: 1.10, 1.25, 1.29
			re := regexp.MustCompile("[0-9]+.[0-9]+")
			usage := re.FindAllString(line, -1)

			load1, _ := strconv.ParseFloat(usage[0], 64)
			load5, _ := strconv.ParseFloat(usage[1], 64)
			load15, _ := strconv.ParseFloat(usage[2], 64)

			stat.Load = &LoadStat{
				Load1:  load1,
				Load5:  load5,
				Load15: load15,
			}
		} else if strings.HasPrefix(line, "PhysMem") {
			// ex) PhysMem: 1293M wired, 3782M active, 1281M inactive, 6357M used, 1834M free.
			re := regexp.MustCompile("[0-9]+")
			mems := re.FindAllString(line, -1)

			wired, _ := strconv.ParseFloat(mems[0], 64)
			active, _ := strconv.ParseFloat(mems[1], 64)
			inactive, _ := strconv.ParseFloat(mems[2], 64)
			free, _ := strconv.ParseFloat(mems[4], 64)

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
