package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{
        Use: "fastanno",
        Version: "1.0.0",
    }

    indexCmd.Flags().StringP("file", "f", "", "input file")
    indexCmd.Flags().IntP("bin", "b", 1000, "bin size")

    annoCmd.Flags().StringP("query", "q", "", "query file")
    annoCmd.Flags().StringP("dbs", "d", "", "one or more dbs file, join by comma")
    annoCmd.Flags().StringP("keycols", "k", "1,2,3,4,5", "key columns, ex: 1,2,4,5 for vcf; 1,2,3,4,5 for annovar tsv")
    annoCmd.Flags().StringP("out", "o", "", "output file")
    annoCmd.Flags().IntP("threads", "t", 1, "set threads")

    rootCmd.AddCommand(annoCmd)
    rootCmd.AddCommand(indexCmd)


    rootCmd.CompletionOptions.DisableDefaultCmd = true
    rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}