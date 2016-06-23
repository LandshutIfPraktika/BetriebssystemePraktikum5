package main

import (
	"github.com/s-gheldd/BuddySimulator"
	"log"
	"gopkg.in/cheggaaa/pb.v1"
	"time"
	"fmt"
)

func main() {

	fmt.Print("How much memory should your new os manage (power of 2):")
	memorySize := uint(0)
	fmt.Scanf("%d", &memorySize)

	OS, err := BuddySimulator.NewOs(memorySize)
	if err != nil {
		log.Fatal(err)
	}

	count := 1000

	bar := pb.New(count)
	bar.RefreshRate = time.Millisecond
	bar.Prefix("Booting OS")
	bar.ShowCounters = false
	bar.Start()
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("Done")

	loop:
	for true {
		fmt.Print("Do you want to:\n allocate a process(1)\n free a process(2)\n quit(3)? ")
		answer := 0
		fmt.Scanf("%d", &answer)
		cases:
		switch answer {
		case 1:
			fmt.Print("Enter memory size: ")
			size := uint(0)
			fmt.Scanf("%d", &size)
			pid, err := OS.AllocNewProcess(size)
			if err != nil {
				log.Print(err)
			}
			fmt.Println("Createt process:", pid)
		case 2:
			fmt.Print("Enter pid: ")
			pid := 0
			fmt.Scanf("%d", &pid)
			err := OS.DeallocProcess(pid)
			if err != nil {
				log.Print(err)
			}
		case 3:
			break loop
		default:
			break cases
		}
		OS.PrintState()
		fmt.Println()
	}

}
