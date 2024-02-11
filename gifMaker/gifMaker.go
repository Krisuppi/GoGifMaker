package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	filetype, fps, dlog := setConfig()
	nameMap := handleInputFiles(filetype, dlog)

	if len(nameMap) == 0 {
		prln("no files with type '"+filetype+"' found. Exiting", dlog)
		return
	}

	cmdString := "ffmpeg -framerate " + fps + " -i %6d" + filetype + " output.gif -y"
	prln("cmd: "+cmdString, dlog)
	parts := strings.Fields(cmdString)
	cmd := exec.Command(parts[0], parts[1:]...)
	outStdAndErr, err := cmd.CombinedOutput()

	if err != nil {
		nameBack(&nameMap)
		prln(string(outStdAndErr), dlog)
		prln(err.Error(), dlog)
		os.Exit(1)
	} else {
		// to view output if needed
		// fmt.Printf("%s\n", outStdAndErr)
		nameBack(&nameMap)
		prln("output.gif generated", dlog)
	}
	if dlog {
		prln("Press Enter to exit", dlog)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		scanner.Text()
	}
}

func setConfig() (string, string, bool) {
	var filetype = ".png"
	var fps = "30"
	var debug = true

	config, err := os.Open("config.txt")
	if err == nil {
		defer config.Close()

		scanner := bufio.NewScanner(config)
		line := 0
		for scanner.Scan() {
			if line == 0 {
				filetypeFound := strings.ToLower(strings.Trim(scanner.Text(), " "))
				if len(filetypeFound) != 0 {
					filetype = filetypeFound
				}
			} else if line == 1 {
				fpsFound := strings.Trim(scanner.Text(), " ")
				if len(fpsFound) != 0 {
					fps = fpsFound
				}
			} else if line == 2 {
				debugFound := strings.Trim(scanner.Text(), " ")
				if len(debugFound) != 0 {
					debug = false
				}
			} else {
				break
			}
			line++
		}
	}
	prln("filetype "+filetype, debug)
	prln("fps "+fps, debug)
	prln("debug "+fmt.Sprint(debug), debug)
	return filetype, fps, debug
}

func handleInputFiles(filetype string, dlog bool) map[string]string {
	entries, err := os.ReadDir("./")
	if err != nil {
		prln(err.Error(), dlog)
		os.Exit(1)
	} else {
		slices.SortFunc(entries, func(a, b fs.DirEntry) int {
			return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
		})
	}

	var nameMap = make(map[string]string)
	i := 0
	for _, e := range entries {
		curName := e.Name()
		if strings.HasSuffix(strings.ToLower(curName), filetype) {
			i++
			newName := fmt.Sprintf("%06d", i) + filetype
			nameMap[newName] = curName
			os.Rename(curName, newName)
		}
	}

	if dlog {
		fmt.Println("filemap to be giffified below")
		fmt.Println(&nameMap)
	}
	return nameMap
}
func nameBack(nameMap *map[string]string) {
	for temp, og := range *nameMap {
		os.Rename(temp, og)
	}
}

func prln(text string, doLog bool) {
	if doLog {
		fmt.Println(text)
	}
}
