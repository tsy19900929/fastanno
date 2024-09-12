package main

import (
    "bufio"
    "fmt"
    "github.com/spf13/cobra"
    "os"
    "strconv"
    "sort"
    "strings"
    "github.com/facette/natsort"
)

var indexCmd = &cobra.Command{
    Use:   "index --file [inputFile] --bin [binSize]",
    Short: "Create index, input file must already be sorted",
    Run: func(cmd *cobra.Command, args []string) {
        inputFile, _ := cmd.Flags().GetString("file")
        binSize, _ := cmd.Flags().GetInt("bin")

        file, err := os.Open(inputFile)
        if err != nil {
            fmt.Println(inputFile, "file not found.")
            return
        }
        defer file.Close()

        index := make(map[string]map[int][]int)
        var previousFilePosition, currentFilePosition int = 0, 0

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := scanner.Text()
            currentFilePosition += len(line) + 1

            if strings.HasPrefix(line, "#") {
                previousFilePosition = currentFilePosition
                continue
            }

            fields := strings.Split(line, "\t")
            chr := fields[0]
            start, _ := strconv.Atoi(fields[1])

            binStart := (start / binSize) * binSize

            if _, exists := index[chr]; !exists {
                index[chr] = make(map[int][]int)
            }
            if _, exists := index[chr][binStart]; !exists {
                index[chr][binStart] = []int{previousFilePosition, currentFilePosition}
            } else {
                index[chr][binStart][1] = currentFilePosition
            }

            previousFilePosition = currentFilePosition
        }

        indexFile := inputFile + ".idx"
        outFile, err := os.Create(indexFile)
        if err != nil {
            fmt.Println("error creating index file:", err)
            return
        }
        defer outFile.Close()

        fmt.Fprintf(outFile, "#BIN\t%d\t%d\n", binSize, getFileSize(inputFile))

        var chrs []string
        for chr := range index {
            chrs = append(chrs, chr)
        }
        natsort.Sort(chrs)

        for _, chr := range chrs {
            regions := index[chr]
            var keys []int
            for key := range regions {
                keys = append(keys, key)
            }
            sort.Ints(keys)
            for _, binStart := range keys {
                start := regions[binStart][0]
                stop := regions[binStart][1]
                fmt.Fprintf(outFile, "%s\t%d\t%d\t%d\n", chr, binStart, start, stop)
            }
        }
    },
}

func getFileSize(filename string) int64 {
    fileInfo, err := os.Stat(filename)
    if err != nil {
        return 0
    }
    return fileInfo.Size()
}