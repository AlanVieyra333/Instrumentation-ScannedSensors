package server

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

var f io.ReadWriteCloser
var pulseH bool
var t0, t1 time.Time

func Connect(port string, timeout uint) error {
	var err error
	pulseH = false

	serialPort := flag.String("serialPort", port, "serial port to test (/dev/ttyUSB0, etc)")
	baud := flag.Uint("baud", 9600, "Baud rate")
	//txData := flag.String("txdata", "31", "data to send in hex format (01ab238b)")
	even := flag.Bool("even", false, "enable even parity")
	odd := flag.Bool("odd", false, "enable odd parity")
	rs485 := flag.Bool("rs485", false, "enable RS485 RTS for direction control")
	rs485HighDuringSend := flag.Bool("rs485_high_during_send", false, "RTS signal should be high during send")
	rs485HighAfterSend := flag.Bool("rs485_high_after_send", false, "RTS signal should be high after send")
	stopbits := flag.Uint("stopbits", 1, "Stop bits")
	databits := flag.Uint("databits", 8, "Data bits")
	chartimeout := flag.Uint("chartimeout", timeout, "Inter Character timeout (ms)")
	minread := flag.Uint("minread", 1, "Minimum read count")
	//rx := flag.Bool("rx", true, "Read data received")

	flag.Parse()

	if *serialPort == "" {
		return &errorString{"Must specify port."}
	}

	if *even && *odd {
		return &errorString{"Can't specify both even and odd parity."}
	}

	parity := serial.PARITY_NONE

	if *even {
		parity = serial.PARITY_EVEN
	} else if *odd {
		parity = serial.PARITY_ODD
	}

	options := serial.OpenOptions{
		PortName:               *serialPort,
		BaudRate:               *baud,
		DataBits:               *databits,
		StopBits:               *stopbits,
		MinimumReadSize:        *minread,
		InterCharacterTimeout:  *chartimeout,
		ParityMode:             parity,
		Rs485Enable:            *rs485,
		Rs485RtsHighDuringSend: *rs485HighDuringSend,
		Rs485RtsHighAfterSend:  *rs485HighAfterSend,
	}

	f, err = serial.Open(options)
	if err != nil {
		return &errorString{"Error opening serial port: " + err.Error()}
	}

	return nil
}

func GetData() (int, int, int, error) {
	var degrees, motor, bpm int
	var data float32

	data, err := getDigitalValue("30")
	if err != nil {
		return -1, -1, -1, &errorString{err.Error()}
	}
	degrees = int(data * 100)
	fmt.Printf("Temperatura: %v grados.\n", degrees)
	//	------------------------------------------------------------------------
	data, err = getDigitalValue("31")
	if err != nil {
		return degrees, -1, -1, &errorString{err.Error()}
	}
	motor = int(data * 100)
	fmt.Printf("Luminosidad: %v Luxes.\n", motor)
	//	------------------------------------------------------------------------
	data, err = getDigitalValue("32")
	if err != nil {
		return degrees, motor, -1, &errorString{err.Error()}
	}
	bpmTmp := int(data * 1)
	bpm = 0
	//fmt.Printf("DATA OK\n")
	if bpmTmp == 1 { // Pulse higth
		if !pulseH {
			pulseH = true
			if t0.IsZero() {
				t0 = time.Now()
			} else {
				t1 = time.Now()

				bpm = round(60 / t1.Sub(t0).Seconds())

				t0 = time.Now()
			}
		}
	} else {
		if pulseH {
			pulseH = false
		}
	}
	fmt.Printf("Latidos: %v lpm.\n", bpm)
	return degrees, motor, bpm, nil
}

func getDigitalValue(msg string) (float32, error) {
	//fmt.Printf("recives\n")
	var txData_ []byte
	buf := make([]byte, 2)
	var err error
	var n int

	txData_, err = hex.DecodeString(msg)
	if err != nil {
		return 0.0, &errorString{"Error decoding hex data:\n" + err.Error()}
	}

	n, err = f.Write(txData_)
	if err != nil {
		return 0.0, &errorString{"Error writing to serial port:\n" + err.Error()}
	}
	//fmt.Printf("Wrote %v bytes: '%s'\n", n, hex.EncodeToString(txData_))

	bufAux := make([]byte, 2)
	tx1 := time.Now()
	for i := 0; i < 2; {
		n, err = f.Read(bufAux)
		if err != nil {
			if err != io.EOF {
				return 0.0, &errorString{"Error reading to serial port:\n" + err.Error()}
			}
		}
		i += n
		if n == 1 && i == 0 {
			buf[0] = bufAux[0]
		} else if n == 1 && i == 1 {
			buf[1] = bufAux[0]
		} else if n == 2 {
			buf[0] = bufAux[0]
			buf[1] = bufAux[1]
		}

		tx2 := time.Now()
		if tx2.Sub(tx1).Seconds() > 1 {
			n = 0
			break
		}
	}
	//fmt.Printf("Read 2 bytes: '%s'\n", hex.EncodeToString(buf))

	if n == 0 {
		return 0.0, &errorString{"No hay conexión al puerto serial."}
	}

	return float32(float32(binary.BigEndian.Uint16(buf)) / float32(1023)), nil
}

func CloseConnection() {
	f.Close()
}

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
