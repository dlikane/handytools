package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

var (
	rootCmd = &cobra.Command{
		Use:   "goshuffle",
		Short: "Convert org timeline.edl file to timeline-shuffle.edl",
		Long:  `Quick and dirty one`,
		Run: func(cmd *cobra.Command, args []string) {
			processFile(inputFilename)
		},
	}
	inputFilename string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputFilename, "inputFile", "i", "./timeline.edl", "input file (default is ./timeline.edl)")
}
func Execute() {
	_ = rootCmd.Execute()
}

// [head                                              ]
// [from                                    ]
// 001  AX       V     C        12:43:48:01 12:43:48:21 00:00:00:00 00:00:00:20
// * FROM CLIP NAME: A041_07111243_C001.braw
type clip struct {
	head     string
	source   string
	duration int64
}
type edl struct {
	title string
	fcm   string
	clips []clip
}

func processFile(filename string) {
	wd, _ := os.Getwd()
	logrus.Infof("Processing file: %s in %s", filename, wd)
	edl := readEDL(filename)

	ext := filepath.Ext(filename)
	filenameOut := filename[:len(filename)-len(ext)] + "_shuffle" + ext

	writeEDL(filenameOut, edl)
	logrus.Infof("Done: %d edl", len(edl.clips))
}
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := strings.TrimRight(scanner.Text(), " \n\r\t")
		lines = append(lines, str)
	}
	return lines, scanner.Err()
}
func readEDL(path string) edl {
	var edl edl
	lines, err := readLines(path)
	if err != nil {
		logrus.WithError(err).Error("can't load file")
	}
	var c *clip
	lnCnt := 0
	for _, str := range lines {
		switch lnCnt {
		case 0:
			edl.title = str
		case 1:
			edl.fcm = str
		default:
			if lnCnt%3 == 0 {
				if str == "" {
					return edl
				}
				c = &clip{
					head: str[0:53],
				}
			} else if lnCnt%3 == 1 {
				c.source = str
				c.parseDuration()
				edl.clips = append(edl.clips, *c)
			}
		}
		lnCnt++
	}
	return edl
}
func writeEDL(path string, edl edl) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	out := bufio.NewWriter(file)
	defer out.Flush()

	_, _ = out.WriteString(edl.title + "\n")
	_, _ = out.WriteString(edl.fcm + "\n")
	_, _ = out.WriteString("\n")
	pos := int64(0)

	shuffled := make([]clip, len(edl.clips))
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(edl.clips))
	for i, v := range perm {
		shuffled[v] = edl.clips[i]
	}
	for i, c := range shuffled {
		s := fmt.Sprintf("%03d%s%s %s", i, c.head[3:], framesToTimeCode(pos), framesToTimeCode(pos+c.duration))
		_, _ = out.WriteString(s + "\n")
		_, _ = out.WriteString(c.source + "\n")
		_, _ = out.WriteString("\n")
		pos += c.duration
	}
}

func (c *clip) parseDuration() {
	// 001  AX       V     C        12:43:48:01 12:43:48:21 00:00:00:00 00:00:00:20
	start := c.head[29:40]
	end := c.head[41:53]

	startFrames := timeCodeToFrames(start)
	endFrames := timeCodeToFrames(end)
	c.duration = endFrames - startFrames
}

func timeCodeToFrames(str string) int64 {
	h := str[0:2]
	hi, _ := strconv.ParseInt(h, 0, 64)
	m := str[3:5]
	mi, _ := strconv.ParseInt(m, 0, 64)
	s := str[6:8]
	si, _ := strconv.ParseInt(s, 0, 64)
	f := str[9:11]
	fi, _ := strconv.ParseInt(f, 0, 64)
	return hi*60*60*25 + mi*60*25 + si*25 + fi
}

func framesToTimeCode(pos int64) string {
	hi := pos / (60 * 60 * 25)
	mi := (pos % (60 * 60 * 25)) / (60 * 25)
	si := (pos % (60 * 25)) / 25
	fi := pos % 25

	return fmt.Sprintf("%02d:%02d:%02d:%02d", hi, mi, si, fi)
}
