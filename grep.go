package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Results struct {
	linenumber int
	pathtofile string
}

var totalCount int = 0

// func main(){
// 	results_Path := os.Args[3]
// 	source_path := os.Args[2]
// 	target := os.Args[1]
// 	readDir(source_path,&target,&results_Path)
// 	fmt.Println("'",os.Args[1],"' found ", totalCount," time(s)")
// }

func IgnoreExtension(path string) bool {
	u_ext := []string{"dll"}
	u_ext_found := false

	for i := 0; i < len(u_ext); i++ {
		if strings.Contains(path, u_ext[i]) {
			u_ext_found = true
			break
		}
	}

	return u_ext_found
}

func readDir(pathtobase_dir string, target *string, resultsPath *string) {
	fp, err := os.Create(*resultsPath)
	if err != nil {
		fmt.Println(err)
	}

	defer fp.Close()

	filepath.Walk(pathtobase_dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !IgnoreExtension(path) {
			//os.Create() returns a pointer to the file.. not the file
			readLine(*target, path, fp)
		}
		return nil
	})
}

func readLine(target string, path string, fp *os.File) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("something went wrong on opening file ", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		k := scanner.Text()
		found := findTarget(target, &k)

		if found {
			res := Results{i, path}
			totalCount++
			go writeFilesWhereTargetFound(fp, res)
		}

		i++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("the scanner had some issues ", err)
		fmt.Println("culprit ", path)
	}
}

func findTarget(t string, sub *string) bool {
	return strings.Contains(*sub, t)
}

func writeFilesWhereTargetFound(fp *os.File, r Results) {
	writable := fmt.Sprintf("Path: %s\nLineNumber: %d\n###################\n", r.pathtofile, r.linenumber)
	fp.WriteString(writable)
}
