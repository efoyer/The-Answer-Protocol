package main

import (
	"bufio"
	"log"
	"net"
	"fmt"
)

func main(){ // Point d’entrée du programme.
	ln, err := net.Listen("tcp", ":8080") // Ouvre l’écoute TCP sur le port 8080.
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close() // Ferme l’écoute quand main se termine.
	for { // Répète sans limite pour accepter de nouveaux clients.
		// Attend un client et récupère sa connexion.
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	scan := bufio.NewScanner(c)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}
	err := scan.Err()
	if err != nil {
		fmt.Println("Erreur: ", err)
	}
}