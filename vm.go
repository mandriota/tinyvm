// Copyright 2021 Mark Mandriota. All right reserved.

// package tinyvm - Tiny Stack-VM.
// This is just a demo work, so there are not many instruction and not complex arch.
package tinyvm

import (
	"fmt"
	"io"
	u "unsafe"
)

type Word [2]byte

func (w *Word) Uint16() uint16 {
	return *(*uint16)(u.Pointer(w))
}

type Machine struct {
	Stack []uint16
	Text  []byte

	Reg struct {
		IP uint16
		CX uint16
	}
}

func (m *Machine) Init(fi []byte) {
	m.Text = fi
	m.Reg.IP = 0
	m.Reg.CX = 0

	m.Stack = make([]uint16, 0, 1<<16)
	m.Push(uint16(len(m.Text) - 1))
}

func (m *Machine) Push(v uint16) {
	m.Stack = append(m.Stack, v)
}

func (m *Machine) Pop() (v uint16) {
	if len(m.Stack) < 1 {
		panic("pop onto empty stack")
	}

	v = m.Stack[len(m.Stack)-1]
	m.Stack = m.Stack[:len(m.Stack)-1]

	return
}

func (m *Machine) NextByte() (v byte) {
	if int(m.Reg.IP) >= len(m.Text) {
		panic("EOF")
	}

	v = m.Text[m.Reg.IP]
	m.Reg.IP++
	return
}

func (m *Machine) NextWord() uint16 {
	var w Word
	w[0] = m.NextByte()
	w[1] = m.NextByte()

	return w.Uint16()
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
)

func (m *Machine) Execute(r io.Reader, w io.Writer) {
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
			fmt.Fprintln(w, m.Pop())
		case INT_GET:
			var x uint16
			fmt.Fscan(r, &x)
			m.Push(x)
		}
	}
}
