// Copyright 2022 MobiusCode GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"io"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Config struct {
	// Host URL on which the emails are received
	ImapHost string `json:"imap-host"`
	// Port for the IMAP protocol
	ImapPort int `json:"imap-port"`
	// Email-User for login
	EmailUser string `json:"username"`
	// Corresponding password
	EmailPassword string `json:"password"`
	// Filter for the "from" email field (Name, not email)
	FromFilter string `json:"from-filter"`
	// Regex to be scanned for, first matched group is used as a result
	ContentRegex string `json:"regexp"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <config.json>", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	imapHostAndPort := fmt.Sprintf("%s:%d", config.ImapHost, config.ImapPort)
	imapClient, err := client.DialTLS(imapHostAndPort, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Don't forget to logout
	defer imapClient.Logout()

	// Login
	if err := imapClient.Login(config.EmailUser, config.EmailPassword); err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := imapClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
	seqset := new(imap.SeqSet)
	rangeStart := uint32(1)
	if mbox.Messages > 10 {
		rangeStart = mbox.Messages - 10
	}
	numMails := (mbox.Messages - rangeStart) + 1
	seqset.AddRange(rangeStart, mbox.Messages)

	// Get the whole message body
	section := &imap.BodySectionName{BodyPartName: imap.BodyPartName{Specifier: imap.EntireSpecifier}, Peek: true}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, numMails)
	done := make(chan error, 1)
	go func() {
		done <- imapClient.Fetch(seqset, items, messages)
	}()

	allMessages := toSlice(messages)

	for i := len(allMessages) - 1; i >= 0; i-- {
		msg := allMessages[i]

		r := msg.GetBody(section)
		if r == nil {
			continue
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			log.Fatal(err)
		}
		from := m.Header.Get("From")
		if !strings.Contains(from, config.FromFilter) {
			continue
		}

		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			log.Fatal(err)
		}
		mailContent, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(string(body))))
		if err != nil {
			log.Fatal(err)
		}

		matchContent := findMailContent(string(mailContent), config.ContentRegex)
		if matchContent != nil {
			fmt.Println(*matchContent)
			os.Exit(0)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

func toSlice(c chan *imap.Message) []*imap.Message {
	s := make([]*imap.Message, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}

func findMailContent(content string, pattern string) *string {
	r, _ := regexp.Compile(pattern)
	matches := r.FindStringSubmatch(content)
	if len(matches) <= 1 {
		return nil
	}

	return &matches[1]
}

func loadConfig(filename string) (Config, error) {
	var config Config

	jsonFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config, err
	}

	json.Unmarshal(byteValue, &config)

	replaceEnvVariables(&config)

	return config, nil
}

func replaceEnvVariables(config *Config) {
	v := reflect.ValueOf(config).Elem()
	vType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldJsonName := vType.Field(i).Tag.Get("json")
		value := field.Interface()
		switch value.(type) {
		case string:
			strValue := field.String()
			if strings.HasPrefix(strValue, "$") {
				env, exists := os.LookupEnv(strValue[1:])
				if !exists {
					log.Fatalf("env variable \"%s\" is not set but required for config field \"%s\"",
						strValue[1:], fieldJsonName)
				}
				field.SetString(env)
			}
		}
	}
}
