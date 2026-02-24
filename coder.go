package snmpc

func CodeError(code byte) uint {
	switch code {
	case 0:
		return 0
	case 1:
		return 44
	case 2:
		return 42
	case 3:
		return 32
	case 4:
		return 35
	case 5:
		return 30
	case 6:
		return 72
	case 7:
		return 73
	case 8:
		return 7
	case 9:
		return 5
	case 10:
		return 25
	case 11:
		return 24
	case 12:
		return 70
	}
	return uint(code)
}
