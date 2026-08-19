package snmpc

import (
	"encoding/binary"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ruraomsk/potop/agent"
	"github.com/ruraomsk/potop/hardware"
)

type SnmpHard struct {
	StateHard hardware.StateHard
	LastVPU   bool
}

func (s *SnmpHard) Update() {
	s.StateHard = hardware.GetStateHard()
	if s.LastVPU && !s.StateHard.VPU {
		exitSpesial = true
	}
	s.LastVPU = s.StateHard.VPU
	if len(s.StateHard.Status) == 0 {
		s.StateHard.Status = make([]int, 4)
	}
}

func (s *SnmpHard) GetStatus() int {
	if !s.StateHard.Connect {
		return 0
	}
	if s.StateHard.Plan == 26 || s.StateHard.Phase == 0 {
		// КК
		return 6 //????
	}
	if s.StateHard.Plan == 27 || s.StateHard.Phase == 33 {
		// ЖМ
		return 4
	}
	if s.StateHard.Plan == 25 || s.StateHard.Phase == 34 {
		// ОС
		return 3
	}

	if s.StateHard.Connect {
		return 1
	}

	return 2
}
func (s *SnmpHard) GetSource() int {
	if s.StateHard.VPU {
		return 12
	}
	if s.StateHard.Source == 0 {
		return 11
	}
	if s.StateHard.Source == 1 {
		return 6
	}
	if s.StateHard.Source == 2 {
		return 1
	}
	if s.StateHard.Source == 3 {
		return 12
	}
	return 10
}
func (s *SnmpHard) GetState() int {
	if autonom {
		return 6
	}
	if s.StateHard.VPU {
		return 12
	}
	if s.StateHard.Source == 0 {
		return 8
	}
	if s.StateHard.Source == 1 {
		return 11
	}
	if s.StateHard.Source == 2 {
		return 9
	}
	if s.StateHard.Source == 3 {
		return 12
	}
	return 0
}
func (s *SnmpHard) GetPlanInteger() int {
	if s.StateHard.VPU {
		return 15
	}
	if GetAutonom() {
		return 0
	}
	// if s.StateHard.Plan == 26 {
	// 	// КК
	// 	return agent.Unsigned32(17) //????
	// }
	if s.StateHard.Plan == 27 || s.StateHard.Phase == 33 {
		// ЖМ
		return 17
	}
	if s.StateHard.Plan == 25 || s.StateHard.Phase == 34 {
		// ОС
		return 18
	}
	if controlPhase && s.StateHard.Phase == 255 {
		return 16
	}
	if controlPhase && s.StateHard.Phase == controlPhaseNumber {
		return 16
	}
	//Планы в контроллере
	if s.StateHard.Plan >= 1 && s.StateHard.Plan <= 24 {
		return s.StateHard.Plan
	}

	// }
	if s.StateHard.Plan == 0 {
		//Ручное управление
		return 15
	}
	return 16
}
func (s *SnmpHard) GetPlan() agent.Unsigned32 {
	return agent.Unsigned32(s.GetPlanInteger())
}
func (s *SnmpHard) GetPhaseInteger() int {
	if s.StateHard.Phase == 0 {
		// КК
		return 0
	}
	if s.StateHard.Plan == 27 || s.StateHard.Phase == 33 {
		// ЖМ
		return 0
	}
	if s.StateHard.Plan == 25 || s.StateHard.Phase == 34 {
		// ОС
		return 0
	}

	if s.StateHard.Plan == 35 {
		//Работаем по Utopia
		return 14
	}
	if s.StateHard.Phase == 255 {
		//Промтакт
		if exitSpesial {
			return 2
		}
		if s.StateHard.LastPhase <= 0 {
			return 2
		}
		if s.StateHard.LastPhase < 8 {
			return s.StateHard.LastPhase + 1
		}
		return 8
	}
	if exitSpesial {
		exitSpesial = false
	}
	if s.StateHard.Phase < 8 {
		return s.StateHard.Phase + 1
	}
	return 8
}
func (s *SnmpHard) GetPhase() agent.Unsigned32 {
	return agent.Unsigned32(s.GetPhaseInteger())
}
func (s *SnmpHard) GetPromtaktInteger() int {
	if s.StateHard.Phase == 255 {
		//Промтакт
		return 1
	}
	return 0
}
func (s *SnmpHard) GetPromtakt() agent.Unsigned32 {
	return agent.Unsigned32(s.GetPromtaktInteger())
}
func getUint32(value any) agent.Unsigned32 {
	v, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
	return agent.Unsigned32(v)

}
func getOctetString(value any) string { //agent.OctetString {
	buf := []byte{0, 0, 0, 0}
	v, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	var res strings.Builder
	for _, v := range buf {
		res.WriteString(fmt.Sprintf(" %02X", v))
	}
	return res.String()
	//agent.OctetString(buf)
}
func getCount32(value any) agent.Counter32 {
	v, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
	return agent.Counter32(v)

}

func (s *SnmpHard) GetSignalGroup() string {
	res := make([]byte, 0)
	var b byte = 0

	for _, v := range hs.StateHard.StatusDirs {
		switch v {
		case 0:
			//   OFF = 0, //все сигналы выключены+
			b = 14
		case 1:
			//   DEACTIV_YELLOW=1, //направление перешло в неактивное состояние, желтый после зеленого+
			b = 6
		case 2:
			//   DEACTIV_RED=2, //направление перешло в неактивное состояние, красный+
			b = 8
		case 3:
			//   ACTIV_RED=3, //направление перешло в активное состояние, красный+
			b = 8
		case 4:
			//   ACTIV_REDYELLOW=4, //направление перешло в активное состояние, красный c желтым +
			b = 0
		case 5:
			//   ACTIV_GREEN=5, //направление перешло в активное состояние, зеленый+
			b = 2
		case 6:
			//   UNCHANGE_GREEN=6, //направление не меняло свое состояние, зеленый+
			b = 3
		case 7:
			//   UNCHANGE_RED=7, //направление не меняло свое состояние, красный+
			b = 8
		case 8:
			//   GREEN_BLINK=8, //зеленый мигающий сигнал+
			b = 5
		case 9:
			//   ZM_YELLOW_BLINK=9, //желтый мигающий в режиме ЖМ+
			b = 6
		case 10:
			//   OS_OFF=10,	//сигналы выключены в режиме ОС+
			b = 13
		case 11:
			//   UNUSED=11 //неиспользуемое направление+
			b = 13
			continue
		default:
			b = 14
		}
		res = append(res, b)
		// res += fmt.Sprintf("%02X-", b)
	}
	// res = strings.TrimSuffix(res, "-")
	return string(res)

}

func SetSystemTimeFromUnix(unixTime int64) error {
	t := time.Unix(unixTime, 0)

	// Формат для команды date: MMDDhhmmYYYY.ss
	dateStr := t.Format("010215042006.05")

	// Выполняем команду
	cmd := exec.Command("date", dateStr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("не удалось установить время: %w", err)
	}

	return nil
}
