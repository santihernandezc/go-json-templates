package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
)

type reqBody struct {
	Template string         `json:"template"`
	Data     map[string]any `json:"data"`
}

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var req reqBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("New request\tmethod: %q\ttemplate: %q\tdata: %v\n", r.Method, string(req.Template), req.Data)

		// TODO: string?
		scanner := newScanner([]byte(req.Template))
		tokens := scanner.scan()

		parser := newParser(tokens)
		statements, err := parser.parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{"err:": err.Error()}); err != nil {
				panic(err)
			}
			return
		}

		interpreter := newInterpreter(statements, req.Data)
		res, err := interpreter.interpret()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{"err": err.Error()}); err != nil {
				panic(err)
			}
			return
		}

		if _, err := w.Write(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))

	http.ListenAndServe(":8080", nil)
}
