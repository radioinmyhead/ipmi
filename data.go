package ipmi

type Product struct {
	name string
	fan  faner
}

type Manufactory struct {
	name string
	list map[int]Product
}

var data map[int]Manufactory

func addFan(mid, pid int, mname, pname string, fan faner) {
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

func getFan(mid, pid int) faner {
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

	addFan(37945, 43707, "inspur", "unknown", &inspurFan{})
	addFan(343, 111, "Intel", "unknown", &intelFan{})
}
