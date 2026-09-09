# 04: Boot Choreography {#1}

Wednesday. After breakfast the crew wheels the big table onto the lawn
and unpacks a crate of palm-sized tracked robots the group has been
building all year: LED strip down the spine, a hobby servo at each end,
and a house-built SBC - the *GT-Nano* - soldered in the middle.

Every unit ships with a factory boot ROM. When the robot powers on it
runs the ROM, streams its `6` "opening pose" numbers onto the
choreography bus, and clunks into position ready for the day's dance.
Nobody has the schematic any more; all you have is the assembly listing
that the boot ROM was compiled from.

The GT-Nano is a tiny virtual machine with:

* *six signed-integer registers* `A B C D E F`, all initially zero;
* an *infinite tape* of signed-integer cells indexed by a signed
  cursor `cur`, all cells initially zero;
* an *output stream* of numbers, initially empty.

## Instruction set

Every instruction takes zero, one or two operands. An operand is either
a **register** (single upper-case letter `A..F`) or a signed
**immediate** integer.

| instruction | effect                          |
| ----------- | ------------------------------- |
| `LEFT`      | `cur -= 1`                      |
| `RIGHT`     | `cur += 1`                      |
| `LOAD R`    | `R = mem[cur]`                  |
| `SAVE V`    | `mem[cur] = V`                  |
| `SET R V`   | `R = V`                         |
| `ADD R V`   | `R += V`                        |
| `SUB R V`   | `R -= V`                        |
| `JMP V`     | jump by `V`                     |
| `JNZ R V`   | if `R != 0`, jump by `V`        |
| `OUT V`     | append `V` to the output stream |

Lines are numbered `0, 1, 2, ...` from the top of the program, and the
program counter `pc` starts at line `0`. Non-jump instructions advance
`pc` by one line. A jump instruction reads a non-zero signed offset `V`
and sets `pc = pc + V`, where `pc` is the address of the *jump
instruction itself*. `JMP 1` is therefore a no-op, `JMP 2` skips the
next line, and `JMP -1` jumps to the previous line. Execution halts
when `pc` walks off the end of the program.

## Input

Your puzzle input is the assembly program, one instruction per line.

## Part 1

Boot the program with all registers, memory cells and the cursor at
`0` and run it to completion. It emits exactly `6` output values
`V0 V1 V2 V3 V4 V5` in that order. Return the *positional
base-10 digest*

    1×V0 + 10×V1 + 100×V2 + 1000×V3 + 10000×V4 + 100000×V5

### Example

Consider a small program that emits `2` outputs (instead of `6`), so its
scaled-down digest is `1×V0 + 10×V1`:

```
SET A 1
LOAD B
SET C 2
ADD A B
SUB C 1
JNZ C -2
RIGHT
LOAD B
ADD A B
OUT A
SET A -1
LOAD B
ADD A B
LEFT
LOAD B
ADD A B
OUT A
```

Booted with all-zero memory the program emits `V0 = 1` and
`V1 = -1`, giving a digest of *1×1 + 10×(-1)* = `-9`.
