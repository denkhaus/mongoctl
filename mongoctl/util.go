package mongoctl

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func LogCombinedLines(prefix string, log *logrus.Entry, data []byte) {
	msgs := strings.Split(string(data), "\n")

	for idx, msg := range msgs {
		log.Infof("%s-%d->%s", prefix, idx, msg)
	}
}
