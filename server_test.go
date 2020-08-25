package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "fmt"
    "bytes"
)


func TestGetAllConfigs(t *testing.T) {
    req, err := http.NewRequest("GET", "/configs", nil)
    if err != nil {
        t.Fatal(err)
    }
    recorder := httptest.NewRecorder()
    confighandle := newConfigHandler()
    handler := http.HandlerFunc(confighandle.configsMethods)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusOK)
    }
    fmt.Printf(recorder.Body.String())
    if body := recorder.Body.String(); body != "[]" {
        t.Errorf("Handler returned wrong body, got %v want []", body)
    }
}

func TestCreateNewConfig(t *testing.T) {
    var jsonStr = []byte(`{"name":"test", "metadata":{"burger":{"calories": "230"}}}`)

    req, err := http.NewRequest("POST", "/configs", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }
    confighandle := newConfigHandler()
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(confighandle.configsMethods)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusOK)
    }
}

func TestGetConfig(t *testing.T) {
    req, err := http.NewRequest("GET", "/configs/test", nil)
    if err != nil {
        t.Fatal(err)
    }
    recorder := httptest.NewRecorder()
    confighandle := newConfigHandler()
    handler := http.HandlerFunc(confighandle.singleConfigMethods)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusNotFound)
    }
}

func TestUpdateConfig(t *testing.T) {
    var jsonStr = []byte(`{"name":"datacenter-1","metadata":{"limits":{"cpu":"1","mem":"500m"}}}`)

    req, err:= http.NewRequest("PUT", "/configs/datacenter-1", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }

    confighandle := newConfigHandler()
    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(confighandle.singleConfigMethods)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusNotFound)
    }
}

func TestPatchConfig(t *testing.T) {
    var jsonStr = []byte(`{"name":"datacenter-2","metadata":{"limits":{"cpu":"1","mem":"500m"}}}`)

    req, err:= http.NewRequest("PATCH", "/configs/datacenter-2", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }

    confighandle := newConfigHandler()
    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(confighandle.singleConfigMethods)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusNotFound)
    }
}

func TestDeleteConfig(t *testing.T) {
    req, err := http.NewRequest("DELETE", "/configs/test", nil)
    if err != nil {
        t.Fatal(err)
    }

    confighandle := newConfigHandler()
    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(confighandle.singleConfigMethods)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusNotFound)
    }
}

func TestQueryConfig(t *testing.T) {
    req, err := http.NewRequest("GET", "/search?metadata.limits.cpu=1", nil)
    if err != nil {
        t.Fatal(err)
    }

    confighandle := newConfigHandler()
    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(confighandle.queryDatabase)
    handler.ServeHTTP(recorder, req)
    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code, got %v want %v", status, http.StatusOK)
    }
}
