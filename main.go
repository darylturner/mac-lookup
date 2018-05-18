package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "log"
        "net/http"
        "os"
)

type MacResult struct {
        Result struct {
                Company  string `json:"company"`
                Prefix   string `json:"mac_prefix"`
                Address  string `json:"address"`
                StartHex string `json:"start_hex"`
                EndHex   string `json:"end_hex"`
                Country  string `json:"country"`
                Type     string `json:"type"`
        } `json:"result"`
}

func lookup(mac string) (macDetail MacResult, err error) {
        res, err := http.Get("https://macvendors.co/api/" + mac)
        if err != nil {
                return macDetail, fmt.Errorf("error during macvendors.co api lookup: %s", err)
        }
        defer res.Body.Close()

        response := json.NewDecoder(res.Body)
        if err = response.Decode(&macDetail); err != nil {
                return
        }
        return
}

func main() {
        version := flag.Bool("v", false, "version info")
        brief := flag.Bool("b", false, "brief output. company name only")
        flag.Parse()

        const usage string = `cli wrapper to http://macvendors.co
usage: mac-lookup [-b] <mac address>
       mac-lookup -v`

        if *version || flag.NArg() != 1 {
                fmt.Println(usage)
                os.Exit(0)
        }

        mac := flag.Arg(0)
        macDetail, err := lookup(mac)
        if err != nil {
                log.Fatal("error during macvendors.co lookup:", err)
        }

        if *brief {
                fmt.Println(macDetail.Result.Company)
        } else {
                stdout := json.NewEncoder(os.Stdout)
                stdout.SetIndent("", " ")

                if err := stdout.Encode(macDetail); err != nil {
                        log.Fatal("error writing to stdout:", err)
                }
        }
}
