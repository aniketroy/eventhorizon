package writer

import (
	"errors"
	"github.com/function61/eventhorizon/metaevents"
	"path"
	"strings"
)

func parentStreamName(streamName string) string {
	return path.Dir(streamName)
}

func stringArrayToRawLines(contentArr []string) (string, error) {
	buf := ""

	for _, line := range contentArr {
		if strings.Contains(line, "\n") {
			return "", errors.New("content cannot contain \\n")
		}

		buf += metaevents.EncodeRegularLine(line) + "\n"
	}

	return buf, nil
}
