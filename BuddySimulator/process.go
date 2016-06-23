package BuddySimulator

import (
	"sort"
	"fmt"
)

var nextPid int = 0

type Process struct {
	pid         int
	memoryUsage uint
	memoryBlock *Block
}

func NewProcess(memoryUsage uint) *Process {
	process := &Process{nextPid, memoryUsage, nil}
	nextPid++
	return process
}

func (p Process)GetPid() int {
	return p.pid
}

func (p Process)GetBlock() Block {
	return *p.memoryBlock
}

func (p *Process)extractBlock() *Block {
	tmp := p.memoryBlock
	p.memoryBlock = nil
	return tmp
}

func (p Process)String() string {
	res := "{pid:" + fmt.Sprint(p.GetPid()) +", "
	res += "memoryUsage:" + fmt.Sprint(p.memoryUsage) +", "
	res += fmt.Sprintf("%+v",p.GetBlock()) + "}"
	return res
}

type processList []*Process

func (p processList)String() string {
	res := "["
	for i := 0; i < p.Len(); i++ {
		res += fmt.Sprint(*p[i])
		if i != p.Len() - 1 {
			res += ", "
		}
	}
	res += "]"
	return res
}

func (p processList)Len() int {
	return len(p)
}

func (p processList)Less(i, j int) bool {
	return p[i].GetPid() < p[j].GetPid()
}

func (p processList)Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type OS struct {
	memory *Memory
	pList  processList
}

func NewOs(memory uint) (*OS, error) {
	var err error
	os := &OS{}
	os.memory, err = NewMemory(memory)
	if err != nil {
		return nil, err
	}
	os.pList = make([]*Process, 0)
	return os, nil
}

func (os *OS)AllocNewProcess(size uint) (int, error) {
	block, err := os.memory.AllocateMemory(size)
	if err != nil {
		return -1, err
	}
	process := NewProcess(size)
	process.memoryBlock = block
	os.pList = append(os.pList, process)
	sort.Sort(os.pList)
	return process.GetPid(), nil
}

func (os *OS)DeallocProcess(pid int) error {
	index := sort.Search(os.pList.Len(), func(n int) bool {
		if os.pList[n].GetPid() < pid {
			return false
		} else {
			return true
		}
	})

	if index == os.pList.Len() || os.pList[index].GetPid() != pid {
		return fmt.Errorf("pid %d is not pressent in OS", pid)
	}

	proc := os.pList[index]
	os.pList = append(os.pList[:index], os.pList[index + 1:]...)
	sort.Sort(os.pList)
	block := proc.extractBlock()
	err := os.memory.FreeBlock(block)
	if err != nil {
		return err
	}
	return nil
}

func (os *OS)PrintState() {
	freeBlocks := make([]Block, 0)
	for i := 0; i < len(os.memory.freeLists); i++ {
		for j := 0; j < len(os.memory.freeLists[i]); j++ {
			freeBlocks = append(freeBlocks, *os.memory.freeLists[i][j])
		}
	}

	fmt.Printf("freeBlocks:%+v, processes:%+v\n",freeBlocks,os.pList)

}



