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
	"regexp"
)

func main() {
	filetype, fps, dlog := setConfig()
	tileName := "tile.png"
	paletteName := "palette.png"
	tileCmd, paletteCmd, gifCmd := genCmds(filetype, fps, dlog, tileName, paletteName)

	nameMap := handleInputFiles(filetype, dlog)
	if len(nameMap) == 0 {
		prln("no files with type '"+filetype+"' found. Exiting", dlog)
		return
		// todo validate sizes
	}

	errTile := runCmd(tileCmd, dlog)
	if errTile != nil {
		nameBack(&nameMap)
		os.Exit(1)
	}
	errPalette := runCmd(paletteCmd, dlog)
	if errPalette == nil {
		errGif := runCmd(gifCmd, dlog)
		if errGif == nil {
			prln("output.gif generated", dlog)
		}
	}

	nameBack(&nameMap)
	delTmpFile(tileName, errTile, dlog)
	delTmpFile(paletteName, errPalette, dlog)
	if dlog {
		prln("Press Enter to exit", dlog)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		scanner.Text()
	}
}

func delTmpFile(name string, genError error, dlog bool) {
	if genError == nil {
		e := os.Remove(name)
		if  e != nil {
			prln("failed to delete temp file. File can be deleted manually " + name, dlog)
		} else {
			prln("tmp file " + name + " deleted successfully", dlog)
		}
	}
}

func genCmds(filetype string, fps string, dlog bool, tileName string, paletteName string) ([]string, []string, []string) {
	tileCmd := []string{
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-i",
		"%6d." + filetype,
		"-vf",
		"format=rgb24,tile=10x10:color=black",
		"-frames:v",
		"1",
		tileName,
	}
	paletteCmd := []string{
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-i",
		tileName,
		"-vf",
		"palettegen=max_colors=64:reserve_transparent=1:stats_mode=single",
		paletteName,
	}
	gifCmd := []string{
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-framerate",
		fps,
		"-i",
		"%6d." + filetype,
		"-i",
		paletteName,
		"-filter_complex",
		"split[s0][s1]; [s0]palettegen= max_colors=256: stats_mode=single: reserve_transparent=on: transparency_color=ffffff[p]; [s1][p]paletteuse=new=1",
		"output.gif",
	}
	return tileCmd, paletteCmd, gifCmd
}

func runCmd(parts []string, dlog bool) error {
	prln(strings.Join(parts[:], " "), dlog)
	cmd := exec.Command(parts[0], parts[1:]...)
	stdAndErr, err := cmd.CombinedOutput()
	if (err != nil && dlog) {
		fmt.Println(string(stdAndErr))
	}
	return err
}

func setConfig() (string, string, bool) {
	var filetype = "png"
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
				if len(filetypeFound) != 0 && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(filetypeFound) {
					filetype = filetypeFound
				}
			} else if line == 1 {
				fpsFound := strings.Trim(scanner.Text(), " ")
				if len(fpsFound) != 0 && regexp.MustCompile(`^[0-9]*$`).MatchString(fpsFound) {
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
	prln("running with parameters:", debug)
	prln("	filetype "+filetype, debug)
	prln("	fps "+fps, debug)
	prln("	debug "+fmt.Sprint(debug), debug)
	return filetype, fps, debug
}

func handleInputFiles(filetype string, dlog bool) map[string]string {
	var nameMap = make(map[string]string)
	entries, err := os.ReadDir("./")
	if err != nil {
		prln(err.Error(), dlog)
		return nameMap
	} else {
		slices.SortFunc(entries, func(a, b fs.DirEntry) int {
			return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
		})
	}

	i := 0
	for _, e := range entries {
		curName := e.Name()
		if strings.HasSuffix(strings.ToLower(curName), filetype) {
			i++
			newName := fmt.Sprintf("%06d.", i) + filetype
			nameMap[newName] = curName
			os.Rename(curName, newName)
		}
	}
	prln("Found files to convert into gif. Renaming them temporarily " + fmt.Sprint(nameMap), dlog)
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
