package ipmi

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

type BasefanSetter interface {
	// get fan speed in %
	GetFanSpeed() (string, error)

	getFanSpeedMin() int
	setFanSpeedPre() (string, error)
	setFanSpeedPost() error
}

type FanSetter interface {
	BasefanSetter
	setFanSpeedPercent(int) (string, error)
}

type loopfanSetter interface {
	BasefanSetter
	getFanIDs() []string
	setFanSpeedOne(string, int) (string, error)
	getFanSpeedOne(string) (string, error)
}

type FanContraller interface {
	FanSetter
	// get fan speed in RPM
	GetFanRPM() (string, error)
	// set fan speed in %
	SetFanSpeed(int) error
}

type fan struct{}

func (f *fan) getFanSpeedMin() int                        { return 20 }
func (f *fan) setFanSpeedPre() (string, error)            { return "", nil }
func (f *fan) setFanSpeedPost() error                     { return nil }
func (f *fan) setFanSpeedPercent(int) (string, error)     { return "", fmt.Errorf("not support") }
func (f *fan) setFanSpeedOne(string, int) (string, error) { return "", fmt.Errorf("not support") }
func (f *fan) getFanSpeedOne(string) (string, error)      { return "", fmt.Errorf("not support") }
func (f *fan) getFanIDs() ([]string, error)               { return nil, fmt.Errorf("not support") }

type loopFan struct {
	loopfanSetter
}

func (f *loopFan) setFanSpeedPercent(speedPercent int) (string, error) {
	for _, id := range f.getFanIDs() {
		info, err := f.setFanSpeedOne(id, speedPercent)
		if err != nil {
			return "", err
		}
		logrus.Debug(info)
	}
	return "", nil
}
func (f *loopFan) getFanSpeedPercent() (string, error) {
	ret := ""
	for _, id := range f.getFanIDs() {
		info, err := f.getFanSpeedOne(id)
		if err != nil {
			return "", err
		}
		ret += info
	}
	return ret, nil
}

type fanContraller struct{ FanSetter }

func (f *fanContraller) GetFanRPM() (string, error) {
	return ipmitool("sdr type Fan")
}
func (f *fanContraller) setFanSpeedCheck(speedPercent int) error {
	min := f.getFanSpeedMin()
	if speedPercent < min || speedPercent > 100 {
		return fmt.Errorf("speed must between %d~100", min)
	}
	return nil
}
func (f *fanContraller) SetFanSpeed(speedPercent int) error {
	if err := f.setFanSpeedCheck(speedPercent); err != nil {
		return err
	}
	if info, err := f.setFanSpeedPre(); err != nil {
		return err
	} else {
		logrus.Debug(info)
	}
	if info, err := f.setFanSpeedPercent(speedPercent); err != nil {
		return err
	} else {
		logrus.Debug(info)
	}
	if err := f.setFanSpeedPost(); err != nil {
		return err
	}
	return nil
}

// intel fan
type intelFan struct{ fan }

func NewIntelFan() FanContraller                 { return &fanContraller{FanSetter: &loopFan{&intelFan{}}} }
func (f *intelFan) GetFanSpeed() (string, error) { return "", fmt.Errorf("intel bmc not support") }
func (f *intelFan) setFanSpeedPre() (string, error) {
	// set factory mode
	return ipmitool("raw 0x06 0x05 0x73 0x28 0x58 0x7A 0x4E 0x57 0x50 0x4F 0x3A 0x6F 0x65 0x2F 0x60 0x57")
}

func (f *intelFan) getFanIDs() []string {
	return []string{"00", "01", "02", "03", "04", "05"}
}
func (f *intelFan) setFanSpeedOne(fanID string, speedPercent int) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x30 0x15 0x05 %s 01 0x%x", fanID, speedPercent))
}

// inspur
//func (f *inspurFan) getFanCtrlMode() (string, error) {
//	return ipmitool("raw 0x3a 0x7b")
//}
type inspurFan struct{ fan }

