package snmpc

import (
	"fmt"
	"time"

	"github.com/ruraomsk/potop/agent"
	"github.com/ruraomsk/potop/asn1"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/journal"
	"github.com/ruraomsk/potop/setup"
)

func registator(agent *agent.Agent) {
	// 1.3.6.1.4.1.1618.3.7.2.10.1
	// По этому OID отвечать находится ли контроллер в локальном режиме ("1") или нет ("0").
	// Т.е. выполняет ли он команды из протокола или нет.
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 10, 1, 0},
		func(oid asn1.Oid) (any, error) {
			r := 0
			if GetAutonom() {
				r = 1
			}
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Локальный режим %d", oid, r))
			return r, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Локальный режим %v", oid, value))
			if GetAutonom() {
				return 1, nil
			}
			return 0, nil
		})

	// 1.3.6.1.4.1.1618.3.5.2.1.5
	// Количество сигнальных групп
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 5, 2, 1, 5, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос числа сигнальных групп %d", oid, hardware.CountSignalGroup()))
			return hardware.CountSignalGroup(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v апрос числа сигнальных групп %v", oid, value))
			return hardware.CountSignalGroup(), nil
		})
	// "1.3.6.1.2.1.1.5.0":
	// Запрос ID
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 2, 1, 1, 5, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос ID устройства %d", oid, setup.Set.SNMP.Uid))
			return setup.Set.SNMP.Uid, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос ID устройства %v", oid, value))
			return setup.Set.SNMP.Uid, nil
		})

	// "1.3.6.1.4.1.1618.3.6.2.1.2.0":
	//indicates the major state of a traffic controller
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 6, 2, 1, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v indicates the major state of a traffic controller %d", oid, hs.GetStatus()))
			return hs.GetStatus(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v indicates the major state of a traffic controller %v", oid, value))
			return hs.GetStatus(), nil
		})
	// case "1.3.6.1.4.1.1618.3.6.2.2.2.0":
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 6, 2, 2, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Состояние устройства %d", oid, hs.GetState()))
			return hs.GetState(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Состояние устройства %v", oid, value))
			return hs.GetState(), nil
		})
	// case "1.3.6.1.4.1.1618.3.5.2.1.7.0":
	// 	//Номер плана
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 5, 2, 1, 7, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v номер плана %d", oid, hs.GetPlan()))
			return hs.GetPlan(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v номер плана %v", oid, value))
			return hs.GetPlan(), nil
		})
	// case "1.3.6.1.4.1.1618.3.7.2.1.3.0":
	// 	//Источник плана
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 1, 3, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Источник плана %d", oid, hs.GetSource()))
			return hs.GetSource(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Источник плана %v", oid, value))
			return hs.GetSource(), nil
		})
	// case "1.3.6.1.4.1.1618.3.7.2.11.2.0":
	//
	//	//Номер фазы
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 11, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Номер фазы %d", oid, hs.GetPhase()))
			return hs.GetPhase(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Номер фазы %v", oid, value))
			return hs.GetPhase(), nil
		})
	// case "1.3.6.1.4.1.1618.3.5.2.1.3.0":
	//
	//	//Номер фазы 2
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 5, 2, 1, 3, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Номер фазы %d", oid, hs.GetPhase()))
			return hs.GetPhase(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Номер фазы %v", oid, value))
			return hs.GetPhase(), nil
		})
	// case "1.3.6.1.4.1.1618.3.7.2.11.2.1.0":
	//
	//	Наличие промтакта
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 11, 2, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Наличие промтакта %d", oid, int(hs.GetPromtakt())))
			return int(hs.GetPromtakt()), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Наличие промтакта %v", oid, value))
			return int(hs.GetPromtakt()), nil
		})
	// case "1.3.6.1.4.1.1618.3.5.2.1.6.0":
	//
	//	//Состояние сигнальных групп
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 5, 2, 1, 6, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Состояние сигнальных групп %v", oid, hs.GetSignalGroup()))
			return hs.GetSignalGroup(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Состояние сигнальных групп %v", oid, value))
			return hs.GetSignalGroup(), nil
		})
	// case "1.3.6.1.4.1.1618.3.1.2.2.2.0":
	//	//Тревоги
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 1, 2, 2, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Тревоги %v", oid, CodeError(byte(hs.StateHard.Status[0]))))
			return getOctetString(CodeError(byte(hs.StateHard.Status[0]))), nil
			// return getUint32(CodeError(byte(hs.StateHard.Status[0]))), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Тревоги %v", oid, value))
			return getOctetString(CodeError(byte(hs.StateHard.Status[0]))), nil
			// return getUint32(CodeError(byte(hs.StateHard.Status[0]))), nil
		})
	// case "1.3.6.1.4.1.1618.3.1.2.2.3.0":
	//	//Тревоги
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 1, 2, 2, 3, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			r := getOctetString(0)
			if len(hs.StateHard.Status) > 2 {
				r = getOctetString(hs.StateHard.Status[1])
			}
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Тревоги %v", oid, r))
			return r, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Тревоги %v", oid, value))
			if len(hs.StateHard.Status) > 2 {
				return getOctetString(hs.StateHard.Status[1]), nil
			}
			return getOctetString(0), nil
		})
	// case "1.3.6.1.4.1.1618.3.1.2.2.4.0":
	//	//Тревоги
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 1, 2, 2, 4, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			r := getOctetString(0)
			if len(hs.StateHard.Status) > 2 {
				r = getOctetString(hs.StateHard.Status[2])
			}
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Тревоги %v", oid, r))
			return r, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Тревоги %v", oid, value))
			if len(hs.StateHard.Status) > 2 {
				return getOctetString(hs.StateHard.Status[2]), nil
			}
			return getOctetString(0), nil
		})
	//	//Команды центра
	//
	// case "1.3.6.1.4.1.1618.3.2.2.4.1.0":
	//	//Установить КК
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 2, 2, 4, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Проверка КК", oid))
			if hs.StateHard.AllRed {
				return getUint32(2), nil
			}
			return getUint32(0), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.StateHard = hardware.GetStateHard()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить КК %v", oid, value))
			if !autonom {
				hs.Update()
				hardware.SetSourceCentralPlans()
				hardware.CommandToKDM(4, int(getUint32(value)))
				if int(getUint32(value)) == 0 {
					exitSpesial = true
					hardware.SetSourceLocal()
				}
			}
			return getUint32(value), nil
		})
	// case "1.3.6.1.4.1.1618.3.2.2.2.1.0":
	//	//Установить ОС
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 2, 2, 2, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Проверка ОС", oid))
			if hs.StateHard.Dark {
				return getUint32(2), nil
			}
			return getUint32(0), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить ОС %v", oid, value))
			if !autonom {
				hardware.SetSourceCentralPlans()
				hardware.CommandToKDM(6, int(getUint32(value)))
				if int(getUint32(value)) == 0 {
					exitSpesial = true
					hardware.SetSourceLocal()
				}
			}
			return getUint32(value), nil
		})
	// case "1.3.6.1.4.1.1618.3.2.2.1.1.0":
	//	//Установить ЖМ
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 2, 2, 1, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Проверка ЖМ", oid))
			if hs.StateHard.Flashing {
				return getUint32(2), nil
			}
			return getUint32(0), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить ЖМ %v", oid, value))
			if !autonom {
				hardware.SetSourceCentralPlans()
				hardware.CommandToKDM(3, int(getUint32(value)))
				if int(getUint32(value)) == 0 {
					exitSpesial = true
					hardware.SetSourceLocal()
				}
			}
			return getUint32(value), nil
		})
	// case "1.3.6.1.4.1.1618.3.2.2.1.2.0":
	//	//Установить ЖМ
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 2, 2, 1, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Проверка ЖМ", oid))
			if hs.StateHard.Flashing {
				return getUint32(2), nil
			}
			return getUint32(0), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить ЖМ %v", oid, value))
			if !autonom {
				hardware.SetSourceCentralPlans()
				hardware.CommandToKDM(3, int(getUint32(value)))
				if int(getUint32(value)) == 0 {
					exitSpesial = true
					hardware.SetSourceLocal()
				}
			}
			return getUint32(value), nil
		})
	// case "1.3.6.1.4.1.1618.3.7.2.11.1.0":
	//	//Установить Фазу
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 11, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Проверить Фазу %d", oid, hs.GetPhase()))
			return hs.GetPhase(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			if hs.StateHard.VPU {
				journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить Фазу %v отказ из-за ВПУ", oid, value))
				return hs.GetPhase(), nil
			}
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить Фазу %v", oid, value))
			if !autonom {
				hardware.SetSourceCentralPlans()
				phase := int(getUint32(value))
				if phase == 1 {
					phase = 8
				} else {
					phase--
				}
				if int(getUint32(value)) == 0 {
					hardware.CommandToKDM(1, 0)
					counter <- ctrlPhase{state: true, phase: 0}
				} else {
					hardware.CommandToKDM(8, phase)
					counter <- ctrlPhase{state: false, phase: phase}
				}
			}
			return hs.GetPhase(), nil
		})
	// case "1.3.6.1.4.1.1618.3.7.2.2.1.0":
	//	//Установить План
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 2, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Проверить План %d", oid, hs.GetPlan()))
			return hs.GetPlan(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить План %v", oid, value))
			if !autonom {
				hardware.SetSourceCentralPlans()
				hardware.CommandToKDM(7, int(getUint32(value)))
				if int(getUint32(value)) == 0 {
					hardware.CommandToKDM(1, 0)
				}
			}
			return hs.GetPlan(), nil
		})
	// PlanNumber = "1.3.6.1.4.1.1618.3.7.2.1.2.0"; // The number of the current traffic plan
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 1, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v The number of the current traffic plan %d", oid, hs.GetPlan()))
			return hs.GetPlan(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v The number of the current traffic plan %v", oid, value))
			return hs.GetPlan(), nil
		})
	// 1.3.6.1.4.1.1206.4.2.6.3.1.0
	//Установить время дату
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1206, 4, 2, 6, 3, 1, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Прочитать время дату %v", oid, getCount32(time.Now().Unix())))
			return getCount32(time.Now().Unix()), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			t := time.Unix(int64(getUint32(value)), 0)
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить время дату %v %s", oid, value, t.String()))
			err := SetSystemTimeFromUnix(int64(getUint32(value)))
			if err != nil {
				journal.SendMessage(journal.LevelSMNP, err.Error())
			}
			return value, nil
		})
	// 1.3.6.1.4.1.1618.3.7.2.2.3.0
	// Когда ставят план
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 2, 3, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v", oid))
			return getUint32(time.Now().Unix()), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v value %v", oid, value))
			return getUint32(value), nil
		})
	// 1.3.6.1.4.1.1618.3.7.2.2.2.0
	// когда ставят план
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 7, 2, 2, 2, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v", oid))
			return uint(1), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v value %v", oid, value))
			return getUint32(value), nil
		})
	// 1.3.6.1.4.1.1618.2.2.1.1.1
	// Запрос текущей таблицы подписок
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 2, 2, 1, 1, 1, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос текущей таблицы подписок", oid))
			return nil, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос текущей таблицы подписок %v", oid, value))
			return getUint32(value), nil
		})
	// 1.3.6.1.4.1.1618.3.5.2.1.1.0
	// Запрос времени 1
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 5, 2, 1, 1, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос времени", oid))
			return getUint32(time.Now().Unix()), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос времени %v", oid, value))
			return getUint32(time.Now().Unix()), nil
		})
	// 1.3.6.1.4.1.1618.3.6.2.1.1.0
	// Запрос времени 1
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 6, 2, 1, 1, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос времени", oid))
			return getUint32(time.Now().Unix()), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос времени %v", oid, value))
			return getUint32(time.Now().Unix()), nil
		})
	// 1.3.6.1.4.1.1618.3.6.2.2.1.0
	// Запрос времени 1
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 6, 2, 2, 1, 0},
		func(oid asn1.Oid) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос времени", oid))
			return getUint32(time.Now().Unix()), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос времени %v", oid, value))
			return getUint32(time.Now().Unix()), nil
		})
	// case "1.3.6.1.4.1.1618.3.6.2.1.3":
	//
	//	The reason for the major status of a traffic controller.
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 6, 2, 1, 3, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v The reason for the major status of a traffic controller", oid))
			return 0, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v The reason for the major status of a traffic controller %v", oid, value))
			return 0, nil
		})
	// case 1.3.6.1.4.1.1618.3.5.2.1.3
	//Здесь мне сначала самому нужно разобраться что это и что по нему Сварка отвечала. Мы его никогда не использовали...
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 6, 2, 1, 3, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v", oid))
			return 0, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v value %v", oid, value))
			return 0, nil
		})
	// stcipExtensionVersion GET 1.3.6.1.4.1.1618.3.11.1.0 i
	//   Указывает, какие версии расширения поддерживаются. Если объект не возвращает значение или значение = 0 — никакие расширения не поддерживается.
	//   1 расширение v1, 2 - расширение v2 и v2, 3 - расширение v1,v2,v3 и т.д.
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 11, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос версии расширения", oid))
			return 2, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Запрос версии расширения %v", oid, value))
			return 2, nil
		})
	//   stcipToLocalTimeout GET SET 1.3.6.1.4.1.1618.3.11.2.1.1.0 i
	//   Таймаут перехода в ЛКУ (локальный календарный режим) при потере связи с центром управления (0 -не переходить, секунды)
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 11, 2, 1, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Время контроля системы %d", oid, timeToLPU))
			return getUint32(timeToLPU), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить время контроля системы %v", oid, value))
			timeToLPU = int(getUint32(value))
			return value, nil
		})
	//   stcipControllerSwVersion GET 1.3.6.1.4.1.1618.3.11.2.3.1.0 s
	//   Версия ПО контроллера (например, "ELC2_v4.2.1")
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 11, 2, 3, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Версия контроллера %s", oid, hs.StateHard.GetStringTypeKDM()))
			return hs.StateHard.GetStringTypeKDM(), nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Версия контроллера установить %s", oid, hs.StateHard.GetStringTypeKDM()))
			return hs.StateHard.GetStringTypeKDM(), nil
		})

	//   stcipAdapterSwVersion GET 1.3.6.1.4.1.1618.3.11.2.3.2.0 s
	//   Версия ПО адаптера/прокси (например, "SwProxy_v1.8")
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 11, 2, 3, 2, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Версия адаптера %s", oid, "potop to "+setup.DateDeploy))
			return "potop to " + setup.DateDeploy, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			hs.Update()
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Версия адаптера установить %s", oid, "potop to "+setup.DateDeploy))
			return "potop to " + setup.DateDeploy, nil
		})
	//   stcipCabinetDoorOpen GET 1.3.6.1.4.1.1618.3.11.2.7.1.0 i
	//   Состояние двери шкафа (0 закрыто, 1 открыто)
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 11, 2, 7, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			i := 0
			if hs.StateHard.Key3 {
				i = 1
			}
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Дверь шкафа %d", oid, i))
			return i, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить дверь шкафа %v", oid, value))
			return value, nil
		})
	//  stcipAvailablePlansCount GET 1.3.6.1.4.1.1618.3.11.2.8.1.0 i
	//   Количество доступных транспортных программ (планов)
	agent.AddRwManagedObject(
		asn1.Oid{1, 3, 6, 1, 4, 1, 1618, 3, 11, 2, 8, 1, 0},
		func(oid asn1.Oid) (any, error) {
			hs.Update()
			config := hardware.GetConfig()
			var plans = 1
			if config.Ready {
				for _, v := range config.Plans.Plans {
					if v.Number > plans {
						plans = v.Number
					}
				}
			}
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Количество планов %d", oid, plans))
			return plans, nil
		},
		func(oid asn1.Oid, value any) (any, error) {
			journal.SendMessage(journal.LevelSMNP, fmt.Sprintf("%v Установить количество планов %v", oid, value))
			return value, nil
		})
}
