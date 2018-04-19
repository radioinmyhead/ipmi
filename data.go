package ipmi

type Product struct {
	name string
	fan  FanContraller
}

type Manufactory struct {
	name string
	list map[int]Product
}

var data map[int]Manufactory

func addFan(mid, pid int, mname, pname string, fan FanContraller) {
	if _, ok := data[mid]; !ok {
		data[mid] = Manufactory{
			name: mname,
			list: make(map[int]Product),
		}
	}
	if _, ok := data[mid].list[pid]; !ok {
		data[mid].list[pid] = Product{
			name: pname,
			fan:  fan,
		}
	}
}

func getFan(mid, pid int) FanContraller {
	if _, ok := data[mid]; !ok {
		return nil
	}
	if _, ok := data[mid].list[pid]; !ok {
		return nil
	}
	return data[mid].list[pid].fan
}

func init() {
	data = make(map[int]Manufactory)

	addFan(343, 111, "Intel Corporation", "S2600WT2R", NewIntelFan())
	addFan(19046, 1087, "Lenovo", "ThinkSystem SR650 -[7X06CTO1WW]-", NewLenovoFan())
	//addFan(19046, 323, "Lenovo", "unknown", NewLenovoFan())
	//addFan(37945, 43707, "inspur", "unknown", NewInspurFan())
	addFan(37945, 514, "Inspur", "SA5212M5", NewInspurFan2())
	//addFan(27500, 0, "Sugon", "unknown", NewSugonFan())
}
