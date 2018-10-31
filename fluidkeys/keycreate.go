package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fluidkeys/fluidkeys/backupzip"
	"github.com/fluidkeys/fluidkeys/colour"
	"github.com/fluidkeys/fluidkeys/humanize"
	"github.com/fluidkeys/fluidkeys/out"
	"github.com/fluidkeys/fluidkeys/pgpkey"

	"github.com/sethvargo/go-diceware/diceware"
)

const DicewareNumberOfWords int = 6
const DicewareSeparator string = "."

type DicewarePassword struct {
	words     []string
	separator string
}

func (d DicewarePassword) AsString() string {
	return strings.Join(d.words, d.separator)
}

type generatePgpKeyResult struct {
	pgpKey *pgpkey.PgpKey
	err    error
}

func keyCreate() exitCode {

	if !gpg.IsWorking() {
		out.Print(colour.Warning("\nGPG isn't working on your system 🤒\n\n"))
		out.Print("You can still use FluidKeys to make a key and then " +
			"later import it from your backup.\n\n" +
			"Alternatively, quit now [ctrl-c], install GPG then " +
			"run FluidKeys again.\n\n")
		promptForInput("Press enter to continue. ", true)
	}
	out.Print("\n")
	email := promptForEmail()
	channel := make(chan generatePgpKeyResult)
	go generatePgpKey(email, channel)

	password := generatePassword(DicewareNumberOfWords, DicewareSeparator)

	out.Print("Your key will be protected with this password:\n\n")
	displayPassword(password)
	if !userConfirmedRandomWord(password) {
		out.Print("Those words did not match. Here it is again:\n\n")
		displayPassword(password)
		if !userConfirmedRandomWord(password) {
			out.Print("Those words didn't match again. Quitting...\n")
			os.Exit(1)
		}
	}

	out.Print("Creating key for " + colour.Info(email) + ":\n\n")

	generateJob := <-channel

	if generateJob.err != nil {
		panic(fmt.Sprint("Failed to generate key: ", generateJob.err))
	}
	printSuccessfulAction("Generate key for " + email)

	pushPrivateKeyBackToGpg(generateJob.pgpKey, password.AsString(), &gpg)
	printSuccessfulAction("Store key in " + colour.Info("gpg"))

	fingerprint := generateJob.pgpKey.Fingerprint()
	db.RecordFingerprintImportedIntoGnuPG(fingerprint)
	if err := tryEnableMaintainAutomatically(generateJob.pgpKey, password.AsString()); err == nil {
		printSuccessfulAction("Store password in " + Keyring.Name())
		printSuccessfulAction("Setup automatic maintenance using " + colour.Info("cron"))
	} else {
		printFailedAction("Setup automatic maintenance")
	}

	filename, err := backupzip.OutputZipBackupFile(fluidkeysDirectory, generateJob.pgpKey, password.AsString())
	if err != nil {
		printFailedAction("Make a backup ZIP file")
	}
	directory, _ := filepath.Split(filename)
	printSuccessfulAction("Make a backup ZIP file in")
	out.Print("        " + colour.Info(directory) + "\n\n")

	printSuccess("Successfully created key for " + email)
	out.Print("\n")
	return 0
}

func generatePgpKey(email string, channel chan generatePgpKeyResult) {
	key, err := pgpkey.Generate(email, time.Now(), nil)

	channel <- generatePgpKeyResult{key, err}
}

func promptForEmail() string {
	var email string
	for {
		out.Print("Enter your email address, this will help other people find your key.\n\n")
		email = promptForInput("[email] : ", true)
		if email != "" {
			break
		}
		printFailed("You must provide an email address\n")
	}
	return email
}

func generatePassword(numberOfWords int, separator string) DicewarePassword {
	return DicewarePassword{
		words:     diceware.MustGenerate(numberOfWords),
		separator: separator,
	}
}

func displayPassword(password DicewarePassword) {
	out.Print(out.NoLogCharacter + "   " + colour.Info(password.AsString()) + "\n\n")
	out.Print("If you use a password manager, save it there now.\n\n")
	out.Print(colour.Warning("Store this safely, otherwise you won’t be able to use your key\n\n"))

	promptForInput("Press enter when you've stored it safely. ", true)
}

func userConfirmedRandomWord(password DicewarePassword) bool {
	clearScreen()
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(password.words))
	correctWord := password.words[randomIndex]
	wordOrdinal := humanize.Ordinal(randomIndex + 1)

	out.Print(fmt.Sprintf("Enter the %s word from your password\n\n", wordOrdinal))
	givenWord := promptForInput("["+wordOrdinal+" word] : ", false)
	return givenWord == correctWord
}

func clearScreen() {
	out.Print("\033[H\033[2J")
}
