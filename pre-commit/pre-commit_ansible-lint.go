package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Global Variables
var lintExceptions = []string{
	"Package installs should not use latest",
	"Commands should not change things if nothing needs doing",
	"become_user requires become to work as expected",
	"Tasks that run when changed should likely be handlers",
	"Deprecated always_run",
}

var repos = []string{
	"playbooks",
	"platform",
}

var reposMapping = map[string]string{
	"playbooks": "playbooks/*.yml",
	"platform":  "playbooks/*.yml",
}

// Exit - Output some text and exit program
func exit() {
	fmt.Printf(`
	Something has gone wrong, exiting the program.
	These could be one or more probable issues.
		a) .git directory is not found
		b) ansible-lint is not installed or not in $PATH
	`)
	os.Exit(1)
}

// FindGitRepo - Function to find the name of the repo
func FindGitRepo() (repoName string) {
	var flag bool
	// Check Current Repo URL
	cmd := `git config --get remote.origin.url`
	gitURL, stderr := RunShellCmd(cmd)

	// Check if there is no error while running the git command and the name of the repo is not blank
	if stderr != "" {
		flag = true
	} else if gitURL == "" {
		flag = true
	}

	if flag == true {
		exit()
	}

	// Finding the repo name from the repo URL
	slicesOfgitURL := strings.Split(gitURL, "/")
	repoNameWithGit := slicesOfgitURL[len(slicesOfgitURL)-1]

	// Find the repo name
	repoName = strings.Replace(repoNameWithGit, ".git", "", -1)
	return repoName
}

// RunShellCmd - Function to Run a Shell Command
func RunShellCmd(command string) (stdout string, stderr string) {
	cmd := exec.Command("sh", "-c", command)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	return outb.String(), errb.String()
}

func stringsInSlice(line string) bool {
	for _, i := range lintExceptions {
		// encapsulating the pattern between \b
		// We need to do this in order to use 'FindAllString'
		compiledPattern, err := regexp.Compile(`\b` + i + `\b`)
		if err != nil {
			log.Panic("Error In Compiling Regex")
		}

		// Check the length of match. This will return 1 even if there
		// are multple matches. To count all matches in a line use -1
		if len(compiledPattern.FindAllString(line, 1)) != 0 {
			return true
		}
	}
	return false
}

// CheckRepoMapping - Check the path of Ansible Playbooks/Roles to be checked
func CheckRepoMapping(repo string) (ansiblePath string) {
	ansiblePath = reposMapping[repo]
	return
}

// RunLint - Run Lint on the Ansible Playbooks
func RunLint(ansiblePath string) (flag bool) {
	flag = true

	// Create a temporary file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-commit-hook")

	if err != nil {
		print("File could not be created")
		panic(err)
	}

	// Run ansible-lint
	fmt.Println("Running Ansible Lint on " + ansiblePath + " This might take some time...")
	cmd := "ansible-lint -p " + ansiblePath + " > " + tmpFile.Name()
	_, stderr := RunShellCmd(cmd)
	if stderr != "" {
		panic(stderr)
	}

	// Exclude Exceptions
	problemFound := ExcludeException(tmpFile.Name())
	if problemFound == true {
		flag = false
	}

	// Remove the file after all functions have completed
	defer os.Remove(tmpFile.Name())

	return flag
}

// ExcludeException - Delete 'Exceptions' Lines from ansible-lint output
func ExcludeException(tmpFile string) (problemFound bool) {
	problemFound = false
	file, err := os.Open(tmpFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		found := stringsInSlice(string(scanner.Text()))
		if found == false {
			fmt.Println(scanner.Text())
			problemFound = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return problemFound
}

// LintRequired - Checking if we need to check the repo or not
func LintRequired(repoName string) (exists bool, dummy string) {
	exists = false

	fmt.Printf("Checking if we need to lint\n\n")

	if strings.Contains(repoName, "-") {
		// Extracting name of the repo before "-"
		repoName = strings.Split(repoName, "-")[0]
	} else {
		// Remove newline character from string
		repoName = strings.Replace(repoName, "\n", "", -1)
	}

	// Check if extracted repo name exists in 'repos' array
	for _, i := range repos {
		if i == repoName {
			fmt.Printf("Yes, Lint is required.\n\n")
			exists = true
			return exists, repoName
		}
	}
	return exists, repoName
}

// Main - Heart of the Code
func main() {
	// Print Current Time
	x := time.Now()
	fmt.Println(x.Format("2006-01-02 03:04 PM (MST)\n"))

	// Find Current Repo
	repoName := FindGitRepo()

	fmt.Println("Current Repo Found : " + repoName)

	needLint, repo := LintRequired(repoName)
	if needLint == true {
		lintSuccess := RunLint(CheckRepoMapping(repo))
		if lintSuccess == false {
			fmt.Printf("\nErrors found while running ansible-lint, fix your stuff...\n")
			os.Exit(1)
		}
	}
	os.Exit(0)
}
