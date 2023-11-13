package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "tu_paquete_proto" // Asegúrate de reemplazar "tu_paquete_proto" con el paquete real generado por el compilador de protobuf
)

type LogEntry struct {
	Action      string
	Sector      string
	Base        string
	NuevoValor  int32
	VectorClock []int32
}

type Nodo struct {
	mu            sync.Mutex
	log           []LogEntry
	relojVector   []int32
	consistente   bool
	ultimoChequeo time.Time
	esDominante   bool
}

type ServicioFulcrum struct {
	nodos         map[string]*Nodo
	mu            sync.Mutex
	nodoDominante *Nodo
	fulcrumNumero int // Agregar el número de fulcrum al servicio
}

func (s *ServicioFulcrum) AgregarBase(ctx context.Context, req *pb.RegistroSoldado, fulcrumNumero int) (*pb.Respuesta, error) {

	filePath := req.Sector + ".txt"
	fileContent := ""

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Si no existe, crear el archivo
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error al crear el archivo: %v", err)
			return nil, err
		}
		defer file.Close()
	}

	// Abrir el archivo existente o recién creado
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error al abrir el archivo: %v", err)
		return nil, err
	}
	defer file.Close()

	// Escribir en el archivo la información de la base
	fileContent += fmt.Sprintf("Base: %s", req.Base)

	// Si req.Cantidad está vacío, considerarlo como 0
	if req.Cantidad == "" {
		req.Cantidad = "0"
	}

	fileContent += fmt.Sprintf(", Valor: %s\n", req.Cantidad)

	if _, err := file.WriteString(fileContent); err != nil {
		// Manejar el error al escribir en el archivo
		log.Printf("Error al escribir en el archivo: %v", err)
		return nil, err
	}

	// Aumentar el valor en el vector de reloj en la posición fulcrumNumero-1
	s.mu.Lock()
	defer s.mu.Unlock()

	if fulcrumNumero <= 0 || fulcrumNumero > len(req.RelojVector) {
		return nil, fmt.Errorf("fulcrumNumero fuera de rango")
	}

	req.RelojVector[fulcrumNumero-1]++

	// Actualizar nodos y lógica consistente aquí si es necesario

	// Devolver la respuesta
	return &pb.Respuesta{Exito: true, Mensaje: "Base agregada exitosamente", RelojVector: req.RelojVector}, nil
}

func (s *ServicioFulcrum) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest, fulcrumNumero int) (*pb.Respuesta, error) {
	filePath := req.Sector + ".txt"

	// Abrir el archivo
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error al abrir el archivo: %v", err)
		return nil, err
	}
	defer file.Close()

	// Escanear el archivo línea por línea
	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Dividir la línea en partes utilizando "," como delimitador
		parts := strings.Split(line, ",")
		if len(parts) >= 2 {
			// Buscar la coincidencia de la baseAntigua y renombrar si es necesario
			if strings.TrimSpace(parts[1]) == req.BaseAntigua {
				// Cambiar el nombre de la base
				parts[1] = strings.TrimSpace(req.BaseNueva)
			}
		}

		// Reconstruir la línea y almacenarla
		lines = append(lines, strings.Join(parts, ","))
	}

	// Volver al inicio del archivo
	file.Seek(0, 0)

	// Truncar el archivo
	file.Truncate(0)

	// Escribir las líneas actualizadas de vuelta al archivo
	for _, line := range lines {
		fmt.Fprintln(file, line)
	}

	// Manejar errores de scanner
	if err := scanner.Err(); err != nil {
		log.Printf("Error al escanear el archivo: %v", err)
		return nil, err
	}

	// Aumentar el valor en el vector de reloj en la posición fulcrumNumero-1
	s.mu.Lock()
	defer s.mu.Unlock()

	if fulcrumNumero <= 0 || fulcrumNumero > len(req.RelojVector) {
		return nil, fmt.Errorf("fulcrumNumero fuera de rango")
	}

	req.RelojVector[fulcrumNumero-1]++

}

func (s *ServicioFulcrum) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	// Implementación del método ActualizarValor
	// ...
}

func (s *ServicioFulcrum) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	// Implementación del método BorrarBase
	// ...
}

