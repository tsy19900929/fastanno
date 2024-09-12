package main

import (
    "bufio"
    "fmt"
    "github.com/spf13/cobra"
    "os"
    "strconv"
    "strings"
    "fastanno/utils"
    "time"
    "sync"
)

func openFile(fileName string) (*os.File, error) {
    file, err := os.Open(fileName)
    if err != nil {
        fmt.Println(fileName, "not found.")
        os.Exit(1)
    }
    return file, err
}

func processLine(line string, colsInt []int)([]string, string){
    if strings.HasPrefix(line, "#") {
        return nil, ""
    }
    fields := strings.Split(line, "\t")
    keys := ""
    for _, col := range colsInt {
        if len(fields) < col {
            fmt.Println("weird format, separator? :", line)
            os.Exit(1)
        }
        keys = keys + "\t" + fields[col - 1]
    }
    return fields, keys
}

func fillWithDot(colsMax int, fields []string) string {
    parts := make([]string, len(fields) - colsMax)
    for i := range parts {
        parts[i] = "."
    }
    return strings.Join(parts, "\t")
}

func processPartFile(queryName string, dbsName string, colsInt []int, colsMax int, logFile *os.File, wg *sync.WaitGroup) {
    muti := strings.HasPrefix(queryName, ".fatmp@")
    if muti {
        defer wg.Done()
    }

    fileA, _ := openFile(queryName)
    defer fileA.Close()

    queryBodyFileName := queryName + "@body"
    if !muti {
        queryBodyFileName = ".fatmp@" + queryBodyFileName
    }

    queryBodyFile, _ := os.Create(queryBodyFileName)
    defer queryBodyFile.Close()
    FilesToPaste := []string{queryBodyFileName}

    dbsNames := strings.Split(dbsName, ",")
    dbIndex := true
    var binSize int
    var indexMap map[string]map[int][]int
    onceA := true
    for _, dbName := range dbsNames {
        now := time.Now()
        indexFile, err := os.Open(dbName + ".idx")
        if err != nil {
            dbIndex = false
            fmt.Fprintln(logFile, dbName + ".idx", "file not found, skip it.")
        }
        defer indexFile.Close()

        if dbIndex {
            indexMap = make(map[string]map[int][]int)
            scanner := bufio.NewScanner(indexFile)
            for scanner.Scan() {
                line := scanner.Text()
                fields := strings.Split(line, "\t")
                chr := fields[0]
                if chr == "#BIN" {
                    binSize, _ = strconv.Atoi(fields[1])
                    continue
                }
                binStart, _ := strconv.Atoi(fields[1])
                start, _ := strconv.Atoi(fields[2])
                end, _ := strconv.Atoi(fields[3])
                if indexMap[chr] == nil {
                    indexMap[chr] = make(map[int][]int)
                }
                indexMap[chr][binStart] = []int{start, end}
            }
        }

        fileB, err := openFile(dbName)
        if err != nil {
            fmt.Println(dbName, "not found.")
            os.Exit(1)
        }
        defer fileB.Close()

        tmpAnnoFileName := queryBodyFileName + "@" + dbName
        FilesToPaste = append(FilesToPaste, tmpAnnoFileName)
        tmpAnnoFile, _ := os.Create(tmpAnnoFileName)
        defer tmpAnnoFile.Close()

        onceBCols := true
        dot := ""
        fileA.Seek(int64(0), 0)
        scannerA := bufio.NewScanner(fileA)
        for scannerA.Scan() {
            lineA := scannerA.Text()
            fieldsA, keysA := processLine(lineA, colsInt)
            if fieldsA == nil {
                continue
            }

            if onceA {
                fmt.Fprintln(queryBodyFile, lineA)
            }

            var offset, indexEnd int
            if dbIndex {
                chr := fieldsA[0]
                start, _ := strconv.Atoi(fieldsA[1])
                binStart := (start / binSize) * binSize
                indexStart := indexMap[chr][binStart][0]
                indexEnd = indexMap[chr][binStart][1]
                offset = indexStart
            }

            if dbIndex {
                fileB.Seek(int64(offset), 0)
            } else {
                fileB.Seek(int64(0), 0)
            }

            anno := ""
            if dot != "" {
                anno = dot
            }

            scannerB := bufio.NewScanner(fileB)
            for scannerB.Scan() {
                //shit，altered by humans
                lineB := scannerB.Text()
                offset += len(lineB) + 1
                fieldsB, keysB := processLine(lineB, colsInt)
                if fieldsB == nil {
                    continue
                }

                if onceBCols {
                    dot = fillWithDot(colsMax, fieldsB)
                    anno = dot
                    onceBCols = false
                }
                if keysA == keysB {
                    anno = strings.Join(fieldsB[colsMax:], "\t")
                    break
                }
                if dbIndex {
                    if offset > indexEnd {
                        break
                    }
                }
            }
            fmt.Fprintln(tmpAnnoFile, anno)

        }
        fmt.Fprintln(logFile, tmpAnnoFileName, "time cost:", time.Since(now))
        if onceA {
            onceA = false
        }
    }
    utils.PasteFiles(FilesToPaste, queryBodyFileName + "@all")
}


var annoCmd = &cobra.Command{
    Use:   "anno --query [queryFile] --dbs [dbsFile] --keycols [keyColumns] --out [outputFile] --threads [threads]",
    Short: "Annotate query file with dbs file",
    Run: func(cmd *cobra.Command, args []string) {

        queryName, _ := cmd.Flags().GetString("query")
        dbsName, _ := cmd.Flags().GetString("dbs")
        keycols, _ := cmd.Flags().GetString("keycols")
        outName, _ := cmd.Flags().GetString("out")
        threads, _ := cmd.Flags().GetInt("threads")

        cols := strings.Split(keycols, ",")
        colsInt := make([]int, len(cols))
        colsMax := 0
        for i, col := range cols {
            colsInt[i], _ = strconv.Atoi(col)
            if colsInt[i] > colsMax {
                colsMax = colsInt[i]
            }
        }

        logFile, _ := os.Create("anno.log")
        defer logFile.Close()

        if threads == 1 {
            processPartFile(queryName, dbsName, colsInt, colsMax, logFile, nil)
            os.Rename(".fatmp@" + queryName + "@body@all", outName)
        } else {
            partFilesNames, _ := utils.SplitFile(queryName, threads)
            partFilesAllNames := []string{}
            var wg sync.WaitGroup
            for _, partFileName := range partFilesNames {
                partFilesAllNames = append(partFilesAllNames, partFileName + "@body@all")
                wg.Add(1)
                go processPartFile(partFileName, dbsName, colsInt, colsMax, logFile, &wg)
            }
            wg.Wait()

            fmt.Fprintln(logFile, "partFilesAllNames:", strings.Join(partFilesAllNames, "\t"))
            utils.MergeFiles(partFilesAllNames, outName)
        }

        files, _ := os.ReadDir(".")
        for _, file := range files {
            if strings.HasPrefix(file.Name(), ".fatmp@") {
                os.Remove(file.Name())
            }
        }
    },
}