# FFLT lang

FFLT lang is essentially isomorphic to esoteric language Whitespace.  
FFLT lang use F, L, T instead of SPACE, TAB, LF, and all other characters are ignored.  
No case sensitivity except for label parameter.

## Usage

```
fflt_lang program.fflt
```

## Building yourself

```
make build
```

## FFLT lang specification

Each command consists of a series of tokens, beginning with the Instruction Modification Parameter (IMP).  
These are listed in the table below.

| IMP | Meaning            |
| --- | ------------------ |
| F   | Stack Manipulation |
| LF  | Arithmeric         |
| LL  | Heap access        |
| T   | Flow Control       |
| LT  | I/O                |

### Stack Manipulation (IMP:[F])

| Command | Parameters | Meaning                                                                            |
| ------- | ---------- | ---------------------------------------------------------------------------------- |
| F       | Number     | Push the number onto the stack                                                     |
| LF      | Number     | Copy the *n*th item on the stack (given by the argument) onto the top of the stack |
| LT      | Number     | Slide _n_ items off the stack, keeping the top item                                |
| TF      | -          | Duplicate the top item on the stack                                                |
| TL      | -          | Swap the top twe item on the stack                                                 |
| TT      | -          | Discard the top item on the stack                                                  |

### Arithmetic (IMP:[LF])

| Command | Parameters | Meaning          |
| ------- | ---------- | ---------------- |
| FF      | -          | Addition         |
| FL      | -          | Subtraction      |
| FT      | -          | Multiplication   |
| LF      | -          | Integer Division |
| LL      | -          | Modulo           |

### Heap Access (IMP:[LL])

| Command | Parameters | Meaning  |
| ------- | ---------- | -------- |
| F       | -          | Store    |
| L       | -          | Retrieve |

### Flow Control (IMP:[T])

| Command | Parameters | Meaning                                                |
| ------- | ---------- | ------------------------------------------------------ |
| FF      | Label      | Mark a location in the program                         |
| FL      | Label      | Call a subroutine                                      |
| FT      | Label      | Jump unconditionally to a label                        |
| LF      | Label      | Jump to a label if the top of stack is zero            |
| LL      | Label      | Jump to a label if the top of stack is negative        |
| LT      | -          | End a subroutine and transfer control back to the call |
| TT      | -          | End the program                                        |

### I/O (IMP:[LT])

| Command | Parameters | Meaning                                                                     |
| ------- | ---------- | --------------------------------------------------------------------------- |
| FF      | -          | Output the character at the top of the stack                                |
| FL      | -          | Output the number at the top of the stack                                   |
| LF      | -          | Read a character and place it in the location given by the top of the stack |
| LL      | -          | Read a number and place it in the location given by the top of the stack    |

### Number and Label

#### Number

|         | Sign | Bits F = 0, L = 1 | Terminated by a T |     |
| ------- | ---- | ----------------- | ----------------- | --- |
| Example | F    | LF                | T                 | = 2 |

#### Label

Label is case sensitive

|         | Label | Terminated by a T |
| ------- | ----- | ----------------- |
| Example | FL    | T                 |

### Example

|        |                                  |
| ------ | -------------------------------- |
| FFFLFT | Push 1 onto the stack            |
| TFFLT  | Set a label "L" at this position |
| TFLLT  | Jump to label "L"                |
