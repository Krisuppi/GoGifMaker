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
	nameMap := handleInputFiles(filetype, dlog)

	if len(nameMap) == 0 {
		prln("no files with type '"+filetype+"' found. Exiting", dlog)
		return
		// todo validate sizes
	}
	tileName, errTile := genGlobalTile(filetype, dlog)
	if errTile != nil {
		prln(errTile.Error(), dlog)
		nameBack(&nameMap)
		os.Exit(1)
	}

	paletteName, errPalette := genPalette(filetype, tileName, dlog)
	if errPalette == nil {
		errGif := genGif(filetype, paletteName, fps, dlog)
		if errGif != nil {
			prln("output.gif generated", dlog)
		} else {
			prln(errGif.Error(), dlog)
		}
	}

	// Clean up. reverse rename of nameMap and delete tmp and palette files
	nameBack(&nameMap)
	if errTile == nil {
		e := os.Remove(tileName)
		if  e != nil {
			prln("failed to delete temp file for tile. File can be deleted manually " + tileName, dlog)
		} else {
			prln("tmp tile file deleted successfully " + tileName, dlog)
		}
	}
	if errPalette == nil {
		e := os.Remove(paletteName)
		if e != nil {
			prln("failed to delete palette file. File can be deleted manually " + paletteName, dlog)
		} else {
			prln("tmp palette file deleted successfully " + paletteName, dlog)	
		}
	}
	if dlog {
		prln("Press Enter to exit", dlog)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		scanner.Text()
	}
}

func genGlobalTile(filetype string, dlog bool) (string, error) {
	tileName := "tmp.png"
	parts := []string{
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
	prln("part 1. create tile: " + strings.Join(parts[:], " "), dlog)
	cmd := exec.Command(parts[0], parts[1:]...)
	outStdAndErr, err := cmd.CombinedOutput()
	if (err != nil && dlog) {
		fmt.Println(string(outStdAndErr))
	}
	return tileName, err
}

func genPalette(filetype string, tileName string, dlog bool) (string, error) {
	palName := "palette.png"
	parts := []string{
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-i",
		tileName,
		"-vf",
		"palettegen=max_colors=64:reserve_transparent=1:stats_mode=single",
		palName,
	}
	prln("part 2. create palette: " + strings.Join(parts[:], " "), dlog)
	cmd := exec.Command(parts[0], parts[1:]...)
	err := cmd.Run()
	return palName, err
}

func genGif(filetype string, paletteName string, fps string, dlog bool) error {
	parts := []string{
		"ffmpeg",
		"-y",
		"hide_banner",
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
	prln("part 3. create Gif: " + strings.Join(parts[:], " "), dlog)
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run()
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
			newName := fmt.Sprintf("%06d.", i) + filetype
			nameMap[newName] = curName
			os.Rename(curName, newName)
		}
	}

	if dlog {
		fmt.Print("Found files to convert into gif. Renaming them temporarily: ")
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
