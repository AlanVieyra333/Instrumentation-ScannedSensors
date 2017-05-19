package main

import (
	"fmt"
	"server"
	//"github.com/tarm/serial"
)

//*
func main() {
	var port string

	fmt.Println("Servidor con conexión serial.\n")
	fmt.Println("Que puerto desea usar para la conexión serial? (Ej. COM1)")
	fmt.Scanf("%s", &port)

	err := server.Connect(port, uint(700))
	if err != nil {
		fmt.Println("Error al conectar con el puerto serial.")
	} else {
		defer server.CloseConnection()
		server.ServerTurnOn()
	}
} //*/

/*
func main() {
	fmt.Println("Ini")
	c := &serial.Config{Name: "COM2", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("Reading...")
	buf := make([]byte, 128)
	n, err := s.Read(buf)
	if err != nil {
		fmt.Print(err.Error())
	}
	//log.Printf("Read:%q", buf[:n])
	fmt.Printf("Read %v bytes: '%q'\n", n, buf[:n])

	n, err = s.Write([]byte("test"))
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("writed %v bytes\n", n)
} //*/
