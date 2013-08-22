package system

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type MemInfo struct {
	Free     int64
	Active   int64
	Wired    int64
	Inactive int64
}

func GetMemInfo() *MemInfo {
	m := new(MemInfo)
	out, err := exec.Command("vm_stat").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	lines = lines[1 : len(lines)-2] // exclude header/footer
	for _, line := range lines {
		rows := strings.Split(line, ":")
		val := rows[1]
		val = strings.TrimLeft(val, " ")
		val = strings.TrimRight(val, ".")
		int_val, _ := strconv.ParseInt(val, 10, 64)
		int_val = int_val * 4096        // page to Byte
		int_val = int_val / 1024 / 1024 // to MByte
		switch rows[0] {
		case "Pages free":
			m.Free = int_val
		case "Pages speculative":
			m.Free += int_val
		case "Pages active":
			m.Active = int_val
		case "Pages inactive":
			m.Inactive = int_val
		case "Pages wired down":
			m.Wired = int_val
		}
	}
	return m
}
