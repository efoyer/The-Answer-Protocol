package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main(){ // Point d’entrée du programme.
	ln, err := net.Listen("tcp", ":8080") // Ouvre l’écoute TCP sur le port 8080.
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close() // Ferme l’écoute quand main se termine.
	register_name := make(map[string]bool)
	for { // Répète sans limite pour accepter de nouveaux clients.
		// Attend un client et récupère sa connexion.
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn, register_name)
	}
}

func handleConnection(c net.Conn, reg map[string]bool) {
	defer c.Close()
	scan := bufio.NewScanner(c)
	_, err := fmt.Fprintln( c, "OK hello proto=1")
	if err != nil {
		log.Print(err)
		return
	}
	name := ""
	// Assure le netoyage du registre en cas de KeyInterrupt
	defer func(){
		if name != ""{
			delete(reg, name)
		}
	}()

	for scan.Scan(){
		ln := scan.Text()
		command, other, _ := strings.Cut(ln, " ")
		subcommand, _, _:= strings.Cut(other, " ")
		subcommand = strings.ToUpper(subcommand)
		command = strings.ToUpper(command)

		val := ""
		leave := false
		if name == "" && command != "CONNECT" && command != "QUIT"{
			val = fmt.Sprintf("The %s command is only accessible when you are logged in.", command)
		}else{
		switch command{
			case "QUIT":
				val = "OK bye"
				leave = true
			case "CONNECT":
				n, err := connectCommand(other, reg)
				if err != nil {
					val = err.Error()
					break
				}
				if name != ""{
					delete(reg, name)
				}
				name = n
				val = "OK connected"
			case "LOOK":
				val = "OK <json-data> (a faire)"
				
			case "MOVE":
				switch other{
					case "north":
						val = "OK room=<room-id> (a faire)"
					default:
						val = "OK (or not) room <room-id>"
				}
			case "CHAT":
				val = "OK"
			case "WHO":
				val = "OK players=<count>"
			case "GROUP":
				switch subcommand{
					case "CREATE":
						val = "OK group=<group-id> (cmd= GROUP CREATE)"
					case "INVITE":
						val = "OK (cmd=GROUP INVITE)"
					case "JOIN":
						val = "OK group=<group-id> (cmd=GROUP JOIN)"
					case "LEAVE":
						val = "OK (cmd=GROUP LEAVE)"
					default:
						val = "Unknown GROUP prefix"
				}
			case "TAKE":
				val = "OK taken=<item-id>"
			case "DROP":
				val = "OK dropped=<item-id>"
			case "INVENTORY":
				val = "OK [<json-item-list>]"
			case "TALK":
				val = "OK <dialogue>"
			case "ATTACK":
				val = "OK <combat-result>"
			case "STATUS":
				val = "oK {hp, max_hp, status}"
			case "QUEST":
				val = "OK <quest-data>"
			case "QUESTS":
				val = "OK <quest-list>"
			default:
				val = "Unknown command"

			}
		}
		_, err := fmt.Fprintln(c, val)
		if err != nil {
			fmt.Println("Erreur: ", err)
			return
		}
		fmt.Printf("key=%q | value=%q | name=%q\n", command, other, name)
		if leave {return}
	}
	err = scan.Err()
	if err != nil {
		fmt.Println("Erreur: ", err)
		return
	}
}

func connectCommand(name string, reg map[string]bool) (string, error) {
	word := strings.Fields(name)
	if len(word) > 1 {
		return "", fmt.Errorf("There are too many arguments")
	} else if len(word) == 0{
		return "", fmt.Errorf("Empty name")
	}else if val, ok := reg[word[0]]; ok && val{
		return "", fmt.Errorf("ERR 201 NAME_IN_USE")
	} else {
		reg[word[0]] = true
		return word[0], nil
	}
}