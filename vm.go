// Copyright 2021 Mark Mandriota. All right reserved.

// package tinyvm - Tiny Stack-VM.
// This is just a demo work, so there are not many instruction and not complex arch.
package tinyvm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	u "unsafe"
)

type Word [2]byte

func (w *Word) Uint16() uint16 {
	return *(*uint16)(u.Pointer(w))
}

type Machine struct {
	Stack [1 << 16]uint16
	Text  []byte

	out *bufio.Writer
	in  *bufio.Reader

	Reg struct {
		IP uint16
		SP uint16
		CX uint16
		AX Word
	}
}

func NewMachine(fi []byte, w io.Writer, r io.Reader) *Machine {
	return &Machine{
		Text: fi,
		out:  bufio.NewWriter(w),
		in:   bufio.NewReader(r),
	}
}

func (m *Machine) Push(v uint16) {
	m.Stack[m.Reg.SP] = v
	m.Reg.SP++
}

func (m *Machine) Pop() uint16 {
	m.Reg.SP--
	return m.Stack[m.Reg.SP]
}

func (m *Machine) NextByte() (v byte) {
	if int(m.Reg.IP) >= len(m.Text) {
		panic(errors.New("end of file"))
	}

	v = m.Text[m.Reg.IP]
	m.Reg.IP++

	return
}

func (m *Machine) NextWord() uint16 {
	m.Reg.AX[0] = m.NextByte()
	m.Reg.AX[1] = m.NextByte()

	return m.Reg.AX.Uint16()
}

const (
	POPI = iota // Pop to m.Reg.IP
	POPC        // Pop to m.Reg.CX

	PUSH // Push value onto stack
	DUP  // Dup value onto stack

	LOOP //  If m.Reg.CX > 0 then POPI with decriment of m.Reg.CX

	ADD // Add 2 value onto stack
	SUB // Sub 2 value onto stack

	INT_GET // Get value from stdout to stack
	INT_PUT // Put value from stack to stdout

	SWP // Swaps 2 value onto stack
)

func (m *Machine) Execute() {
	defer m.out.Flush()
	for {
		switch m.NextByte() {
		case POPI:
			m.Reg.IP = m.Pop()
		case POPC:
			m.Reg.CX = m.Pop()
		case PUSH:
			m.Push(m.NextWord())
		case DUP:
			x := m.Pop()
			m.Push(x)
			m.Push(x)
		case SWP:
			x, y := m.Pop(), m.Pop()
			m.Push(x)
			m.Push(y)
		case LOOP:
			if x := m.Pop(); m.Reg.CX > 0 {
				m.Reg.CX--
				m.Reg.IP = x
			}
		case ADD:
			m.Push(m.Pop() + m.Pop())
		case SUB:
			m.Push(m.Pop() - m.Pop())
		case INT_PUT:
			fmt.Fprintln(m.out, m.Pop())
		case INT_GET:
			x := uint16(0)
			fmt.Fscan(m.in, &x)
			m.Push(x)
		default:
			panic(errors.New("unsupported opcode"))
		}
	}
}
