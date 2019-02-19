package main

import (
	"flag"
	"log"
	"mg-coms/service"
	"os"
	"runtime/pprof"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file.")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {

		f, err := os.Create(*cpuprofile)

		if err != nil {

			log.Fatal(err)

		}

		pprof.StartCPUProfile(f)

		defer pprof.StopCPUProfile()

	}

	server1 := service.New()
	server1.Start()

	// coms, err := coms.New()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// coms.Start()
	// time.Sleep(1 * time.Second)
	// coms.TestSend()
	// time.Sleep(1 * time.Second)
	// coms.DeviceMap.Print()
	// coms.Stop()
	//
	// time.Sleep(1 * time.Second)

}

//Process　进程占用cpu
// type Process struct {
// 	pid int
// 	cpu float64
// 	mem float64
// }
//
// func process() {
// 	cmd := exec.Command("ps", "aux")
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	processes := make([]*Process, 0)
// 	for {
// 		line, err := out.ReadString('\n')
// 		if err != nil {
// 			break
// 		}
// 		tokens := strings.Split(line, " ")
// 		ft := make([]string, 0)
// 		for _, t := range tokens {
// 			if t != "" && t != "\t" {
// 				ft = append(ft, t)
// 			}
// 		}
// 		// log.Println(len(ft), ft)
// 		pid, err := strconv.Atoi(ft[1])
// 		if err != nil {
// 			continue
// 		}
// 		cpu, err := strconv.ParseFloat(ft[2], 64)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		mem, err := strconv.ParseFloat(ft[3], 64)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		processes = append(processes, &Process{pid, cpu, mem})
// 	}
// 	var sum, num float64
// 	for _, p := range processes {
// 		sum += p.cpu
// 		num += p.mem
// 	}
// 	log.Println("CPU ", sum, " MEM ", num)
//
// }

// func main() {
// 	// coms := coms.New()
// 	// coms.Commander = command.Command{}
// 	var c1 channel.ChannelConfig
// 	var c2 channel.ChannelConfig
// 	var m command.Message
// 	// c.Name = "通道1"
// 	// ch := make(chan int)
// 	// go func() {
// 	// 	ch <- 1
// 	// }()
// 	// for {
// 	// 	select {
// 	// 	case <-ch:
// 	// 		c.Start()
// 	// 		return
// 	// 	}
// 	// }
// 	m.Init()
// 	switch m.Cmd {
// 	case 0:
// 		fmt.Println("退出")
// 	case 1:
// 		fmt.Println(m.Value)
// 		fmt.Println(m.Time)
// 		c1.Start()
// 		c2.Start()
// 	case 2:
// 		fmt.Println("重新启动")
// 	}
// }
