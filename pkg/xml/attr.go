package xml

import (
	"strings"
)

func ParseAttrMap(attr string) map[string]string {
	attrMap := make(map[string]string)
	for _, keyValPair := range strings.Split(attr, ",") {
		splitted := strings.Split(strings.TrimSpace(keyValPair), ":")
		attrMap[strings.TrimSpace(splitted[0])] = strings.TrimSpace(splitted[1])
	}

	return attrMap
}
