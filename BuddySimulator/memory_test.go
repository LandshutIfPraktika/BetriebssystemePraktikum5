package BuddySimulator

import (
	"testing"
)

func TestBlock_GetSize(t *testing.T) {

	memory, _ := NewMemory(512)

	blockHeight := uint(6)

	if size := CalculateBlockSize(memory, blockHeight); size != 128 {
		t.Errorf("size should be 128 got %d", size)
	}

}

func TestBlock_getBuddyAddress(t *testing.T) {

	leftBuddy := NewBlock(256,128)
	rightBuddy := NewBlock(384,128)

	if addr,left := leftBuddy.getBuddyAddressAndOwnPos(); addr != 384 || !left{
		t.Errorf("budy address should be 384 got %d", addr)
	}

	if addr,left := rightBuddy.getBuddyAddressAndOwnPos(); addr != 256 || left {
		t.Errorf("budy address should be 256 got %d", addr)
	}
}

func Test_calculateNextPowerOfTwo(t *testing.T) {

	if erg, height := calculateNextPowerOfTwoAndHeight(111); erg != 128 || height != 7 {
		t.Errorf("next power should be 128 got %d", erg)
	}

	if erg, height := calculateNextPowerOfTwoAndHeight(3); erg != 4  || height != 2{
		t.Errorf("next power should be 4 got %d", erg)
	}

	if erg,height := calculateNextPowerOfTwoAndHeight(16); erg != 16  || height != 4{
		t.Errorf("next power should be 16 got %d", erg)
	}
}

