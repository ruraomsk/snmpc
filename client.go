package snmpc

import (
	"fmt"
	"net"
	"time"

	gtrap "github.com/gosnmp/gosnmp"
	"github.com/ruraomsk/potop/agent"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/journal"
	"github.com/ruraomsk/potop/logger"
	"github.com/ruraomsk/potop/setup"
)

var connect bool
var hs SnmpHard
var autonom bool
var controlPhase = false
var controlPhaseNumber = 0
var exitSpesial = false
var timeToLPU = 90

func Connect() bool {
	return connect
}
func SetAutonom(value bool) {
	autonom = value
}
func GetAutonom() bool {
	return autonom
}
func getDuration() time.Duration {
	return time.Duration(timeToLPU) * time.Second
}

type ctrlPhase struct {
	state bool
	phase int
}

var live chan any
var counter chan ctrlPhase

func counterPhases() {
	var timer *time.Timer
	timer = time.NewTimer(time.Hour)
	timer.Stop()
	controlPhase = false
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if controlPhase {
				if hardware.GetVPU() {
					controlPhase = false
					timer.Stop()
					continue
				}
			}
		case s := <-counter:
			if s.state {
				// Нормальное завершение
				if controlPhase {
					timer.Stop()
					controlPhase = false
				}
			} else {
				// Start отсчета
				if controlPhase {
					timer.Stop()
				}
				timer = time.NewTimer(getDuration())
				controlPhase = true
				controlPhaseNumber = s.phase
			}
		case <-timer.C:
			hardware.CommandToKDM(1, 0)
			controlPhase = false
		}
	}

}
func waiter() {
	var timer *time.Timer
	for {
		logger.Info.Print("Нет управления от STCIP")
		<-live
		hardware.SetWork <- 1
		timer = time.NewTimer(getDuration())
		logger.Info.Print("Есть управление от STCIP")
		connect = true
	loop:
		for {
			select {
			case <-timer.C:
				if !autonom {
					hardware.SetWork <- 0
					hardware.CommandToKDM(0, 1)
					break loop
				}
			case <-live:
				timer.Stop()
				if !hardware.StateHardware.GetCentral() {
					hardware.SetWork <- 1
				}
				timer = time.NewTimer(getDuration())
			}
		}
		connect = false
		timer.Stop()
		logger.Error.Print("Потеряно управление от STCIP")
	}

}
func senderTrapsSnmp() {

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	hsold = hs
	gtrap.Default.Target = setup.Set.SNMP.Host
	gtrap.Default.Port = uint16(setup.Set.SNMP.Port)
	gtrap.Default.Version = gtrap.Version2c
	gtrap.Default.Community = "public"
	// gtrap.Default.Logger = gtrap.NewLogger(log.New(os.Stdout, "", 0))
	connect := false
	journal.SendMessage(journal.LevelSMNP, fmt.Sprintln("TRAP NOT READY"))
	for !connect {
		err := gtrap.Default.Connect()
		if err != nil {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("TRAP err:%v", err))
			time.Sleep(time.Second)
		} else {
			connect = true
		}
	}
	defer gtrap.Default.Conn.Close()
	journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("TRAP READY to %s:%d", gtrap.Default.Target, gtrap.Default.Port))
	for {
		trap := gtrap.SnmpTrap{
			Variables:    maketrapSnmp(),
			Enterprise:   ".1.3.6.1.4.1.1618",
			AgentAddress: "127.0.0.1",
			GenericTrap:  0,
			SpecificTrap: 0,
			Timestamp:    uint(time.Now().Unix()),
		}
		if len(trap.Variables) > 0 {
			for _, v := range trap.Variables {
				// logger.Debug.Printf("trap %v", v)
				journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("TRAP %v %v", v.Name, v.Value))
			}
			_, err := gtrap.Default.SendTrap(trap)
			if err != nil {
				logger.Error.Printf("SendTrap() err: %v", err)
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func Start() {
	connect = false
	autonom = false
	for {
		if len(hardware.StateHardware.Status) > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	hs.Update()
	live = make(chan any)
	counter = make(chan ctrlPhase)
	go counterPhases()
	agent := agent.NewAgent()
	if !setup.Set.SNMP.Trap {
		go senderTrapsSnmp()
	}
	// Set the read-only and read-write communities
	agent.SetCommunities("public", "private")
	registator(agent)
	addr, err := net.ResolveUDPAddr("udp", ":161")
	if err != nil {
		logger.Error.Printf("SNMP %v", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Error.Printf("SNMP %v", err)
		return
	}
	connect = true
	go waiter()
	journal.Setter <- journal.SetLevel{Level: 2, Double: true}
	logger.Info.Printf("Starting as SNMP Server on: %s\n", setup.ExtSet.SNMP.Listen)
	for {
		buffer := make([]byte, 10240)
		n, source, err := conn.ReadFrom(buffer)
		if err != nil {
			logger.Error.Printf("SNMP %v", err)
			continue
		}
		// logger.Debug.Printf("req %v", buffer[:n])
		buffer, err = agent.ProcessDatagram(buffer[:n])
		if err != nil {
			logger.Error.Printf("SNMP %v", err)
			continue
		}
		// logger.Debug.Printf("rep %v", buffer)
		_, err = conn.WriteTo(buffer, source)
		if err != nil {
			logger.Error.Printf("SNMP %v", err)
			continue
		}
		live <- 1
	}
}
