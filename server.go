package main

import (
    "encoding/json"
    "net/http"
    "sync"
    "os"
    "log"
    "io/ioutil"
    "strings"
)


// Config schema
type Config struct {
    Name string `json:"name,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// in memory database of configs mapped with a string Name
// sync.Mutex for locking and unlocking the database while writing and reading it.
type configHandler struct {
    sync.Mutex
    database map[string]Config
}

// constructor for configHandler which returns a initialised database
func newConfigHandler() *configHandler {
    return &configHandler {
        database: map[string]Config{},
    }
}

// Switch case logic for routing multiple methods for the same route /configs
func (c *configHandler) configsMethods(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "GET":
            c.getAllConfigs(w, r)
            return
        case "POST":
            c.createNewConfig(w, r)
            return
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            w.Write([]byte("Method not allowed"))
            log.Print("Method not allowed")
            return
    }
}

// Switch case logic for routing multiple methods for the same route /configs/{name}
func (c *configHandler) singleConfigMethods(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "GET":
            c.getConfig(w, r)
            return
        case "PUT", "PATCH":
            c.updateConfig(w, r)
            return
        case "DELETE":
            c.deleteConfig(w, r)
            return
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            w.Write([]byte("Method not allowed"))
            log.Print("Method not allowed")
            return

    }
}


// Query database and respond with matching config
// URL: /search?metadata.key=value
// Method: GET
// Output: json of configs matching the query
func (c *configHandler) queryDatabase(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
        return
    }
    params := strings.Split(strings.Split(r.URL.String(), "?")[1], ".")
    if params[0] != "metadata" {
        log.Print("wrong query")
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("bad Request. Please query metadata"))
        return
    }
    lastparam := strings.Split(params[len(params)-1], "=")[0]
    searchvalue := strings.Split(params[len(params)-1], "=")[1]

    // list of configs that matches the query
    reqConfigs := []Config{}

    c.Lock()
    for _, config := range c.database {
        meta := config.Metadata
        var parse func(metadata map[string]interface{})
        parse = func(metadata map[string]interface{}) {
            for key, val := range metadata {
                // switch case for matching data type
                switch actualVal := val.(type) {
                    case map[string]interface{}:
                        if key == params[len(params)-2] {
                            if actualVal[lastparam] == searchvalue {
                                reqConfigs = append(reqConfigs, config)
                                log.Print(reqConfigs)
                                break
                            }
                        }
                        parse(val.(map[string]interface{}))
                    default:
                        if len(params) == 2 {
                            if key == lastparam {
                                if actualVal == searchvalue {
                                    reqConfigs = append(reqConfigs, config)
                                    log.Print(reqConfigs)
                                    break
                                }
                            }
                        }
                }
            }
        }
        parse(meta)
    }
    c.Unlock()

    responsejson, err := json.Marshal(reqConfigs)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responsejson)
}

// Fetch all exisitng configs and respond
// URL: /configs
// Method: GET
// Output: json of all configs in the database
func (c *configHandler) getAllConfigs(w http.ResponseWriter, r *http.Request) {
    log.Printf("received /configs GET request from %v", r.Host)
    configs := make([]Config, len(c.database))

    c.Lock()
    i :=0
    for _, config := range c.database {
        configs[i] = config
        i++
    }
    c.Unlock()

    responsejson, err := json.Marshal(configs)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }

    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responsejson)
}

//Fetch the config that needs to be displayed, updated or deleted
func (c *configHandler) filterConfig(w http.ResponseWriter, r *http.Request) *Config {
    urllen := strings.Split(r.URL.String(), "/")
    if len(urllen) !=3 {
        w.WriteHeader(http.StatusNotFound)
        log.Print("Bad request")
    }

    c.Lock()
    reqconfig, ok := c.database[urllen[2]]
    c.Unlock()
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        log.Print("Failed fetching the desired config")
    }

    return &reqconfig
}

// Fetch the desired config and respond
// URL: /configs/{name}
// Method: GET
// Output: json of requested config
func (c *configHandler) getConfig(w http.ResponseWriter, r *http.Request) {
    reqconfig := c.filterConfig(w, r)
    log.Printf("received /configs/%v GET request from %v", reqconfig.Name, r.Host)

    responsejson, err := json.Marshal(reqconfig)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }

    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responsejson)
}

// Update the desire config and respond
// URL: /configs/{name}
// Method: PUT/PATCH
// Output: updates the existing config with new config
func (c *configHandler) updateConfig(w http.ResponseWriter, r *http.Request) {
    reqconfig := c.filterConfig(w, r)
    log.Printf("received /configs/%v PUT/PATCH request from %v", reqconfig.Name, r.Host)

    requestbody, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    var updatedconfig Config
    err = json.Unmarshal(requestbody, &updatedconfig)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    c.Lock()
    c.database[reqconfig.Name] = updatedconfig
    c.Unlock()

    log.Printf("config %v updated successfully", reqconfig.Metadata)
    w.WriteHeader(http.StatusOK)
}


// Delete the desire config and respond
// URL: /configs/{name}
// Method: DELETE
// Output: deletes the desired config from database
func (c *configHandler) deleteConfig(w http.ResponseWriter, r *http.Request) {
    reqconfig := c.filterConfig(w, r)
    log.Printf("received /configs/%v DELETE request from %v", reqconfig.Name, r.Host)
    delete(c.database, reqconfig.Name)
    log.Printf("deleted %v successfully from database", reqconfig.Name)
    w.WriteHeader(http.StatusOK)
}

// Create a new config from requestbody
// URL: /configs
// Method: POST
// Output: create a new config
func (c *configHandler) createNewConfig(w http.ResponseWriter, r *http.Request) {
    log.Printf("received /configs POST request from %v", r.Host)
    requestbody, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    // check request's content-type
    content := r.Header.Get("content-type")
    if content != "application/json" {
        w.WriteHeader(http.StatusUnsupportedMediaType)
        w.Write([]byte(err.Error()))
        return
    }

    var newconfig Config
    err = json.Unmarshal(requestbody, &newconfig)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    c.Lock()
    c.database[newconfig.Name] = newconfig
    c.Unlock()
    log.Print("Created new config successfully and added to database ", c.database)

    w.WriteHeader(http.StatusOK)
}

func main() {
    // Letting user run the server on desired port by setting an environment variable SERVER_PORT
    server_port := os.Getenv("SERVE_PORT")
    if server_port == "" {
        panic("SERVE_PORT env variable is required to run server on the desired port.")
    }

    confighandle := newConfigHandler()
    http.HandleFunc("/configs", confighandle.configsMethods)
    http.HandleFunc("/configs/", confighandle.singleConfigMethods)
    http.HandleFunc("/search", confighandle.queryDatabase)
    err := http.ListenAndServe(":"+string(server_port), nil)
    if err != nil {
        log.Fatal(err)
    }
}
