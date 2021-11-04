package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"

	"github.com/tmax-cloud/cicd-util/pkg/utils"
)

func updateImageUrl() {
	sourcePath := os.Getenv("SOURCE_PATH")
	if sourcePath == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "SOURCE_PATH should be given")
	}
	originalUrl := os.Getenv("IMAGE_URL")
	if originalUrl == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "IMAGE_URL should be given")
	}
	targetFilePath := os.Getenv("TARGET_FILE")
	if targetFilePath == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "TARGET_FILE should be given")
	}

	reg, err := regexp.Compile("([^:/]*(:[0-9]*)?/[^:]*)(:.*)?")
	if err != nil {
		utils.ExitError(log, err, "error while compiling regexp")
	}

	match := reg.FindStringSubmatch(originalUrl)

	baseUrl := match[1]
	tag := match[3]

	if tag != "" {
		log.Info(fmt.Sprintf("image url %s already contains tag... skipping tagging", originalUrl))
		writeUrlToFile(targetFilePath, originalUrl, "")
		os.Exit(0)
	}

	repo, err := git.PlainOpen(sourcePath)
	if err != nil {
		utils.ExitError(log, err, "cannot open git repository")
	}

	ref, err := repo.Head()
	if err != nil {
		utils.ExitError(log, err, "cannot get head of repository")
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		utils.ExitError(log, err, "cannot get log of repository")
	}

	commit, err := cIter.Next()
	if err != nil {
		utils.ExitError(log, err, "cannot get latest commit")
	}

	hash := commit.Hash.String()
	if len(hash) < 7 {
		writeUrlToFile(targetFilePath, baseUrl, hash)
	} else {
		writeUrlToFile(targetFilePath, baseUrl, hash[:7])
	}
}

func writeUrlToFile(filePath, baseUrl, hash string) {
	imageUrl := fmt.Sprintf("%s:%s", baseUrl, hash)
	if hash == "" {
		imageUrl = baseUrl
	}
	log.Info(fmt.Sprintf("Updated image url : %s", imageUrl))
	if err := ioutil.WriteFile(filePath, []byte(imageUrl), 0777); err != nil {
		utils.ExitError(log, err, "cannot write to file")
	}
}
