package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-3/proto"
	"google.golang.org/grpc"
)

// Struct para poder hacer un objeto de base dependiendo si es que esta creado.
type base struct {
	nombre      string //Nombre del planeta manejado (registro)
	relojx      int    //Dimension X del reloj de vector
	relojy      int    //Dimension Y del reloj de vector
	relojz      int    //Dimension Z del reloj de vector
	lastfulcrum string //ip del ultimo fulcrum consultado para este planeta
}

// Lista de structs que almacenará de manera eficiente los bases.
var bases []base
var direccionBroker = "dist053:50051"
var direccionFulcrum = ""

// Constructor para el planeta, cosa de poder almacenar en memoria la info de los planetas manejados por la consola del informante.
func Cbase(name string, x int, y int, z int, ip string) (basedata base) {
	basedata = base{
		nombre:      name,
		relojx:      x,
		relojy:      y,
		relojz:      z,
		lastfulcrum: ip,
	}
	return
}

// Funcion ejecutada por gRPC para enviar el mensaje
func Solicitud(serviceClient pb.BrokerClient, msg string) string {
	res, err := serviceClient.RedirectInformant(context.Background(), &pb.InformantRequest{
		Command: msg,
	})
	if err != nil {
		panic("Mensaje no pudo ser creado ni enviado: " + err.Error())
	}
	//fmt.Println(res.Body)
	return res.Address
}

// Funcion que toma la IP del Broker, se conecta y realiza el envio del mensaje. Retorna la IP del Fulcrum
func enviarMsg(ip string, msg string) (answer string) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())

	if err != nil {
		panic("No se puede conectar al servidor " + err.Error())
	}

	serviceClient := pb.NewBrokerClient(conn)

	answer = Solicitud(serviceClient, msg)

	defer conn.Close()

	return
}

func sendToFulcrum(ip string, commandParts []string) string {
	// Connect to Fulcrum
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		panic("No se puede conectar al Fulcrum " + err.Error())
	}
	serviceClient := pb.NewFulcrumClient(conn)

    // Initialize a CommandRequest with default values
    req := &pb.CommandRequest{
        Action:  commandParts[0],
        Sector:  commandParts[1],
        Base:    commandParts[2],
        NewBase: "",
        Value:   0,
    }

    // Handle different actions
    switch commandParts[0] {
    case "RenombrarBase":
        if len(commandParts) > 3 {
            req.NewBase = commandParts[3]
        }
    case "AgregarBase", "ActualizarValor":
        if len(commandParts) > 3 {
            value, err := strconv.Atoi(commandParts[3])
            if err != nil {
                panic("Invalid value: " + err.Error())
            }
            req.Value = int32(value)
        }
    case "BorrarBase":
        // Do nothing, NewBase and Value should remain empty
    default:
        panic("Invalid action: " + commandParts[0])
    }

    res, err := serviceClient.ApplyCommand(context.Background(), req)
    if err != nil {
        panic("Message could not be created or sent: " + err.Error())
    }

	strClock := make([]string, len(res.VectorClock))
    for i, num := range res.VectorClock {
        strClock[i] = strconv.Itoa(int(num))
    }

    return strings.Join(strClock, ",")
}

// Procesa los comandos del usuario (Consulta a broker, luego a Fulcrum).
func processMsg(command string) {
	//Comando = ["AgregarBase sector base 100"]
	var comando = strings.Split(command, " ")

	//Se recibe la ip para el fulcrum
	respuesta := enviarMsg(direccionBroker, command)
	fmt.Println("[*] Ip recibida desde el Broker:")
	fmt.Println(respuesta)
	direccionFulcrum = respuesta

	//Se consulta al Fulcrum
	fmt.Println("[*] Ejecutando consulta al servidor fulcrum...")
	respuestaFulcrum := sendToFulcrum(respuesta, comando)
	fmt.Println("[*] Respuesta recibida!, datos:")
	fmt.Println(respuestaFulcrum)

	//Se analiza si no hay error
	data := strings.Split(respuestaFulcrum, ",")
	if len(data) == 3 {
		//Se recibieron los valores del reloj, se verifica consistencia y se actualiza data en struct del planeta.
		dataX, _ := strconv.Atoi(data[0])
		dataY, _ := strconv.Atoi(data[1])
		dataZ, _ := strconv.Atoi(data[2])
		flag := 1
		for i := range bases {
			if bases[i].nombre == comando[1] {
				if (dataX >= bases[i].relojx) && (dataY >= bases[i].relojy) && (dataZ >= bases[i].relojz) {
					bases[i].relojx = dataX
					bases[i].relojy = dataY
					bases[i].relojz = dataZ
					bases[i].lastfulcrum = direccionFulcrum
					flag = 0
					fmt.Println("\n[*] Sin Error de consistencia!")
					break
				} else {
					fmt.Println("[*] Error de consistencia!")
					flag = 0
					break
				}
			}
		}
		if flag == 1 {
			//Quiere decir que no se maneja info del planeta y el archivo fue creado.
			bases = append(bases, Cbase(comando[1], dataX, dataY, dataZ, direccionFulcrum))
		}
	} else {
		fmt.Println("Process message: Error en la respuesta del servidor Fulcrum")
	}
}

func scanMsg() (mensaje string) {
	scanner := bufio.NewScanner(os.Stdin)
	var PromptC = ""
	fmt.Println("Escriba el comando a ejecutar (0 para cerrar programa)")
	fmt.Println("Recuerde ser consistente con mayúsculas y minúsculas para los comandos")
	if scanner.Scan() {
		PromptC = scanner.Text()
	}
	mensaje = PromptC
	return
}

func main() {
	mensaje := "-1"
	for mensaje != "0" {
		mensaje := scanMsg()
		if mensaje == "0" {
			break
		}
		processMsg(mensaje)
	}
}
