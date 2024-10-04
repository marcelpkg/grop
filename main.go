package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var Reset = "\033[0m" 
var Red = "\033[31m" 
var Green = "\033[32m" 
var Yellow = "\033[33m" 
var Blue = "\033[34m" 
var Magenta = "\033[35m" 
var Cyan = "\033[36m" 
var Gray = "\033[37m" 
var White = "\033[97m"

func skipFile(ext string) bool {
    skip := map[string]bool{
        ".class": true,
		".exe": true,
		".bin": true,
		".dll": true,
		".so": true,
		".o": true,
        ".jar": true,
        ".tar": true,
        ".gz": true,
        ".zip": true,
        "": true,
    }

    return skip[ext]
}

func searchDirectory(dirpath string, target string, resultChan chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    files, err := os.ReadDir(dirpath)
	if err != nil {
		fmt.Println(err)
        return
	}

	for _, file := range files {
        path := filepath.Join(dirpath, file.Name())

        if file.IsDir() {
            wg.Add(1)
            go searchDirectory(path, target,  resultChan, wg)
        } else {
            if !skipFile(filepath.Ext(path)) {
                wg.Add(1)
                go searchFile(path, target, resultChan, wg)
            }
        } 
    }
}

func searchFile(path string, target string, resultChan chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    file, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        
        if strings.Contains(line, target) {
            resultChan <- strings.ReplaceAll(fmt.Sprintln(Magenta + path + ": " + strings.TrimSpace(line)), target, Red + target + Reset)
        }
    }
}

func main() {
    wg := sync.WaitGroup{}
    args := os.Args

    if len(args) < 3 {
        log.Fatal("usage: grop <string> <directory>")
        os.Exit(3)
    }

    target := args[1]
    rootDir := args[2]

    resultChan := make(chan string)

    go func() {
        for res := range resultChan {
            fmt.Print(res)
        }
    }()

    files, err := os.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
        path := filepath.Join(rootDir, file.Name())

        if file.IsDir() {
            wg.Add(1)
            go searchDirectory(path, target, resultChan, &wg)
        } else {
            wg.Add(1)
            go searchFile(path, target, resultChan, &wg)
        } 
    }
    
    wg.Wait()
    close(resultChan)
}