package main

import (
    "github.com/jackpal/bencode-go"
    "fmt"
    "os"
    "io/ioutil"
    "strings"
    "path/filepath"
    "time"
)

// Nested structs with the metadata we need
type FileDict struct {
    Path   []string "path"
}

type InfoDict struct {
    Name        string     "name"
    Files       []FileDict "files"
}

type MetaInfo struct {
    Info         InfoDict "info"
    CreationDate int64    "creation date"
    Comment      string   "comment"
    CreatedBy    string   "created by"
}

// Comodity panic
func check(e error) {
    if e != nil {
        fmt.Println("ERROR:", e)
        os.Exit(1)
    }
}

// Check for empty dirs
func IsEmptyDir(name string) (bool, error) {
    entries, err := ioutil.ReadDir(name)
    if err != nil {
        return false, err
    }
    return len(entries) == 0, nil
}

// Recursively delete orphaned dirs
func recRmOrphanedDir(name string) (error) {
    if empty, _ := IsEmptyDir(name); empty {
        fmt.Printf("Deleting empty dir '%s'\n", name)
        e := os.Remove(name)
        if e != nil {
            return e
        }
        previous := filepath.Dir(name)
        return recRmOrphanedDir(previous)
    }
    return nil
}

func main() {
    // Init vars
    progName := filepath.Base(os.Args[0])
    torrent, searchDir, backupDir := "", "", "_BACKUP"
    args := os.Args[1:]
    switch len(args) {
        case 1:
            torrent = args[0]
        case 2:
            torrent, searchDir = args[0], args[1]
        default:
            fmt.Printf("Usage: %s <.torrent file> <Directory to Check>\n", progName)
            os.Exit(1)
    }

    // Check exntension.
    if fileExt := filepath.Ext(torrent); fileExt != ".torrent" {
        fmt.Printf("'%s' Is not a .torrent file!\n", torrent)
        os.Exit(1)
    }

    // Open file now.s
    fmt.Printf("Reading '%s'.\n", torrent)
    f, e := os.Open(torrent)
    check(e)
    defer f.Close()

    // Unmarshal the file into our struct
    mi := MetaInfo{}
    e = bencode.Unmarshal(f, &mi)
    check(e)

    // Print some info
    fmt.Printf("\tName: %s \n", mi.Info.Name)
    if mi.Info.Name != mi.Comment {
        fmt.Printf("\tComment: %s \n", mi.Comment)
    }
    fmt.Printf("\tCreation Date: %s \n", time.Unix(mi.CreationDate, 0))
    fmt.Printf("\tCreated By: %s\n", mi.CreatedBy)

    // Get all the files with their path
    torrentFileMap := map[string]bool{}
    for _, file := range mi.Info.Files {
        filePath := strings.Join(file.Path, "/")
        torrentFileMap[filePath] = true
    }

    // Print one of the files to help figure layout if needed
    fmt.Printf("\tHas %d files.\n", len(torrentFileMap))
    for k, _ := range torrentFileMap {
        fmt.Printf("\tExample: '%s'\n", k)
        break
    }

    // Use the torrent Name as a Dir if it was not specified
    if searchDir == "" {
        searchDir = mi.Info.Name
    }

    // Change the working dir
    e = os.Chdir(searchDir)
    check(e)
    fmt.Printf("Checking directory '%s'\n", searchDir)

    // Get all the files with their path on our working dir
    localFileMap := map[string]bool{}
    e = filepath.Walk(".", func(localFile string, f os.FileInfo, e error) error {
        // Skip dirs and everything under our backup dir
        if !f.IsDir() && !strings.HasPrefix(localFile, backupDir) {
           // Replace each file separator with a slash to be OS consitent
           localFile = filepath.ToSlash(localFile)
           localFileMap[localFile] = true
        }
        return nil
    })
    check(e)
    fmt.Printf("\tHas %d files.\n", len(localFileMap))

    // Get the Files we have locally that aren't on the torrent file
    unwantedFiles := []string{}
    for localFile, _ := range localFileMap {
        if _, ok := torrentFileMap[localFile]; !ok {
            unwantedFiles = append(unwantedFiles, localFile)
        }
    }

    if len(unwantedFiles) > 0 {
        // Show the unwanted files
        fmt.Printf("\tGot %d unwanted Files.\n", len(unwantedFiles))
        fmt.Printf("\tExample: '%s'\n", unwantedFiles[0])
        // Ask user for confirmation
        fmt.Printf("\nEnter YES to move them to '%s': ", backupDir)
        res := ""
        _, e := fmt.Scanln(&res)
        check(e)
        if res != "YES" {
            os.Exit(0)
        }
        // Backup the files
        for _, file := range unwantedFiles {
            newFile := backupDir + "/" + file
            newBase := filepath.Dir(newFile)
            os.MkdirAll(newBase, os.ModePerm);
            fmt.Printf("Moving '%s'\n", file)
            e := os.Rename(file, newFile)
            check(e)
            // Clean orphaned directories
            oldBase := filepath.Dir(file)
            e = recRmOrphanedDir(oldBase)
            check(e)
        }
    }
}
