package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		if err := startServer(":8080"); err != nil {
			fmt.Fprintln(os.Stderr, "server error:", err)
			os.Exit(1)
		}
		return
	}

	output, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}
	fmt.Print(output)
}

func run(source string) (string, error) {
	bytes := make([]byte, 2048)
	ptr := 0
	var output strings.Builder

	for i := 0; i < len(source); i++ {
		switch source[i] {
		case '>':
			ptr++
			if ptr >= len(bytes) {
				ptr = 0
			}
		case '<':
			ptr--
			if ptr < 0 {
				ptr = len(bytes) - 1
			}
		case '+':
			bytes[ptr]++
		case '-':
			bytes[ptr]--
		case '.':
			output.WriteByte(bytes[ptr])
		case '[':
			if bytes[ptr] == 0 {
				nest := 1
				for nest > 0 {
					i++
					if i >= len(source) {
						return "", fmt.Errorf("unmatched '[' at position %d", i)
					}
					if source[i] == '[' {
						nest++
					} else if source[i] == ']' {
						nest--
					}
				}
			}
		case ']':
			if bytes[ptr] != 0 {
				nest := 1
				for nest > 0 {
					i--
					if i < 0 {
						return "", fmt.Errorf("unmatched ']' at position %d", i)
					}
					if source[i] == ']' {
						nest++
					} else if source[i] == '[' {
						nest--
					}
				}
			}
		}
	}

	return output.String(), nil
}

func startServer(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/run", runHandler)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	fmt.Println("Brainfuck UI available at http://localhost" + addr)
	return http.ListenAndServe(addr, mux)
}

type runRequest struct {
	Code string `json:"code"`
}

type runResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, runResponse{Error: "invalid JSON payload"})
		return
	}

	if len(req.Code) == 0 {
		writeJSON(w, http.StatusBadRequest, runResponse{Error: "brainfuck code is required"})
		return
	}

	output, err := run(req.Code)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, runResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, runResponse{Output: output})
}

func writeJSON(w http.ResponseWriter, status int, payload runResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
