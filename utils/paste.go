package utils

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func PasteFiles(fileNames []string, outputFile string) error {
    scanners := make([]*bufio.Scanner, len(fileNames))

    for i, fileName := range fileNames {
        file, err := os.Open(fileName)
        if err != nil {
            return fmt.Errorf("error opening file %s: %v", fileName, err)
        }
        defer file.Close()
        scanners[i] = bufio.NewScanner(file)
    }

    var output *os.File
    if outputFile == "" {
        output = os.Stdout
    } else {
        // shit, not :=
        output, _ = os.Create(outputFile)
        defer output.Close()
    }

    for {
        allDone := true
        var lines []string
        lines = lines[:0]
        for _, scanner := range scanners {
            if scanner.Scan() {
                lines = append(lines, scanner.Text())
                allDone = false
            } else {
                lines = append(lines, "")
            }
        }
        if allDone {
            break
        }
        fmt.Fprintln(output, joinWithTabs(lines))
    }
    return nil
}

func joinWithTabs(lines []string) string {
    return strings.Join(lines, "\t")
}