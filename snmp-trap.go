package snmpc

import (
	gtrap "github.com/gosnmp/gosnmp"
	"github.com/ruraomsk/potop/setup"
)

var hsold SnmpHard
var hn SnmpHard

func maketrapSnmp() []gtrap.SnmpPDU {
	hn.Update()
	result := make([]gtrap.SnmpPDU, 0)
	if hn.GetStatus() != hsold.GetStatus() {
		//Изменился статус  	// "1.3.6.1.4.1.1618.3.6.2.1.2.0":
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.6.2.1.2.0", Type: gtrap.Integer, Value: int(hn.GetStatus())})
	}
	if hn.GetState() != hsold.GetState() {
		//Изменился статус  	"1.3.6.1.4.1.1618.3.6.2.2.2.0":
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.6.2.2.2.0", Type: gtrap.Integer, Value: int(hn.GetState())})
	}
	if hn.GetPhase() != hsold.GetPhase() {
		//Сменилась фаза
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.7.2.11.2.0", Type: gtrap.Integer, Value: int(hn.GetPhaseInteger())})
	}
	if hn.GetPlan() != hsold.GetPlan() {
		//Сменился план 1, 3, 6, 1, 4, 1, 1618, 3, 5, 2, 1, 7, 0
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.5.2.1.7.0", Type: gtrap.Integer, Value: int(hn.GetPlanInteger())})
	}
	if hn.GetSource() != hsold.GetSource() {
		//Сменился источник "1.3.6.1.4.1.1618.3.7.2.1.3.0"
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.7.2.1.3.0", Type: gtrap.Integer, Value: int(hn.GetSource())})
	}
	if hn.GetSignalGroup() != hsold.GetSignalGroup() {
		//Сминились сигнальные группы 1.3.6.1.4.1.1618.3.5.2.1.6.0
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.5.2.1.6.0", Type: gtrap.OctetString, Value: hn.GetSignalGroup()})
	}
	if CodeError(byte(hn.StateHard.Status[0])) != CodeError(byte(hsold.StateHard.Status[0])) {
		//Сменились тревоги 1.3.6.1.4.1.1618.3.1.2.2.2.0
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.1.2.2.2.0", Type: gtrap.Integer, Value: int(CodeError(byte(hn.StateHard.Status[0])))})
	}
	if hn.GetPromtaktInteger() != hsold.GetPromtaktInteger() {
		// Сменился промтакт "1.3.6.1.4.1.1618.3.7.2.11.2.1.0":
		result = append(result, gtrap.SnmpPDU{Name: "1.3.6.1.4.1.1618.3.7.2.11.2.1.0", Type: gtrap.Integer, Value: hn.GetPromtaktInteger()})
	}
	hsold = hn
	if len(result) > 0 {
		r := make([]gtrap.SnmpPDU, 0)
		r = append(r, gtrap.SnmpPDU{Name: "1.3.6.1.2.1.1.5.0", Type: gtrap.Integer, Value: setup.Set.SNMP.Uid})
		result = append(r, result...)
	}
	return result
}
