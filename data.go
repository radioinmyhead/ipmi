package ipmi

type product struct {
	name string
	fan  FanContraller
}

type manufactory struct {
	name string
	list map[int]product
}

var data map[int]manufactory

func RegisterFan(mid, pid int, mname, pname string, fan FanContraller) {
	if _, ok := data[mid]; !ok {
		data[mid] = manufactory{
			name: mname,
			list: make(map[int]product),
		}
	}
	if _, ok := data[mid].list[pid]; !ok {
		data[mid].list[pid] = product{
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
	data = make(map[int]manufactory)
}
