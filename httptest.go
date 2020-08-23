package main

import (
    "io"
    "net/http"
)
func queryParamDisplayHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, "name: "+req.FormValue("name"))
    io.WriteString(res, "\nphone: "+req.FormValue("phone"))
}

func main() {
    http.HandleFunc("/example", func(res http.ResponseWriter, req *http.Request) {
        queryParamDisplayHandler(res, req)
    })
    println("Enter this in your browser:  http://localhost:8080/example?name=jenny&phone=867-5309")
    http.ListenAndServe(":8080", nil)
}
