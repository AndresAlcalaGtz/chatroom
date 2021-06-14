package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"../env"
	"../msg"
)

var (
	users []net.Conn
	chat  []msg.Message
)

func main() {
	go server()

	var option string

	fmt.Println("\nSERVER")
	fmt.Println("[1] Print chat")
	fmt.Println("[2] Save chat")
	fmt.Println("[0] Close server")

	for option != "0" {
		fmt.Println("\n")
		fmt.Scanln(&option)
		fmt.Println()

		switch option {
		case "1":
			fmt.Println("Chat conversation! \n")
			printChat()

		case "2":
			fmt.Println("Chat has been backed up!")
			saveChat()

		case "0":
			fmt.Println("Server closed!")

		default:
			fmt.Println("Invalid option!")
		}
	}
}

func server() {
	listener, err := net.Listen(env.ConnectionType, env.ConnectionPort)
	if err != nil {
		fmt.Println("LISTENING ERROR:", err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("CONNECTION ERROR:", err)
			continue
		}

		users = append(users, connection)

		go handleClient(connection)
	}
}

func handleClient(clientConnection net.Conn) {
	var message msg.Message
	var active bool = true

	for {
		err := gob.NewDecoder(clientConnection).Decode(&message)
		if err != nil {
			if active {
				fmt.Println("DECODING ERROR:", err)
			}
			return
		}

		chat = append(chat, message)

		active = message.Text != msg.Logout

		if !active {
			for index, client := range users {
				if clientConnection == client {
					users = append(users[:index], users[index+1:]...)
				}
			}
		}

		for _, client := range users {
			if clientConnection != client {
				err := gob.NewEncoder(client).Encode(message)
				if err != nil {
					fmt.Println("ENCODING ERROR:", err)
					return
				}
			}
		}
	}
}

func printChat() {
	for _, message := range chat {
		fmt.Println(message.ToString())
	}
}

func saveChat() {
	file, err := os.Create("backup.txt")
	if err != nil {
		fmt.Println("BACKUP ERROR:", err)
		return
	}
	defer file.Close()

	for _, message := range chat {
		file.WriteString(message.ToString() + "\n")
	}
}
