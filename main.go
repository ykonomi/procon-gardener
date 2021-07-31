package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/PuerkitoBio/goquery"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
	"github.com/thoas/go-funk"
	cli "github.com/urfave/cli/v2"
)

const APP_NAME = "procon-gardener"
const ATCODER_API_SUBMISSION_URL = "https://kenkoooo.com/atcoder/atcoder-api/results?user="

type AtCoderSubmission struct {
	ID            int     `json:"id"`
	EpochSecond   int64   `json:"epoch_second"`
	ProblemID     string  `json:"problem_id"`
	ContestID     string  `json:"contest_id"`
	UserID        string  `json:"user_id"`
	Language      string  `json:"language"`
	Point         float64 `json:"point"`
	Length        int     `json:"length"`
	Result        string  `json:"result"`
	ExecutionTime int     `json:"execution_time"`
}

func isDirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type Service struct {
	RepositoryPath string `json:"repository_path"`
	DirectoryPath  string `json:"directory_path"`
	KeyFileName    string `json:"key_file_name"`
	UserID         string `json:"user_id"`
	UserEmail      string `json:"user_email"`
}
type Config struct {
	Atcoder Service `json:"atcoder"`
}

func pow(a, x int) int {
	r := 1
	for x > 0 {
		if x&1 == 1 {
			r *= a
		}
		a *= a
		x >>= 1
	}
	return r
}

// a -> 1, b -> 2,...,z -> 26
// aa -> 27...
func toNumber(s string) int {
	power := pow(26, len(s)-1)

	r := 0
	for i := range s {
		switch {
		case 97 <= s[i] && s[i] <= 122:
			r += power * int(s[i]-97+1)
		}

		power /= 26
	}

	return r
}

func languageExtension(language string) string {
	//e.g C++14 (GCC 5.4.1)
	//C++14
	language = strings.Split(language, "(")[0]
	//remove extra last whitespace
	language = language[:len(language)-1]
	if strings.HasPrefix(language, "C++") {
		return ".cpp"
	}
	if strings.HasPrefix(language, "Bash") {
		return ".sh"
	}

	//C (GCC 5.4.1)
	//C (Clang 3.8.0)
	if language == "C" {
		return ".c"
	}

	if language == "C#" {
		return ".cs"
	}

	if language == "Clojure" {
		return ".clj"
	}

	if strings.HasPrefix(language, "Common Lisp") {
		return ".lisp"
	}

	//D (DMD64 v2.070.1)
	if language == "D" {
		return ".d"
	}

	if language == "Fortran" {
		return ".f08"
	}

	if language == "Go" {
		return ".go"
	}

	if language == "Haskell" {
		return ".hs"
	}

	if language == "JavaScript" {
		return ".js"
	}
	if language == "Java" {
		return ".java"
	}
	if language == "OCaml" {
		return ".ml"
	}

	if language == "Pascal" {
		return ".pas"
	}

	if language == "Perl" {
		return ".pl"
	}

	if language == "PHP" {
		return ".php"
	}

	if strings.HasPrefix(language, "Python") {
		return ".py"
	}

	if language == "Ruby" {
		return ".rb"
	}

	if language == "Scala" {
		return ".scala"
	}

	if language == "Scheme" {
		return ".scm"
	}

	if language == "Main.txt" {
		return ".txt"
	}

	if language == "Visual Basic" {
		return ".vb"
	}

	if language == "Objective-C" {
		return ".m"
	}

	if language == "Swift" {
		return ".swift"
	}

	if language == "Rust" {
		return ".rs"
	}

	if language == "Sed" {
		return ".sed"
	}

	if language == "Awk" {
		return ".awk"
	}

	if language == "Brainfuck" {
		return ".bf"
	}

	if language == "Standard ML" {
		return ".sml"
	}

	if strings.HasPrefix(language, "PyPy") {
		return ".py"
	}

	if language == "Crystal" {
		return ".cr"
	}

	if language == "F#" {
		return ".fs"
	}

	if language == "Unlambda" {
		return ".unl"
	}

	if language == "Lua" {
		return ".lua"
	}

	if language == "LuaJIT" {
		return ".lua"
	}

	if language == "MoonScript" {
		return ".moon"
	}

	if language == "Ceylon" {
		return ".ceylon"
	}

	if language == "Julia" {
		return ".jl"
	}

	if language == "Octave" {
		return ".m"
	}

	if language == "Nim" {
		return ".nim"
	}

	if language == "TypeScript" {
		return ".ts"
	}

	if language == "Perl6" {
		return ".p6"
	}

	if language == "Kotlin" {
		return ".kt"
	}

	if language == "COBOL" {
		return ".cob"
	}

	log.Printf("Unknown ... %s", language)
	return ""
}

