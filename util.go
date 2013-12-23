package main

import (
    "fmt"
    "bytes"
    "os"
    "os/exec"
    "strings"
)

func readFile(path string) (data []byte, err error) {
    var fi *os.File
    var fiStat os.FileInfo

    if fi, err = os.Open(path); err != nil {
        return nil, err
    }
    defer fi.Close()

    if fiStat, err = fi.Stat(); err != nil {
        return nil, err
    }

    var raw []byte

    raw = make([]byte, fiStat.Size())

    if _, err = fi.Read(raw); err != nil {
        return nil, err
    }

    return raw, nil
}

func expand(list map[string]string, str string) string {
    ret := str
    for key, val := range list {
        ret = strings.Replace(ret, "{"+key+"}", val, -1)
    }
    return ret
}

func (v *vcsCmd) runCmd(dir string, cmdLine string, verbose bool, keyvals map[string]string) error {
    if keyvals != nil {
        // expand cmd
        cmdLine = expand(keyvals, cmdLine)
    }

    args := strings.Fields(cmdLine)

    // TODO Check if v.cmd exists
    if _, err := exec.LookPath(v.cmd); err != nil {
        fmt.Fprintf(os.Stderr, "is: missing %s command.", v.cmd)
        return err
    }

    if verbose {
        fmt.Println("Executing", v.cmd, "with", cmdLine, "...")
    }
    // Execute
    cmd := exec.Command(v.cmd, args...)
    cmd.Dir = dir
    var buf bytes.Buffer
    cmd.Stdout = &buf
    cmd.Stderr = &buf
    err := cmd.Run()
    out := buf.Bytes()
    if verbose {
        os.Stdout.Write(out)
    }
    return err
}