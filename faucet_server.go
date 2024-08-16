package main  

import (  
    "fmt"  
    "io/ioutil"  
    "net/http"  
    "os/exec" 
)  

var globalData string

func execute(params []string) bool {  
    if len(params) == 0 {  
        fmt.Println("No command provided")  
        return false
    }

    cmd := params[0]  
    args := params[1:]  

    out, err := exec.Command(cmd, args...).Output()  

    if err != nil {  
        fmt.Printf("%s\n", err)  
        return false
    }  

    fmt.Println("Command Successfully Executed")  
    output := string(out[:])  
    fmt.Println(output)
	return true
}  

// CORS middleware  
func withCORS(handler http.HandlerFunc) http.HandlerFunc {  
    return func(w http.ResponseWriter, r *http.Request) {  
        w.Header().Set("Access-Control-Allow-Origin", "*")  
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")  
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  
		w.Header().Set("Content-Type", "text/plain") 
        if r.Method == http.MethodOptions {  
            w.WriteHeader(http.StatusOK)
            return
        }
        handler(w, r)
    }
}

func handleConnection(writer http.ResponseWriter, request *http.Request) {
    if request.Method != "POST" {
        http.Error(writer, "Only POST method allowed", http.StatusMethodNotAllowed)
        return
    }

    // Read the request body
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Error reading request body", http.StatusInternalServerError)
        return
    }

	globalData = string(body)

    // Print the received data to the console
    fmt.Printf("Received data: %s\n", globalData)

    if globalData != "evmos17u6aw9l89myt7mmfr3vfluzkst4w7ths0sa9ru" {
        // Respond with a success message
        executeState := execute([]string{"evmosd", "tx", "bank", "send", "evmos1tt5kdszrefc6gy3ty535h8aj23jh7h3l9ymhtl", globalData, "2000000000000aevmos", "--fees", "1400000aevmos", "-y"})
        // execute([]string{"evmosd", "q", "bank", "balances", globalData})

        response := "Data received successfully!"
        if executeState == false {
            response = "Something went wrong. Please try again later."
        }
        writer.Write([]byte(response))
    }
}

func main() {
    http.HandleFunc("/", withCORS(handleConnection))

    fmt.Println("Server listening on port 2000...")
    err := http.ListenAndServe(":2000", nil)
    if err != nil {
        fmt.Println("Error starting HTTP server:", err)
    }
}