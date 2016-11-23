package mode

import (
	"fmt"
	"time"
)

type Grains struct {
	Input      string
	AzureStyle string
	Duration   time.Duration
}

var supported = map[string]Grains{
	"5m": {Input: "5m", AzureStyle: "PT5M", Duration: 5 * time.Minute},
	"1h": {Input: "1h", AzureStyle: "PT1H", Duration: 1 * time.Hour},
}

func GetSupportedGrains() []string {
	keys := make([]string, 0, len(supported))
	for s := range supported {
		keys = append(keys, s)
	}
	return keys
}

func (g Grains) Start() string {
	return time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
}

func (g Grains) End() string {
	return time.Now().Format(time.RFC3339)
}

func NewGrain(input string) (*Grains, error) {
	if val, ok := supported[input]; ok {
		return &val, nil
	}
	return nil, fmt.Errorf("This timeunit(%s) is not supported! Supported are: %s", input, GetSupportedGrains())
}

func ParseTimestamp(t string) (string, error) {
	t1, err := time.Parse(time.RFC3339, t)
	return t1.Local().String(), err
}
