package main

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
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

func runCmd(dir string, cmdName string, cmdLine string, verbose bool, keyvals map[string]string) error {
    if keyvals != nil {
        // expand cmd
        cmdLine = expand(keyvals, cmdLine)
    }

    args := strings.Fields(cmdLine)

    // TODO Check if v.cmd exists
    if _, err := exec.LookPath(cmdName); err != nil {
        fmt.Fprintf(os.Stderr, "is: missing %s command.", cmdName)
        return err
    }

    if verbose {
        fmt.Println("Executing", cmdName, "with", cmdLine, "...")
    }
    // Execute
    cmd := exec.Command(cmdName, args...)
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

func copyFile(src string, dstDirectory string) (err error) {
    var sf, df *os.File
    var sfStat os.FileInfo

    if sf, err = os.Open(src); err != nil {
        return err
    }
    defer sf.Close()

    var dst string
    if sfStat, err = sf.Stat(); err != nil {
        return err
    }
    dst = filepath.Join(dstDirectory, sfStat.Name())


    if df, err = os.Create(dst); err != nil {
        return err
    }
    defer df.Close()

    _, err = io.Copy(sf, df)
    return
}