package ipmi

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestVpdParse(t *testing.T) {
	assert := assert.New(t)

	var fan FanContraller

	fan = NewIntelFan()
	assert.NotNil(fan)
	logrus.Info(fan)

	fan = NewLenovoFan()
	assert.NotNil(fan)
	logrus.Info(fan)

	fan = NewInspurFan()
	assert.NotNil(fan)
	logrus.Info(fan)

	fan = NewInspurFan2()
	assert.NotNil(fan)
	logrus.Info(fan)

	fan = NewSugonFan()
	assert.NotNil(fan)
	logrus.Info(fan)

	logrus.Info(data)
}
