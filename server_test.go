package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "fmt"
)

func TestgetAllConfigs(t *testing.T) {
    req, err := http.NewRequest("GET", "/configs", nil)
    if err != nil {
        t.Fatal(err)
    }
    recorder := httptest.NewRecorder()
    confighandle := newConfigHandler()
    handler := http.HandlerFunc(confighandle.configsOptions)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusOK)
    }
    fmt.Printf(recorder.Body.String())
    if body := recorder.Body.String(); body != "[{'name': 'hari'}]" {
        t.Errorf("Handler returned wrong body, got %v want []", body)
    }
}
