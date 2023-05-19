package lzw

import (
    "os"
    "fmt"
)

func Zip(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    content := make([]byte, 1024)
    file.Read(content)
    fmt.Println(content)
    
    return nil
}