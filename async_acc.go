package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type AsyncAcc []string

type WitnessPath []WitnessNode
type WitnessNode struct {
	hash string
	dir  Direction
}

type Direction int

const (
	LEFT Direction = iota
	RIGHT
)

func NewAsyncAcc() AsyncAcc {
	return AsyncAcc{"-"}
}

func (acc *AsyncAcc) Add(x string) (witPath WitnessPath) {
	// Copy previous accumulator
	newAcc := *acc

	d := 0
	z := hashSHA256(x)

	for newAcc[d] != "-" {
		if len(newAcc) < d+2 {
			newAcc = append(newAcc, "-")
		}

		z = hashSHA256(newAcc[d] + z)
		witPath = append(witPath, WitnessNode{newAcc[d], LEFT})
		newAcc[d] = "-"

		d++
	}
	newAcc[d] = z

	*acc = newAcc
	return
}

func (acc AsyncAcc) Verify(x string, witPath WitnessPath) bool {
	for _, a := range getAncestors(x, witPath) {
		for _, root := range acc {
			if root == a {
				return true
			}
		}
	}
	return false
}

func UpdateWitness(y string, witPathY, witPathX WitnessPath) (newWitPathX WitnessPath) {
	dx := len(witPathX)
	dy := len(witPathY)
	// No updates to witness
	if dy < dx {
		return witPathX
	}

	ancestorsY := getAncestors(y, witPathY)
	// Add ancestor and append rest of `witPathY`s path
	newWitPathX = witPathX
	newWitPathX = append(newWitPathX, WitnessNode{ancestorsY[dx], RIGHT})
	if dx+1 < len(witPathY) {
		newWitPathX = append(newWitPathX, witPathY[dx+1:]...)
	}

	return
}

func getAncestors(x string, witPath WitnessPath) []string {
	c := hashSHA256(x)
	ancestors := []string{c}

	for _, node := range witPath {
		if node.dir == LEFT {
			c = hashSHA256(node.hash + c)
		} else {
			c = hashSHA256(c + node.hash)
		}
		ancestors = append(ancestors, c)
	}

	return ancestors
}

func hashSHA256(x string) string {
	h256 := sha256.New()
	io.WriteString(h256, x)
	return base64.StdEncoding.EncodeToString(h256.Sum(nil))
}

func main() {
	numValues := (1 << 16) - 1
	acc := NewAsyncAcc()
	witnesses := []WitnessPath{}

	fmt.Print("[ADDING VALUES] |")
	for i := 0; i < numValues; i++ {
		wit := acc.Add(fmt.Sprintf("%d", i))
		witnesses = append(witnesses, wit)

		// Print progress bar
		divisor := numValues / 100
		if i%(10*divisor) == 0 && i != 0 {
			if i/divisor != 100 {
				fmt.Print(i / divisor)
			}
		} else if i%divisor == 0 {
			fmt.Print("=")
		}
	}

	for i, wit := range witnesses {
		fmt.Println("witness", i, wit)
	}

	fmt.Println("|\naccumulator for", numValues, "values:", acc)
}
