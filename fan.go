package ipmi

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

type faner interface {
	GetFanSpeed() error
	SetFanSpeed(int) error
	CheckSpeed(int) error
}

type Fan struct{}

func (f *Fan) CheckSpeed(speedPercent int) error {
	if speedPercent < 20 || speedPercent > 100 {
		return fmt.Errorf("speed must between 20~100")
	}
	return nil
}

// inspur
type inspurFan struct{ Fan }

func (f *inspurFan) GetFanSpeed() error {
	for _, id := range f.getFanIDs() {
		info, err := f.getFanSpeed(id)
		if err != nil {
			return err
		}
		logrus.Info(info)
	}
	return nil
}
func (f *inspurFan) SetFanSpeed(speedPercent int) error {
	if _, err := f.setFanCtrlMode(1); err != nil {
		return err
	}
	for _, id := range f.getFanIDs() {
		info, err := f.setFanSpeed(id, speedPercent)
		if err != nil {
			return err
		}
		logrus.Info(info)
	}
	return nil
}
func (f *inspurFan) getFanIDs() []string {
	return []string{"0", "2", "4", "6"}
}
func (f *inspurFan) setFanSpeed(fanID string, speedPercent int) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x3a 0x78 %s %d", fanID, speedPercent))
}
func (f *inspurFan) getFanSpeed(fanID string) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x3a 0x79 %s", fanID))
}
func (f *inspurFan) setFanCtrlMode(ctrl int) (string, error) {
	// 0 auto mode ; 1 manual mode ;
	return ipmitool(fmt.Sprintf("raw 0x3a 0x7a %d", ctrl))
}
func (f *inspurFan) getFanCtrlMode() (string, error) {
	return ipmitool("raw 0x3a 0x7b")
}

// intel fan
type intelFan struct{ Fan }

func (f *intelFan) GetFanSpeed() error {
	return fmt.Errorf("intel bmc not support")
}
func (f *intelFan) SetFanSpeed(speedPercent int) error {
	if _, err := f.setFactory(); err != nil {
		return err
	}
	for _, id := range f.getFanIDs() {
		info, err := f.setFanSpeed(id, speedPercent)
		if err != nil {
			return err
		}
		logrus.Info(info)
	}
	return nil
}
func (f *intelFan) getFanIDs() []string {
	return []string{"00", "01", "02", "03", "04", "05"}
}
func (f *intelFan) setFactory() (string, error) {
	return ipmitool("raw 0x06 0x05 0x73 0x28 0x58 0x7A 0x4E 0x57 0x50 0x4F 0x3A 0x6F 0x65 0x2F 0x60 0x57")
}
func (f *intelFan) setFanSpeed(fanID string, speedPercent int) (string, error) {
	return ipmitool(fmt.Sprintf("raw 0x30 0x15 0x05 %s 01 0x%x", fanID, speedPercent))
}
