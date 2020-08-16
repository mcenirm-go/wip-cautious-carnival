/**
 * @license
 * Copyright Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// [START sheets_quickstart]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mcenirm-go/wip-cautious-carnival/carnival"
	"github.com/vharitonsky/iniflags"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

const (
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokenFile       = "token.json"
	credentialsFile = "credentials.json"
)

var (
	clientScopes = []string{sheets.SpreadsheetsReadonlyScope}
)

var (
	spreadsheetID, readRange, reportHeader string
	reportIndices                          = carnival.NewListOfUintsForFlag(0, 4)
)

func init() {
	flag.StringVar(&spreadsheetID, "spreadsheetID", "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms", "Identifier of spreadsheet to read")
	flag.StringVar(&readRange, "readRange", "Class Data!A2:E", "Range of cells to read")
	flag.StringVar(&reportHeader, "reportHeader", "Name, Major", "Heading to print before data")
	flag.Var(reportIndices, "reportIndices", "Indices of columns to print, zero-based relative to readRange")
}

func main() {
	iniflags.Parse()
	b, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, clientScopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, tokenFile)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		indices := reportIndices.Values()
		reportFormatPieces := make([]string, len(indices))
		reportValues := make([]interface{}, len(indices))
		for i := range indices {
			reportFormatPieces[i] = "%s"
		}
		reportFormat := strings.Join(reportFormatPieces, ", ") + "\n"

		fmt.Println(reportHeader + ":")
		for _, row := range resp.Values {
			for i, rowIndex := range indices {
				reportValues[i] = row[rowIndex]
			}
			fmt.Printf(reportFormat, reportValues...)
		}
	}
}

// [END sheets_quickstart]
