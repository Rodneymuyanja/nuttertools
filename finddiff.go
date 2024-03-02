package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	// "path/filepath"
)

var totalMatches int = 0
var totalLinesFromFileA int = 0
var totalLinesFromFileB int = 0
var percentage_ float64 = 0
var endString string = "___end___/"

func handler(fp *os.File, channel chan string) {
	scanner := bufio.NewScanner(fp)

	defer endl(channel)

	for scanner.Scan() {
		channel <- scanner.Text()
		time.Sleep(time.Microsecond * 100000)
	}

	if err := scanner.Err(); err != nil {
		defer errHandler("Scanner in handler had a problem")
		panic(err)
	}
}

func endl(channel chan string) {
	channel <- endString
}

func similarity(args ...*int) {
	///0 - input from file A
	///1 - input from file B
	///2 - total matches
	///percentage = ( 2 / (0 + 1)/2 ) * 100
	inputFromFileA := float64(*args[0])
	inputFromFileB := float64(*args[1])
	inputTotalMatches := float64(*args[2])

	averageNumberOflines := (inputFromFileA + inputFromFileB) / 2
	percentage := float64((inputTotalMatches / averageNumberOflines) * 100.0)
	percentage_ = percentage
	fmt.Printf("\rPercentage match: %.1f %%", percentage_)
	//fmt.Println("")
}

func compareLines(p_lineFromFileA *string, p_lineFromFileB *string) bool {
	return strings.ToUpper(*p_lineFromFileA) == strings.ToUpper(*p_lineFromFileB)
}

func readFile(args ...string) {
	///two paths
	fileA := args[0]
	fileB := args[1]

	fp1, fp2 := openFiles(fileA, fileB)
	fp3, fp4 := openFiles(fileA, fileB)

	fileTotalLines(fp1, fp2)

	bufferedChannel := make(chan string, 2)

	go handler(fp3, bufferedChannel)
	go handler(fp4, bufferedChannel)

	go func() {

		defer fp1.Close()
		defer fp2.Close()
		defer fp3.Close()
		defer fp4.Close()

		for {

			//fmt.Println("Waiting...")
			lineFromChannelOne := <-bufferedChannel
			lineFromChannelTwo := <-bufferedChannel

			if lineFromChannelOne == endString || lineFromChannelTwo == endString {
				break
			} else {
				matchOrNot := compareLines(&lineFromChannelOne, &lineFromChannelTwo)

				if matchOrNot {
					totalMatches++
					go similarity(&totalLinesFromFileA, &totalLinesFromFileA, &totalMatches)

				}

			}

		}

		fmt.Printf("\rFinal percentage match:  %.1f %%\n", percentage_)
		fmt.Println("Total lines in ", fileA, ": ", totalLinesFromFileA)
		fmt.Println("Total lines in ", fileB, ": ", totalLinesFromFileB)
		fmt.Println("Total matches in both files: ", totalMatches)
	}()

}

func fileTotalLines(fp1 *os.File, fp2 *os.File) {
	scannerOne := bufio.NewScanner(fp1)
	scannerTwo := bufio.NewScanner(fp2)

	for scannerOne.Scan() {
		totalLinesFromFileA++
	}

	for scannerTwo.Scan() {
		totalLinesFromFileB++
	}

	if err := scannerOne.Err(); err != nil {
		defer errHandler("scannerOne had some problem")
		panic(err)
	}

	if err := scannerTwo.Err(); err != nil {
		defer errHandler("scannerTwo had some problem")
		panic(err)
	}
}

func errHandler(str string) {
	err := recover()
	fmt.Println(str, err)
}

func openFiles(args ...string) (*os.File, *os.File) {
	fileA := args[0]
	fileB := args[1]

	fp1, err := os.Open(fileA)
	if err != nil {
		defer errHandler("Something went wrong on opening file " + fileA)
		panic(err)
	}

	fp2, err := os.Open(fileB)
	if err != nil {
		defer errHandler("Something went wrong on opening file " + fileB)
		panic(err)
	}

	return fp1, fp2
}

func myProgressBar() {
	a := []byte{45}
	b := []byte{35}

	//set up plain bar
	for i := 0; i < 20; i++ {
		fmt.Print(string(a))
	}

	fmt.Print(" 0%")
	fmt.Print("\r")

	//start traversing plain bar
	for i := 0; i < 20; i++ {
		fmt.Print(string(b))
		time.Sleep(time.Second * 1)
	}

	fmt.Print(" 100%\t\t")
}

/*
	my test functions
*/

//total line tests

func main() {

	fileOne := os.Args[1]
	fileTwo := os.Args[2]

	readFile(fileOne, fileTwo)

	var input string
	fmt.Scanln(&input)

}
