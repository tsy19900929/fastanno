package utils

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

func SplitFile(filename string, parts int) ([]string, error) {
    var filenames []string
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println(filename, "not found.")
        os.Exit(1)
    }
    defer file.Close()
    scanner1 := bufio.NewScanner(file)
    totalLines := 0
    for scanner1.Scan() {
        totalLines++
    }

    linesPerFile := (totalLines + parts - 1) / parts

    file.Seek(0, 0)
    scanner2 := bufio.NewScanner(file)
    for i := 0; i < parts; i++ {
        partFilename := ".fatmp@" + filename + "@p" + strconv.Itoa(i)
        partFile, err := os.Create(partFilename)
        if err != nil {
            return nil, err
        }
        defer partFile.Close()

        for j := 0; j < linesPerFile; j++ {
            if !scanner2.Scan() {
                break
            }
            line := scanner2.Text()
            fmt.Fprintln(partFile, line)
        }
        filenames = append(filenames, partFilename)
    }
    return filenames, nil
}

func MergeFiles(fileNames []string, outputFile string) error {
	out, _ := os.Create(outputFile)
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for _, filename := range fileNames {
		file, _ := os.Open(filename)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			writer.WriteString(scanner.Text() + "\n")
		}
		file.Close()
	}
	return nil
}