func (s *ServicioFulcrum) VerificarConsistencia(ctx context.Context, req *pb.ConsistenciaRequest) (*pb.Respuesta, error) {
	// Implementación del método VerificarConsistencia
	// ...
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: tu_programa <fulcrum_numero>")
		os.Exit(1)
	}

	fulcrumArg := os.Args[1]
	fulcrumNumero, err := strconv.Atoi(fulcrumArg)
	if err != nil {
		fmt.Println("Error: el argumento debe ser un número entero")
		os.Exit(1)
	}

	servidor := grpc.NewServer()
	servicio := &ServicioFulcrum{
		nodos:         make(map[string]*Nodo),
		nodoDominante: nil,
		fulcrumNumero: fulcrumNumero, // Agregar el número de fulcrum al servicio
	}

	pb.RegisterFulcrumServer(servidor, servicio)

	direccion := ":50051"
	lis, err := net.Listen("tcp", direccion)
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	log.Printf("Servidor gRPC escuchando en %s para Fulcrum %d", direccion, fulcrumNumero)
	if err := servidor.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}

/*package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "tu_paquete_proto" // Reemplaza esto con el paquete real generado por el compilador de protobuffer
)

// Nodo representa la información y el reloj de un nodo.
type Nodo struct {
	mu           sync.Mutex
	informacion  map[string]int32
	relojVector  []int32
	consistente  bool
	ultimoChequeo time.Time
	esDominante   bool

	// Nueva lista para almacenar operaciones pendientes.
	operacionesPendientes []Operacion
	limiteOperaciones     int
}

// Operacion representa una operación pendiente.
type Operacion struct {
	Tipo     string
	Registro *pb.RegistroSoldado
	Renombrar *pb.RenombrarBaseRequest
	Actualizar *pb.ActualizarValorRequest
	Borrar   *pb.BorrarBaseRequest
}

// ServicioFulcrum implementa el servicio gRPC Fulcrum.
type ServicioFulcrum struct {
	nodos        map[string]*Nodo
	mu           sync.Mutex
	nodoDominante *Nodo
}

// AgregarBase implementa el método AgregarBase del servicio Fulcrum.
func (s *ServicioFulcrum) AgregarBase(ctx context.Context, req *pb.RegistroSoldado) (*pb.Respuesta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodo, err := s.obtenerNodo(req.Sector)
	if err != nil {
		return nil, err
	}

	nodo.mu.Lock()
	defer nodo.mu.Unlock()

	// Lógica de actualización del reloj de vectores y la información del nodo.
	if nodo.esDominante {
		// Agregar la operación a la lista de operaciones pendientes.
		nodo.operacionesPendientes = append(nodo.operacionesPendientes, Operacion{
			Tipo:     "AgregarBase",
			Registro: req,
		})
	} else {
		// Lógica para nodos no dominantes.
		// ...
	}

	// Verificar y aplicar cambios si se supera el límite de operaciones pendientes.
	s.aplicarCambiosPendientes(nodo)

	return &pb.Respuesta{Exito: true}, nil
}

// RenombrarBase implementa el método RenombrarBase del servicio Fulcrum.
func (s *ServicioFulcrum) RenombrarBase(ctx context.Context, req *pb.RenombrarBaseRequest) (*pb.Respuesta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodo, err := s.obtenerNodo(req.Sector)
	if err != nil {
		return nil, err
	}

	nodo.mu.Lock()
	defer nodo.mu.Unlock()

	// Lógica de actualización del reloj de vectores y la información del nodo.
	if nodo.esDominante {
		// Agregar la operación a la lista de operaciones pendientes.
		nodo.operacionesPendientes = append(nodo.operacionesPendientes, Operacion{
			Tipo:     "RenombrarBase",
			Renombrar: req,
		})
	} else {
		// Lógica para nodos no dominantes.
		// ...
	}

	// Verificar y aplicar cambios si se supera el límite de operaciones pendientes.
	s.aplicarCambiosPendientes(nodo)

	return &pb.Respuesta{Exito: true}, nil
}

// ActualizarValor implementa el método ActualizarValor del servicio Fulcrum.
func (s *ServicioFulcrum) ActualizarValor(ctx context.Context, req *pb.ActualizarValorRequest) (*pb.Respuesta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodo, err := s.obtenerNodo(req.Sector)
	if err != nil {
		return nil, err
	}

	nodo.mu.Lock()
	defer nodo.mu.Unlock()

	// Lógica de actualización del reloj de vectores y la información del nodo.
	if nodo.esDominante {
		// Agregar la operación a la lista de operaciones pendientes.
		nodo.operacionesPendientes = append(nodo.operacionesPendientes, Operacion{
			Tipo:       "ActualizarValor",
			Actualizar: req,
		})
	} else {
		// Lógica para nodos no dominantes.
		// ...
	}

	// Verificar y aplicar cambios si se supera el límite de operaciones pendientes.
	s.aplicarCambiosPendientes(nodo)

	return &pb.Respuesta{Exito: true}, nil
}

// BorrarBase implementa el método BorrarBase del servicio Fulcrum.
func (s *ServicioFulcrum) BorrarBase(ctx context.Context, req *pb.BorrarBaseRequest) (*pb.Respuesta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodo, err := s.obtenerNodo(req.Sector)
	if err != nil {
		return nil, err
	}

	nodo.mu.Lock()
	defer nodo.mu.Unlock()

	// Lógica de actualización del reloj de vectores y la información del nodo.
	if nodo.esDominante {
		// Agregar la operación a la lista de operaciones pendientes.
		nodo.operacionesPendientes = append(nodo.operacionesPendientes, Operacion{
			Tipo:   "BorrarBase",
			Borrar: req,
		})
	} else {
		// Lógica para nodos no dominantes.
		// ...
	}

	// Verificar y aplicar cambios si se supera el límite de operaciones pendientes.
	s.aplicarCambiosPendientes(nodo)

	return &pb.Respuesta{Exito: true}, nil
}

// VerificarConsistencia implementa el método VerificarConsistencia del servicio Fulcrum.
func (s *ServicioFulcrum) VerificarConsistencia(ctx context.Context, req *pb.ConsistenciaRequest) (*pb.Respuesta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodo, err := s.obtenerNodo(req.Sector)
	if err != nil {
		return nil, err
	}

	nodo.mu.Lock()
	defer nodo.mu.Unlock()

	// Lógica de verificación y corrección de consistencia eventual.
	if nodo.esDominante {
		// Lógica específica para el nodo dominante.
		// ...
	} else {
		// Lógica para nodos no dominantes.
		// ...
	}

	return &pb.Respuesta{Exito: true}, nil
}

// obtenerNodo obtiene el nodo correspondiente al sector.
func (s *ServicioFulcrum) obtenerNodo(sector string) (*Nodo, error) {
	nodo, ok := s.nodos[sector]
	if !ok {
		return nil, status.Error(codes.NotFound, "Sector no encontrado")
	}
	return nodo, nil
}

// aplicarCambiosPendientes verifica y aplica cambios si se supera el límite de operaciones pendientes.
func (s *ServicioFulcrum) aplicarCambiosPendientes(nodo *Nodo) {
	if len(nodo.operacionesPendientes) >= nodo.limiteOperaciones {
		// Lógica para aplicar cambios al archivo correspondiente.
		// Puedes iterar sobre la lista de operaciones pendientes y aplicar cada cambio.
		for _, operacion := range nodo.operacionesPendientes {
			switch operacion.Tipo {
			case "AgregarBase":
				// Lógica para aplicar la operación de agregar base.
				// ...
			case "RenombrarBase":
				// Lógica para aplicar la operación de renombrar base.
				// ...
			case "ActualizarValor":
				// Lógica para aplicar la operación de actualizar valor.
				// ...
			case "BorrarBase":
				// Lógica para aplicar la operación de borrar base.
				// ...
			}
		}

		// Limpiar la lista de operaciones pendientes después de aplicar los cambios.
		nodo.operacionesPendientes = nil
	}
}

func main() {
	s := grpc.NewServer()
	servicio := &ServicioFulcrum{
		nodos: make(map[string]*Nodo),
	}

	// Configura un nodo como dominante (por ejemplo, el primer nodo creado).
	servicio.nodoDominante = &Nodo{
		informacion:           make(map[string]int32),
		relojVector:           make([]int32, 3), // Ajusta el tamaño según tus necesidades.
		esDominante:           true,
		operacionesPendientes: make([]Operacion, 0),
		limiteOperaciones:     10, // Ajusta el límite según tus necesidades.
	}

	// Registra el servicio Fulcrum en el servidor gRPC.
	pb.RegisterFulcrumServer(s, servicio)

	// Inicia el servidor en el puerto 50051.
	log.Println("Iniciando servidor en el puerto 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
*/
