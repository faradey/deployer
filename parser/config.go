package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type ConfigRow struct {
	Name string
	Data string
}

func GetConfig(dir string) []ConfigRow {
	lines := getLines(dir)
	var commands []ConfigRow
	for _, val := range lines {
		if val != "" {
			rowBlock := strings.SplitN(val, " ", 2)
			if len(rowBlock) > 1 {
				commands = append(commands, ConfigRow{Name: rowBlock[0], Data: rowBlock[1]})
			} else if len(rowBlock) == 1 {
				commands = append(commands, ConfigRow{Name: rowBlock[0], Data: ""})
			}
		}
	}

	return commands
}

func getLines(dir string) []string {
	fileName := "deployer-config"
	var rows []string
	file, err := os.Open(dir + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[len(line)-2:] != "/" {
			rows = append(rows, line)
			i++
		} else {
			rows[i] += line[:len(line)-2]
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return rows
}
