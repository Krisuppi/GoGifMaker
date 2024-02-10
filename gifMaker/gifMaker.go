package main

import (
    "fmt"
    "os"
    "log"
	"strings"
	"os/exec"
	"slices"
	"cmp"
	"io/fs"
	"bufio"
)

func main() {
	var filetype = ".png"
	config, err := os.Open("config.txt")
	if err == nil {
		defer config.Close()
 		scanner := bufio.NewScanner(config)
		scanner.Scan()
	    filetype = strings.Trim(scanner.Text(), " ")
		fmt.Println("config found. selecting for file type '" + filetype + "' for gif generation")
	} else {
		fmt.Println("config not found or unable to open. using default .png for gif generation")
	}
    entries, err := os.ReadDir("./")
    if err != nil {
        log.Fatal(err)
    } else {
		slices.SortFunc(entries, func(a, b fs.DirEntry) int {
			return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
		})
	}
	var nameMap = make(map[string]string)
    for i, e := range entries {
    	curName := e.Name()
		if (strings.HasSuffix(curName, filetype)) {
			newName := fmt.Sprintf("%06d", i) + filetype
			nameMap[newName] = curName
			os.Rename(curName, newName)
		}
    }
	if (len(nameMap) == 0) {
	fmt.Println("no files with type '" + filetype + "' found. Exiting")
		return
	}
	cmdString := "ffmpeg -i %6d" + filetype +" output.gif -y"
	parts := strings.Fields(cmdString)
	cmd := exec.Command(parts[0], parts[1:]...)
	outStdAndErr, err := cmd.CombinedOutput()
	if err != nil {
		nameBack(&nameMap)
		fmt.Printf("%s\n", outStdAndErr)
		fmt.Println("cmd was: " + cmdString)
		log.Fatal(err)
	} else {
		// to view output if needed
		// fmt.Printf("%s\n", outStdAndErr)
		nameBack(&nameMap)
		fmt.Println("output.gif generated")
	}
}

func nameBack(nameMap *map[string]string) {
	for temp, og := range *nameMap {
		os.Rename(temp, og)
	}
}
