// Copyright 2021 Mark Mandriota. All right reserved.

// package tinyvm - Tiny Stack-VM.
// This is only a demo work, so there are not many instruction and not complex arch.
package tinyvm

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

var ErrUnsupportedOpcode = errors.New("unsupported opcode")

type Machine struct {
	Stack [1 << 16]uint16
	Text  []byte

	PC uint16 // Program Counter
	SC uint16 // Stack Counter
	CC uint16 // Cycle Counter
	AR uint16 // Accumulator Register
}

func NewMachine(text []byte) *Machine {
	return &Machine{Text: text}
}

func (m *Machine) Push(w uint16) {
	m.Stack[m.SC] = w
	m.SC++
}

func (m *Machine) Pop() uint16 {
	m.SC--
	return m.Stack[m.SC]
}

func (m *Machine) nextByte() byte {
	b := m.Text[m.PC]
	m.PC++

	return b
}

func (m *Machine) NextWord() (uint16, error) {
	if int(m.PC)+1 >= len(m.Text) {
		return 0, io.EOF
	}

	return *(*uint16)(unsafe.Pointer(
		&[2]byte{m.nextByte(), m.nextByte()},
	)), nil
}

const (
	POPI = iota // Pop to m.PC
	POPC        // Pop to m.CC

	PUSH // Push value onto stack
	DUP  // Dup value onto stack

	LOOP //  If m.CC > 0 then POPI with decriment of m.CC

	ADD // Add 2 value onto stack
	SUB // Sub 2 value onto stack

	INT_GET // Get value from stdout to stack
	INT_PUT // Put value from stack to stdout

	SWP // Swaps 2 value onto stack
)

func (m *Machine) Execute(r io.Reader, w io.Writer) error {
	for {
		if int(m.PC) >= len(m.Text) {
			return io.EOF
		}

		switch m.nextByte() {
		case POPI:
			m.PC = m.Pop()
		case POPC:
			m.CC = m.Pop()
		case PUSH:
			w, err := m.NextWord()
			if err != nil {
				return err
			}
			m.Push(w)
		case DUP:
			x := m.Pop()
			m.Push(x)
			m.Push(x)
		case SWP:
			x, y := m.Pop(), m.Pop()
			m.Push(x)
			m.Push(y)
		case LOOP:
			if x := m.Pop(); m.CC > 0 {
				m.CC--
				m.PC = x
			}
		case ADD:
			m.Push(m.Pop() + m.Pop())
		case SUB:
			m.Push(m.Pop() - m.Pop())
		case INT_PUT:
			fmt.Fprintln(w, m.Pop())
		case INT_GET:
			fmt.Fscan(r, &m.AR)
			m.Push(m.AR)
		default:
			return ErrUnsupportedOpcode
		}
	}
}
