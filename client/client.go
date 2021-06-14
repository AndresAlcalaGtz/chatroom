package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"../env"
	"../msg"
)

func main() {
	send := make(chan msg.Message, 1)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("\nUser: ")
	scanner.Scan()
	user := scanner.Text()

	go client(send, user)
	sendMessage(send, user, msg.Login, []byte{}, false)

	var temp string
	var option string

	fmt.Println("\nCLIENT")
	fmt.Println("[1] Send message")
	fmt.Println("[2] Send file")
	fmt.Println("[0] Log out")

	for option != "0" {
		fmt.Println("\n")
		fmt.Scanln(&option)
		fmt.Println()

		switch option {
		case "1":
			fmt.Println("Write a message!")
			scanner.Scan()
			temp = scanner.Text()
			sendMessage(send, user, temp, []byte{}, false)

		case "2":
			fmt.Println("Write a file name!")
			scanner.Scan()
			temp = scanner.Text()
			sendFile(send, user, temp)

		case "0":
			fmt.Println("Goodbye!")
			sendMessage(send, user, msg.Logout, []byte{}, false)
			time.Sleep(time.Second)

		default:
			fmt.Println("Invalid option!")
		}
	}
}

func client(send chan msg.Message, username string) {
	var input msg.Message

	connection, err := net.Dial(env.ConnectionType, env.ConnectionPort)
	if err != nil {
		fmt.Println("CONNECTION ERROR:", err)
		return
	}
	defer connection.Close()

	for {
		select {
		case output := <-send:
			err := gob.NewEncoder(connection).Encode(output)
			if err != nil {
				fmt.Println("ENCODING ERROR:", err)
				return
			}

		default:
			continue
		}

		err := gob.NewDecoder(connection).Decode(&input)
		if err != nil {
			fmt.Println("DECODING ERROR:", err)
			return
		}
		if input.File {
			saveFile(input, username)
		}
		fmt.Println(input.ToString())
	}
}

func sendMessage(send chan msg.Message, author string, text string, bytes []byte, file bool) {
	message := msg.Message{Author: author, Text: text, Bytes: bytes, File: file}

	send <- message
}

func sendFile(send chan msg.Message, author string, text string) {
	bytes, err := ioutil.ReadFile(text)
	if err != nil {
		fmt.Println("READING ERROR:", err)
		return
	}

	filepath := strings.Split(text, "/")
	filename := filepath[len(filepath)-1]

	sendMessage(send, author, filename, bytes, true)
}

func saveFile(message msg.Message, username string) {
	filepath := "client/" + strings.ToLower(username)

	_, err := os.Stat(filepath)
	if err != nil {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			fmt.Println("FOLDER ERROR:", err)
			return
		}
	}

	err = ioutil.WriteFile(filepath+"/"+message.Text, message.Bytes, os.ModePerm)
	if err != nil {
		fmt.Println("WRITING ERROR:", err)
		return
	}
}
