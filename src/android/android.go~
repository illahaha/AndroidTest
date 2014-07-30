package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Script_item struct {
	Name  string `xml:"name,attr"`
	Value string
}
type Scripts struct {
	XMLName xml.Name      `xml:"scripts"`
	Script  []Script_item `xml:"script"`
}

func getPackName(path string) string {

	cmd := exec.Command("aapt", "dump", "badging", path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	out_bytes := new(bytes.Buffer)
	output_done_channel := make(chan bool)
	go func() {
		io.Copy(out_bytes, stdout)
		//fmt.Printf("%s\n", out_bytes)
		output_done_channel <- true
	}()
	<-output_done_channel
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	//regex package name
	//r, _ := regexp.Compile("package:\\sname='([a-zA-Z0-9\\.]*)'")
	r, _ := regexp.Compile("package:\\sname='(.+?)'")

	packageName := r.FindStringSubmatch(out_bytes.String())

	return packageName[1]
}

func runCmd(name string, arg ...string) (int, string) {

	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return -1, ""
	}
	cmd.Wait()
	return 0, ""
}

func runScript(appName string, packageName string, args [6]string) (int, string) {

	outbuf := ""

	//read script
	xmlFile, err := os.Open(args[2])
	if err != nil {
		panic(err)
		return -1, ""
	}

	defer xmlFile.Close()

	fi, _ := xmlFile.Stat()
	buf := make([]byte, fi.Size())
	_, err = xmlFile.Read(buf)
	if err != nil {
		panic(err)
		return -1, ""
	}

	v := Scripts{}
	err = xml.Unmarshal(buf, &v)
	if err != nil {
		panic(err)
		return -1, ""
	}

	for _, value := range v.Script {
		text := strings.Replace(value.Value, "[package]", packageName, -1)
		fmt.Println("testin script name ", value.Name)

		params := make([]string, len(strings.Fields(text))+2)
		params[0] = args[0]
		params[1] = args[1]

		for i, s := range strings.Fields(text) {
			params[i+2] = s
		}

		parts := params[:]

		runCmd(appName, parts...)
	}

	return 0, outbuf
}

func main() {

	fmt.Println("******android it******")

	argsWithoutProg := os.Args[1:]

	//fmt.Println(argsWithoutProg)

	argsCmd := [6]string{}

	//fetch package name
	packageName := getPackName(argsWithoutProg[2])

	//base command
	appName := "adb"
	argsCmd[0] = "-s"
	argsCmd[1] = argsWithoutProg[0] //device

	//install command
	argsCmd[2] = "install"
	argsCmd[3] = "-r"
	argsCmd[4] = argsWithoutProg[2] //apk
	parts := argsCmd[:5]
	fmt.Println("install apk by adb command", parts)
	runCmd(appName, parts...)

	//monitor logcat
	logCat := exec.Command("logfilter", appName, argsWithoutProg[0], packageName, "filters.xml", argsWithoutProg[2]+"-"+argsWithoutProg[0], strconv.Itoa(os.Getpid()))

	go func() {
		//clear logcat
		runCmd("adb", "-s", argsWithoutProg[0], "logcat", "-c")
		logCat.Start()
		defer logCat.Wait()
	}()

	//run script
	argsCmd[2] = argsWithoutProg[3]
	if argsCmd[2] != "none" {
		runScript(appName, packageName, argsCmd)
	}

	//uninstall command
	argsCmd[2] = "uninstall"
	argsCmd[3] = packageName
	parts = argsCmd[:4]
	fmt.Println("uninstall apk by adb command", parts)
	runCmd(appName, parts...)
	fmt.Println("android log out")
}
