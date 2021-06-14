package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/long-schlong-gang/turing"
	"github.com/long-schlong-gang/turing/cypher"

	"github.com/andybalholm/crlf"
	"github.com/teris-io/cli"
)

const DECYPH_CYPHER = 0
const DECYPH_INPUT = 1

var Decypher = cli.NewCommand("decypher", "Decypher a string passed in the console or a file").
	WithArg(
		cli.NewArg("cypher", "The cypher to use (use > to chain cyphers e.g.: 'caesar>vigniere')").
			WithType(cli.TypeString)).
	WithArg(
		cli.NewArg("input", "The plaintext to decypher; Reads from STDIN if left empty").
			WithType(cli.TypeString).
			AsOptional()).
	WithOption(
		cli.NewOption("direct-input", "Decyphers the input directly instead of opening it as a file.").
			WithChar('d').
			WithType(cli.TypeBool)).
	WithOption(
		cli.NewOption("key", "Key to use for decyphering").
			WithChar('k').
			WithType(cli.TypeString)).
	WithOption(
		cli.NewOption("output", "File to put decyphered text").
			WithChar('o').
			WithType(cli.TypeString)).
	WithAction(func(args []string, options map[string]string) int {

		// Verify Args/Options
		cyphName := args[DECYPH_CYPHER]
		input := ""
		if len(args) > 1 {
			input = args[DECYPH_INPUT]
		}

		keyStr, keyExists := options["key"]
		key := turing.ParseKey(keyStr)
		if !keyExists {
			key = turing.NilKey
		}

		_, directInput := options["direct-input"]
		if !directInput && input != "" {
			bytes, err := ioutil.ReadFile(input)
			if err != nil {
				fmt.Println("Error: ", err)
				return STATUS_ERR
			}
			input = string(bytes)
		}

		if input == "" {
			crlfReader := crlf.NewReader(bufio.NewReader(os.Stdin))
			bytes, err := ioutil.ReadAll(crlfReader)
			if err != nil {
				fmt.Println("Error: ", err)
				return STATUS_ERR
			}
			input = string(bytes)

			// Remove final newline
			input = strings.TrimSuffix(input, "\n")
		}

		cyph, err := cypher.RegistryGetCypher(cyphName)
		if err != nil {
			fmt.Println("Error: ", err)
			return STATUS_ERR
		}

		// Decypher Text
		// TODO: Handle chained encypherings (plus encyphered keys)
		plaintext := cyph.Decypher(input, key)

		outFile, isOutToFile := options["output"]
		if isOutToFile {
			ioutil.WriteFile(outFile, []byte(plaintext), OS_FILE_MODE)
		} else {
			fmt.Printf("Decyphered Plaintext: \"%v\"\n", plaintext)
		}

		// TODO: Add time taken

		return STATUS_OK
	})
