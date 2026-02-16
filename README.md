## 🧠 Brainfuck Interpreter

This project implements an interpreter for **Brainfuck**, a famous **esoteric programming language (esolang)**. Created by Urban Müller in 1993, Brainfuck is notable for its extreme minimalism. Despite having only **eight simple commands**, it is **Turing complete**, serving as a classic intellectual puzzle for programmers.

---

### ✨ Implementation Details

The interpreter operates based on the core components defined by the Brainfuck specification:

* **Memory Array:** A simple array of **2048 bytes**, with all values initialized to **0**. This array serves as the program's primary memory.
* **Data Pointer:** A single pointer that starts at the first byte (index 0) and is used to select the currently active cell in the memory array.
* **Source Code:** The program accepts the Brainfuck source code as a string input, processing it character by character. The input is assumed to be valid and contains fewer than 4096 operations.



---

### 🕹️ The Eight Brainfuck Commands

The language consists of eight single-character operators that manipulate the data pointer and the value of the pointed-to byte.

| Command | Name | Action |
| :------ | :--- | :--- |
| **`>`** | Pointer Increment | Move the data pointer one cell to the **right**. |
| **`<`** | Pointer Decrement | Move the data pointer one cell to the **left**. |
| **`+`** | Value Increment | **Increment** the value of the byte at the current pointer location. |
| **`-`** | Value Decrement | **Decrement** the value of the byte at the current pointer location. |
| **`.`** | Output | **Output** the value of the pointed byte as an **ASCII character** to standard output. |
| **`,`** | Input | **Input** a single byte and store its value in the pointed cell (standard Brainfuck, often omitted in simple interpreters). |
| **`[`** | Loop Start | If the value of the pointed byte is **0**, skip execution forward to the command after the matching `]`. |
| **`]`** | Loop End | If the value of the pointed byte is **non-zero**, jump execution back to the command after the matching `[`. |

> **Note:** Any character in the source code that is not one of the eight defined commands is treated as a **comment** and is ignored.

---

### 🚀 Usage

The Brainfuck source code is provided as the first command-line argument when executing the interpreter (assuming a Go project structure for the examples).

**Example 1: "Hello World!"**

```console
$ go run . "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>." | cat -e
Hello World!$
```
```console
$ go run . "+++++[>++++[>++++H>+++++i<<-]>>>++\n<<<<-]>>--------.>+++++.>." | cat -e
Hi$
```
```console
$ go run . "++++++++++[>++++++++++>++++++++++>++++++++++<<<-]>---.>--.>-.>++++++++++." | cat -e
abc$
```
```console
$ go run .
$
```

### Web UI Mode

If you run the project without arguments, it now starts a web server and serves a frontend:

```console
$ go run .
Brainfuck UI available at http://localhost:8080
```

Open `http://localhost:8080` in your browser, paste Brainfuck code, and click **Run Program**.

The frontend calls:

- `POST /api/run` with JSON body: `{"code":"<brainfuck program>"}`
- JSON response: `{"output":"..."}` or `{"error":"..."}`
