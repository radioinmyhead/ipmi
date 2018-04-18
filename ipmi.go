package ipmi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/radioinmyhead/shell"
)

func ipmitool(cmd string) (string, error) {
	return shell.Out(fmt.Sprintf("ipmitool %s", cmd))
}

type IPMI struct {
	faner
}

func GetLocalIPMI() (ret *IPMI, err error) {
	ret = &IPMI{}
	str, err := shell.Out("dmidecode  -t 38")
	if err != nil {
		return
	}
	if !strings.Contains(str, "KCS") {
		err = fmt.Errorf("no support; info=%s", str)
		return
	}
	str, err = ipmitool("mc info")
	if err != nil {
		return
	}
	mid, pid, err := parseIPMIMcInfo(str)
	if err != nil {
		return
	}
	fan := getFan(mid, pid)
	if fan == nil {
		err = fmt.Errorf("not support! info=%s", str)
	}
	ret.faner = fan
	return
}

func parseIPMIMcInfo(info string) (mid, pid int, err error) {
	for _, line := range strings.Split(info, "\n") {
		if strings.Contains(line, "Manufacturer ID") {
			mid, err = getNum(line)
			if err != nil {
				return
			}
		}
		if strings.Contains(line, "Product ID") {
			pid, err = getNum(line)
			if err != nil {
				return
			}
		}

	}
	return
}
func getNum(line string) (int, error) {
	tmp := strings.Split(line, ":")
	if len(tmp) != 2 {
		return 0, fmt.Errorf("failed; line=%s", line)
	}
	tmp = strings.Split(tmp[1], "(")
	return strconv.Atoi(strings.Trim(tmp[0], " "))
}
