package process

import (
	"bytes"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//Process　进程占用cpu
type Process struct {
	Pid int
	Cpu float64
	Mem float64
}

//ProcessGet 获取进程状态
func (p Process) ProcessGet() (temp, data float64) {
	if runtime.GOOS == "windows" {
		return 0, 0
	}
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		// log.Println(len(ft), ft)
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		mem, err := strconv.ParseFloat(ft[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		processes = append(processes, &Process{pid, cpu, mem})
	}
	var sum, num float64
	for _, p := range processes {
		if p.Cpu >= sum {
			sum = p.Cpu
		}
		// sum += p.Cpu
		num += p.Mem
	}
	// log.Println("CPU ", sum, " MEM ", num)
	return sum, num

}