func initCmd(strict bool) {

	log.Println("Initialize your config...")
	home, err := homedir.Dir()
	if err != nil {
		log.Println(err)
		return
	}
	configDir := filepath.Join(home, "."+APP_NAME)
	if !isDirExist(configDir) {
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			log.Println(err)
			return
		}
	}

	configFile := filepath.Join(configDir, "config.json")
	if strict || !isFileExist(configFile) {
		//initial config
		atcoder := Service{RepositoryPath: "", UserID: ""}

		config := Config{Atcoder: atcoder}

		jsonBytes, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			log.Println(err)
			return
		}
		json := string(jsonBytes)
		file, err := os.Create(filepath.Join(configDir, "config.json"))
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		file.WriteString(json)
	}
	log.Println("Initialized your config at ", configFile)
}

func loadConfig() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(home, "."+APP_NAME)
	configFile := filepath.Join(configDir, "config.json")
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Config
	if err = json.Unmarshal(bytes, &config); err != nil {
		log.Println(err)
		return nil, err
	}
	return &config, nil
}

func archiveFile(config *Config, code, filePath string) error {
	abstPath := filepath.Join(config.Atcoder.RepositoryPath, filePath)
	dirPath, _ := filepath.Split(abstPath)

	if err := os.MkdirAll(dirPath, 0700); err != nil {
		log.Fatal(err)
		return err
	}

	file, err := os.OpenFile(abstPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(code); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// in-place に詰め直す
func extractLatestSubmissionsPerSubmission(s []AtCoderSubmission) []AtCoderSubmission {
	sort.Slice(s, func(i, j int) bool {
		return s[i].EpochSecond > s[j].EpochSecond
	})

	r := s[:0]
	set := map[string]struct{}{}

	for i := range s {
		tmp := s[i]
		key := tmp.ContestID + "_" + tmp.ProblemID

		_, ok := set[key]
		if ok {
			continue
		}

		set[key] = struct{}{}
		r = append(r, tmp)
	}
	return r

}

func extractUnArchivedSubmission(s []AtCoderSubmission, archivedKeys map[string]struct{}) []AtCoderSubmission {
	r := s[:0]

	for i := range s {
		tmp := s[i]
		_, ok := archivedKeys[tmp.ProblemID]
		if !ok {
			r = append(r, tmp)
		}
	}
	return r
}

func extractAc(s []AtCoderSubmission) []AtCoderSubmission {
	r := s[:0]

	for i := range s {
		tmp := s[i]
		if tmp.Result == "AC" {
			r = append(r, tmp)
		}
	}

	return r
}

func loadArchivedProgramId(c *Config) (map[string]struct{}, error) {
	r := map[string]struct{}{}
	fp, err := os.Open(filepath.Join(c.Atcoder.RepositoryPath, c.Atcoder.DirectoryPath, c.Atcoder.KeyFileName))
	if err != nil {
		log.Println(err)
		return r, err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		r[scanner.Text()] = struct{}{}
	}

	return r, nil
}

func dryRunCmd() {
	config, err := loadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := http.Get(ATCODER_API_SUBMISSION_URL + config.Atcoder.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var submissions []AtCoderSubmission
	err = json.Unmarshal(bytes, &submissions)
	if err != nil {
		log.Println(err)
		return
	}

	archivedKeys, err := loadArchivedProgramId(config)
	if err != nil {
		log.Println(err)
		return
	}

	//only ac
	submissions = extractAc(submissions)
	//skip the already archived code
	submissions = extractUnArchivedSubmission(submissions, archivedKeys)
	// filter latest submission for each problem
	submissions = extractLatestSubmissionsPerSubmission(submissions)

	for i := range submissions {
		fmt.Printf("%v\n", submissions[i])
	}

	for i := range submissions {
		tmp := strings.Split(submissions[i].ProblemID, "_")
		if tmp[0] != submissions[i].ContestID {
			fmt.Printf("info: a part of %v problems are saved in %v \n", submissions[i].ContestID, tmp[0])
		}
	}

}

func archiveCmd() {
	config, err := loadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := http.Get(ATCODER_API_SUBMISSION_URL + config.Atcoder.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var submissions []AtCoderSubmission
	err = json.Unmarshal(bytes, &submissions)
	if err != nil {
		log.Println(err)
		return
	}

	archivedKeys, err := loadArchivedProgramId(config)
	if err != nil {
		log.Println(err)
		return
	}

	//only ac
	submissions = extractAc(submissions)
	//skip the already archived code
	submissions = extractUnArchivedSubmission(submissions, archivedKeys)
	// filter latest submission for each problem
	submissions = extractLatestSubmissionsPerSubmission(submissions)

	startTime := time.Now()
	log.Printf("Archiving %d code...", len(submissions))

	successFileName := []string{}
	funk.ForEach(submissions, func(s AtCoderSubmission) {
		url := fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%s", s.ContestID, strconv.Itoa(s.ID))

		//log.Printf("Requesting... %s", url)
		elapsedTime := time.Now().Sub(startTime)
		if elapsedTime.Milliseconds() < 1500 {
			sleepTime := time.Duration(1500 - elapsedTime.Milliseconds())
			time.Sleep(time.Millisecond * sleepTime)
		}
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}

		defer resp.Body.Close()

		startTime = time.Now()
		if err != nil {
			log.Println(err)
			return
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		language := s.Language
		epochSecond := s.EpochSecond
		// ex) abc002_a -> abc002 a
		contestId := strings.Split(s.ProblemID, "_")[0]
		programName := strings.Split(s.ProblemID, "_")[1]
		// ex) tenkei90 cl -> 90
		if contestId == "typical90" {
			programName = strconv.Itoa(toNumber(programName))
		}

		fileName := programName + languageExtension(language)
		filePath := filepath.Join(config.Atcoder.DirectoryPath, contestId, fileName)
		message := fmt.Sprintf("[AC] %s %s", contestId, programName)

		doc.Find(".linenums").Each(func(i int, gs *goquery.Selection) {
			code := gs.Text()
			if code == "" {
				log.Print("Empty string...")
				return
			}

			if err = archiveFile(config, code, filePath); err != nil {
				log.Println("Fail to archive the code at", filePath)
				return
			}
			fmt.Println("archived the code at ", filePath)

			if err := commit(config, filePath, message, epochSecond); err != nil {
				log.Println("Error: fail to commit ", filePath)
				return
			}

			successFileName = append(successFileName, s.ProblemID)
			return
		})

	})

	filePath := filepath.Join(config.Atcoder.DirectoryPath, config.Atcoder.KeyFileName)
	tmp := ""
	for i := range successFileName {
		tmp += successFileName[i] + "\n"
	}
	if err = archiveFile(config, tmp, filePath); err != nil {
		log.Println("Fail to archive the code at", filePath)
		return
	}

	if err := commit(config, filePath, "Update a Key file", time.Now().Unix()); err != nil {
		log.Println("Error: fail to commit a key file")
		return
	}

}

func commit(config *Config, filePath, message string, epochSecond int64) error {
	r, err := git.PlainOpen(config.Atcoder.RepositoryPath)
	if err != nil {
		log.Println(err)
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := w.Add(filePath); err != nil {
		log.Println(err)
		return err
	}

	if _, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Atcoder.UserID,
			Email: config.Atcoder.UserEmail,
			When:  time.Unix(epochSecond, 0),
		},
	}); err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("commited the code at ", filePath)

	return nil
}

func validateConfig(config Config) bool {
	//TODO check path
	return false
}

func editCmd() {

	home, err := homedir.Dir()
	if err != nil {
		log.Println(err)
		return
	}
	configFile := filepath.Join(home, "."+APP_NAME, "config.json")
	//Config file not found, force to run an init cmd
	if !isFileExist(configFile) {
		initCmd(true)
	}

	editor := os.Getenv("EDITOR")
	if editor != "" {
		c := exec.Command(editor, configFile)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	} else {
		open.Run(configFile)
	}

}

func main() {

	app := cli.App{Name: "procon-gardener", Usage: "archive your AC submissions",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "initialize your config",
				Action: func(c *cli.Context) error {
					initCmd(true)
					return nil
				},
			},
			{
				Name:    "archive",
				Aliases: []string{"a"},
				Usage:   "archive your AC submissions",
				Action: func(c *cli.Context) error {
					archiveCmd()
					return nil
				},
			},
			{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "list your AC commmit.",
				Action: func(c *cli.Context) error {
					dryRunCmd()
					return nil
				},
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "edit your config file",
				Action: func(c *cli.Context) error {
					editCmd()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
