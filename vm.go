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
	stack [1 << 16]uint16
	text  []byte

	pc uint16 // Program Counter
	sc uint16 // Stack Counter
	cc uint16 // Cycle Counter
	ar uint16 // Accumulator Register
}

func NewMachine(text []byte) *Machine {
	return &Machine{text: text}
}

func (m *Machine) push(w uint16) {
	m.stack[m.sc] = w
	m.sc++
}

func (m *Machine) pop() uint16 {
	m.sc--
	return m.stack[m.sc]
}

func (m *Machine) nextByte() byte {
	b := m.text[m.pc]
	m.pc++

	return b
}

func (m *Machine) NextWord() (uint16, error) {
	if int(m.pc)+1 >= len(m.text) {
		return 0, io.EOF
	}

	return *(*uint16)(unsafe.Pointer(
		&[2]byte{m.nextByte(), m.nextByte()},
	)), nil
}

const (
	POPI = iota // Pop to m.pc
	POPC        // Pop to m.cc

	PUSH // Push value onto stack
	DUP  // Dup value onto stack

	LOOP //  If m.cc > 0 then POPI with decriment of m.cc

	ADD // Add 2 value onto stack
	SUB // Sub 2 value onto stack

	INT_GET // Get value from stdout to stack
	INT_PUT // Put value from stack to stdout

	SWP // Swaps 2 value onto stack
)

func (m *Machine) Execute(r io.Reader, w io.Writer) error {
	for {
		if int(m.pc) >= len(m.text) {
			return io.EOF
		}

		switch m.nextByte() {
		case POPI:
			m.pc = m.pop()
		case POPC:
			m.cc = m.pop()
		case PUSH:
			w, err := m.NextWord()
			if err != nil {
				return err
			}
			m.push(w)
		case DUP:
			x := m.pop()
			m.push(x)
			m.push(x)
		case SWP:
			x, y := m.pop(), m.pop()
			m.push(x)
			m.push(y)
		case LOOP:
			if x := m.pop(); m.cc > 0 {
				m.cc--
				m.pc = x
			}
		case ADD:
			m.push(m.pop() + m.pop())
		case SUB:
			m.push(m.pop() - m.pop())
		case INT_PUT:
			fmt.Fprintln(w, m.pop())
		case INT_GET:
			fmt.Fscan(r, &m.ar)
			m.push(m.ar)
		default:
			return ErrUnsupportedOpcode
		}
	}
}