func NewInspurFan() FanContraller { return &fanContraller{FanSetter: &loopFan{&inspurFan{}}} }
func (f *inspurFan) getFanIDs() []string {
	return []string{"0", "2", "4", "6"}
}
func (f *inspurFan) GetFanSpeed() (string, error) { return "", fmt.Errorf("inspur bmc not support") }
func (f *inspurFan) setFanSpeedOne(fanID string, speedPercent int) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x3a 0x78 %s %d", fanID, speedPercent))
}
func (f *inspurFan) getFanSpeedOne(fanID string) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x3a 0x79 %s", fanID))
}
func (f *inspurFan) setFanSpeedPre() (string, error) {
	// 0 auto mode ; 1 manual mode ;
	return ipmitool(fmt.Sprintf("raw 0x3a 0x7a %d", 1))
}

// inspur fan 2
type inspurFan2 struct{ fan }

func NewInspurFan2() FanContraller { return &fanContraller{FanSetter: &inspurFan2{}} }
func (f *inspurFan2) setFanSpeedPre() (string, error) {
	return f.setFanCtrlMode()
}
func (f *inspurFan2) setFanSpeedPercent(speedPercent int) (string, error) {
	// 0xFF : all fans
	return ipmitool(fmt.Sprintf("raw 0x3c 0x2d %s 0x%x", "0xff", speedPercent))
}
func (f *inspurFan2) GetFanSpeed() (string, error) {
	// 0xFF : all fans
	return ipmitool(fmt.Sprintf("raw 0x3c 0x2e %s", "0xff"))
}
func (f *inspurFan2) setFanCtrlMode() (string, error) {
	// 0x00 auto mode ; 0x01 manual mode ;
	return ipmitool("raw 0x3c 0x2f 0x01")
}
func (f *inspurFan2) getFanCtrlMode() (string, error) {
	// respons: 0x00 success
	// other failed
	return ipmitool("raw 0x3c 0x30")
}

// lenovo
type lenovoFan struct{ fan }

func NewLenovoFan() FanContraller                 { return &fanContraller{FanSetter: &lenovoFan{}} }
func (f *lenovoFan) getFanSpeedMin() int          { return 30 }
func (f *lenovoFan) GetFanSpeed() (string, error) { return "", fmt.Errorf("lenovo bmc not support") }
func (f *lenovoFan) setFanSpeedPercent(speedPercent int) (string, error) {
	xs := speedPercent * 255 / 100
	ids := "0xFF" // all fans
	cmd := fmt.Sprintf("raw 0x3a 0x07 %s 0x%X 0x01 0x00 0x00 0x00 0x00 0x00", ids, xs)
	return ipmitool(cmd)
}

// sugon
/* TODO:
 *  3.4 退出风扇调速手动模式，进入自动控制模式
 * Balance Mode:
 * ipmitool -H x.x.x.x –U  admin  -P  admin  raw  0x3a   0xb   0x0  0x1
 *
 * Silence Mode:
 * ipmitool -H x.x.x.x –U  admin  -P  admin  raw  0x3a   0xb   0x0  0x2
 *
 * Performance Mode:
 * ipmitool -H x.x.x.x –U  admin  -P  admin  raw  0x3a   0xb   0x0  0x3
 */
type sugonFan struct{ fan }

func NewSugonFan() FanContraller                 { return &fanContraller{FanSetter: &sugonFan{}} }
func (f *sugonFan) GetFanSpeed() (string, error) { return "", fmt.Errorf("sugon bmc not support") }
func (f *sugonFan) setFanSpeedPercent(speedPercent int) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x3a 0xd 0xFF %d", speedPercent))
}
func (f *sugonFan) setFanSpeedPre() (string, error) {
	return ipmitool("raw 0x3a 0xb 0x0 0x0")
}
