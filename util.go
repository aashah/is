package main

import "os"

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