package BuddySimulator

import (
	"math"
	"fmt"
)

type Memory struct {
	size      uint
	height    uint
	freeLists [][]*Block
}

func NewMemory(size uint) (*Memory, error) {
	floatDepth := math.Log2(float64(size))

	if math.Floor(floatDepth) != floatDepth {
		return nil, fmt.Errorf("size must be a power of 2; got %d", size)
	}

	newMemory := &Memory{}
	newMemory.size = size
	newMemory.height = uint(floatDepth) + 1
	newMemory.freeLists = make([][]*Block, newMemory.height)
	newMemory.freeLists[newMemory.height - 1] = append(newMemory.freeLists[newMemory.height - 1], NewBlock(0, size))

	return newMemory, nil
}

func (m Memory)String() string {
	res := ""
	res += "{"
	res += "size:" + fmt.Sprint(m.size)
	res += " height:" + fmt.Sprint(m.height)
	res += " freeLists:["
	for i := uint(0); i < m.height; i++ {
		res += "["
		for j, b := range m.freeLists[i] {
			res += fmt.Sprintf("%+v", *b)
			if j != len(m.freeLists[i]) - 1  {
				res += ", "
			}
		}
		res += "]"
	}
	res += "]}"
	return res

}

func (m *Memory)FreeBlock(block *Block) error {
	buddyAddress, left := block.getBuddyAddressAndOwnPos()
	_, height := calculateNextPowerOfTwoAndHeight(block.size)

	var (
		b *Block
		pos int
		found bool = false
	)
	for pos, b = range m.freeLists[height] {
		if b.GetAddress() == buddyAddress {
			found = true
			break
		}
	}

	if found {
		m.freeLists[height] = append(m.freeLists[height][:pos], m.freeLists[height][pos + 1:]...)
		if left {
			return m.FreeBlock(NewBlock(block.GetAddress(), 2 * block.GetSize()))
		} else {
			return m.FreeBlock(NewBlock(buddyAddress, 2 * block.GetSize()))
		}
	}
	m.freeLists[height] = append(m.freeLists[height], block)

	return nil
}

func (m *Memory)AllocateMemory(size uint) (block *Block, err error) {
	nextPower, height := calculateNextPowerOfTwoAndHeight(size)

	if (nextPower > m.size) {
		return nil, fmt.Errorf("can't allocate block of %d > memory capacity %d", size, m.size)
	}
	return m.makeBlockOfHeight(height)
}

func (m *Memory)makeBlockOfHeight(height uint) (*Block, error) {
	if len(m.freeLists[height]) == 0 {
		if err := m.splitBlockOfHeight(height + 1); err != nil {
			return nil, err
		}
	}
	block := m.freeLists[height][0]
	m.freeLists[height] = m.freeLists[height][1:]
	return block, nil

}

func (m*Memory)splitBlockOfHeight(height uint) error {
	if height >= m.height {
		return fmt.Errorf("end of memory height reached got %d, max height %d", height, m.height - 1)
	}
	if len(m.freeLists[height]) == 0 {
		if err := m.splitBlockOfHeight(height + 1); err != nil {
			return err
		}
	}

	blockSize := CalculateBlockSize(m, height - 1)
	splitBlock := m.freeLists[height][0]
	m.freeLists[height] = m.freeLists[height][1:]

	leftBuddy := NewBlock(splitBlock.GetAddress(), blockSize)
	rightBuddy := NewBlock(splitBlock.GetAddress() + blockSize, blockSize)
	m.freeLists[height - 1] = append(m.freeLists[height - 1], leftBuddy, rightBuddy)
	return nil
}

func calculateNextPowerOfTwoAndHeight(size uint) (nextPower, height uint) {
	nextPower = 1
	height = 0
	for nextPower < size {
		nextPower = nextPower << 1
		height++
	}
	return nextPower, height
}

type Block struct {
	address uint
	size    uint
}

func NewBlock(address, size uint) *Block {
	var b Block = Block{address, size}
	return &b
}

func (b Block)getBuddyAddressAndOwnPos() (buddyAddress uint, left bool) {
	blockSize := b.GetSize()
	address := b.GetAddress()
	if (address / blockSize) % 2 == 0 {
		return address + blockSize, true
	} else {
		return address - blockSize, false
	}
}

func CalculateBlockSize(memory *Memory, blockHeight uint) uint {
	return memory.size >> (memory.height - 1 - blockHeight)
}

func (b Block)GetAddress() uint {
	return b.address
}

func (b Block) GetSize() uint {
	return b.size
}

