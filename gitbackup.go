package main

import (
	"fmt"
	"os"
    "os/exec"
)

func listDir(path string) ([]string, []string, error) {
	files := []string{}
	dirs := []string{}
	f, err := os.Open(path)
	if err != nil {
		return dirs, files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		return dirs, files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else {
			files = append(files, file.Name())
		}
	}
	return dirs, files, nil
}

func findGit(path string) ([]string, error) {
	gitDirs := []string{}
	subDir, _, err := listDir(path)
	if err != nil {
		return gitDirs, err
	}
	for _, dir := range subDir {
		if dir == ".git" {
			gitDirs = append(gitDirs, path)
		} else {
			nextLevelDir, err := findGit(path + "/" + dir)
			if err == nil {
				gitDirs = append(gitDirs, nextLevelDir...)
			}
		}
	}
	return gitDirs, nil
}

func gitbackup(path string) {
	os.Chdir(path)
	pwd := exec.Command("pwd")
	stdout, _ := pwd.CombinedOutput()
	fmt.Println(string(stdout))
	fetch := exec.Command("git", "fetch", "--all")
	pull := exec.Command("git", "pull", "--all")
	stdout, _ = fetch.CombinedOutput()
	fmt.Println(string(stdout))
	stdout, _ = pull.CombinedOutput()
	fmt.Println(string(stdout))
}

func absolutePath (arg string) (string) {
	currentDir, _ := os.Getwd()
	if string([]byte(arg)[0:1]) == "/" {
		return arg
	} else {
		return currentDir + "/" + arg
	}
}

func main() {
    path, _ := os.Getwd()
	if len(os.Args) > 1 {
		path = absolutePath(os.Args[1])
	} 

	dirs, err := findGit(path)
	if err == nil {
		for _, dir := range dirs {
			gitbackup(dir)
		}
	} else {
		fmt.Println(err)
	}
}
