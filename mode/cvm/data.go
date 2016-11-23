package cvm

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/griesbacher/check_azure/azureHttp"
	"github.com/griesbacher/check_azure/mode"
	"github.com/griesbacher/check_x"
	"github.com/griesbacher/check_azure/helper"
)

//chidley.exe -G  -e "" -t vm-cpu-5m.xml

func Cpu(ac azureHttp.AzureConnector, res, name, givenGrain string, warn, crit *check_x.Threshold) error {
	err, state, text := requestSingleData(ac, res, name, givenGrain, warn, crit, "Percentage CPU", "%")
	check_x.ExitOnError(err)
	check_x.Exit(*state, text)
	return nil
}
func Network(ac azureHttp.AzureConnector, res, name, givenGrain, warn, crit  string) error {
	thresholds, err := helper.ParseCommaThresholds(warn, crit)
	check_x.ExitOnError(err)
	states := check_x.States{}
	msg := ""
	for i, typ := range ([]string{"Network In", "Network Out"}) {
		var w *check_x.Threshold
		var c *check_x.Threshold
		if len((*thresholds)["warning"]) - 1 >= i {
			w = (*thresholds)["warning"][i]
		}
		if len((*thresholds)["critical"]) - 1 >= i {
			c = (*thresholds)["critical"][i]
		}
		err, state, text := requestSingleData(ac, res, name, givenGrain, w, c, typ, "B")
		check_x.ExitOnError(err)
		states = append(states, *state)
		msg += text + "\n"
	}
	worst, err := states.GetWorst()
	check_x.ExitOnError(err)
	check_x.LongExit(*worst, "Network", msg)
	return nil
}

func Disk(ac azureHttp.AzureConnector, res, name, givenGrain, warn, crit  string) error {
	thresholds, err := helper.ParseCommaThresholds(warn, crit)
	check_x.ExitOnError(err)
	states := check_x.States{}
	msg := ""
	for i, typ := range ([]string{"Disk Read Bytes/sec", "Disk Write Bytes/sec"}) {
		var w *check_x.Threshold
		var c *check_x.Threshold
		if len((*thresholds)["warning"]) - 1 >= i {
			w = (*thresholds)["warning"][i]
		}
		if len((*thresholds)["critical"]) - 1 >= i {
			c = (*thresholds)["critical"][i]
		}
		err, state, text := requestSingleData(ac, res, name, givenGrain, w, c, typ, "Bps")
		check_x.ExitOnError(err)
		states = append(states, *state)
		msg += text + "\n"
	}
	worst, err := states.GetWorst()
	check_x.ExitOnError(err)
	check_x.LongExit(*worst, "Disk", msg)
	return nil
}

func requestSingleData(ac azureHttp.AzureConnector, res, name, givenGrain string, warn, crit *check_x.Threshold, typ, unit string) (error, *check_x.State, string) {
	grain, err := mode.NewGrain(givenGrain)
	check_x.ExitOnError(err)

	urlPath := fmt.Sprintf(
		"resourceGroups/%s/providers/Microsoft.ClassicCompute/virtualMachines/%s/metrics",
		res, name,
	)
	filter := fmt.Sprintf(
		"name.value eq '%s' and startTime eq %s and endTime eq %s and timeGrain eq duration'%s'",
		typ, grain.Start(), grain.End(), grain.AzureStyle,
	)
	body, typ, err := ac.RequestWithSub("2014-04-01", urlPath, filter)
	check_x.ExitOnError(err)

	if typ != azureHttp.ContentTypeXML {
		return azureHttp.ContentError(azureHttp.ContentTypeXML, typ), nil, ""
	}

	var item MetricValueSetCollection
	decoder := xml.NewDecoder(bytes.NewReader(body))
	token, err := decoder.Token()
	check_x.ExitOnError(err)

	switch se := token.(type) {
	case xml.StartElement:
		err := decoder.DecodeElement(&item, &se)
		if err != nil {
			return err, nil, ""
		}
	default:
		return errors.New("Could not decode XML Data"), nil, ""
	}
	lastOne := item.Value.MetricValueSet.MetricValues.MetricValue[len(item.Value.MetricValueSet.MetricValues.MetricValue) - 1]
	t, err := mode.ParseTimestamp(lastOne.Timestamp.Text)
	check_x.ExitOnError(err)

	perf := check_x.NewPerformanceData(item.Value.MetricValueSet.Name.Value.Text, lastOne.Average.Text).Warn(warn).Crit(crit)
	if unit != "" {
		perf.Unit(unit)
	}
	result := check_x.Evaluator{Warning: warn, Critical: crit, }.Evaluate(lastOne.Average.Text)
	return nil,
		&result,
		fmt.Sprintf("%s last checked: %s", item.Value.MetricValueSet.Name.Value.Text, t)
}

type Root struct {
	MetricValueSetCollection *MetricValueSetCollection `xml:"http://schemas.microsoft.com/windowsazure MetricValueSetCollection,omitempty" json:"MetricValueSetCollection,omitempty"`
}

type MetricValue struct {
	Average    *Average    `xml:"http://schemas.microsoft.com/windowsazure average,omitempty" json:"average,omitempty"`
	Count      *Count      `xml:"http://schemas.microsoft.com/windowsazure count,omitempty" json:"count,omitempty"`
	Maximum    *Maximum    `xml:"http://schemas.microsoft.com/windowsazure maximum,omitempty" json:"maximum,omitempty"`
	Minimum    *Minimum    `xml:"http://schemas.microsoft.com/windowsazure minimum,omitempty" json:"minimum,omitempty"`
	Properties *Properties `xml:"http://schemas.microsoft.com/windowsazure properties,omitempty" json:"properties,omitempty"`
	Timestamp  *Timestamp  `xml:"http://schemas.microsoft.com/windowsazure timestamp,omitempty" json:"timestamp,omitempty"`
	Total      *Total      `xml:"http://schemas.microsoft.com/windowsazure total,omitempty" json:"total,omitempty"`
	XMLName    xml.Name    `xml:"http://schemas.microsoft.com/windowsazure MetricValue,omitempty" json:"MetricValue,omitempty"`
}

type MetricValueSet struct {
	EndTime      *EndTime      `xml:"http://schemas.microsoft.com/windowsazure endTime,omitempty" json:"endTime,omitempty"`
	Id           *Id           `xml:"http://schemas.microsoft.com/windowsazure id,omitempty" json:"id,omitempty"`
	MetricValues *MetricValues `xml:"http://schemas.microsoft.com/windowsazure metricValues,omitempty" json:"metricValues,omitempty"`
	Name         *Name         `xml:"http://schemas.microsoft.com/windowsazure name,omitempty" json:"name,omitempty"`
	Properties   *Properties   `xml:"http://schemas.microsoft.com/windowsazure properties,omitempty" json:"properties,omitempty"`
	ResourceId   *ResourceId   `xml:"http://schemas.microsoft.com/windowsazure resourceId,omitempty" json:"resourceId,omitempty"`
	StartTime    *StartTime    `xml:"http://schemas.microsoft.com/windowsazure startTime,omitempty" json:"startTime,omitempty"`
	TimeGrain    *TimeGrain    `xml:"http://schemas.microsoft.com/windowsazure timeGrain,omitempty" json:"timeGrain,omitempty"`
	Unit         *Unit         `xml:"http://schemas.microsoft.com/windowsazure unit,omitempty" json:"unit,omitempty"`
	XMLName      xml.Name      `xml:"http://schemas.microsoft.com/windowsazure MetricValueSet,omitempty" json:"MetricValueSet,omitempty"`
}

type MetricValueSetCollection struct {
	Attr_i     string   `xml:"xmlns i,attr"  json:",omitempty"`
	Attr_xmlns string   `xml:" xmlns,attr"  json:",omitempty"`
	Id         *Id      `xml:"http://schemas.microsoft.com/windowsazure id,omitempty" json:"id,omitempty"`
	Value      *Value   `xml:"http://schemas.microsoft.com/windowsazure value,omitempty" json:"value,omitempty"`
	XMLName    xml.Name `xml:"http://schemas.microsoft.com/windowsazure MetricValueSetCollection,omitempty" json:"MetricValueSetCollection,omitempty"`
}

type Average struct {
	Text    float64  `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure average,omitempty" json:"average,omitempty"`
}

type Count struct {
	Text    int8     `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure count,omitempty" json:"count,omitempty"`
}

type EndTime struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure endTime,omitempty" json:"endTime,omitempty"`
}

type Id struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure id,omitempty" json:"id,omitempty"`
}

type LocalizedValue struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure localizedValue,omitempty" json:"localizedValue,omitempty"`
}

type Maximum struct {
	Text    float64  `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure maximum,omitempty" json:"maximum,omitempty"`
}

type MetricValues struct {
	MetricValue []*MetricValue `xml:"http://schemas.microsoft.com/windowsazure MetricValue,omitempty" json:"MetricValue,omitempty"`
	XMLName     xml.Name       `xml:"http://schemas.microsoft.com/windowsazure metricValues,omitempty" json:"metricValues,omitempty"`
}

type Minimum struct {
	Text    float64  `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure minimum,omitempty" json:"minimum,omitempty"`
}

type Name struct {
	LocalizedValue *LocalizedValue `xml:"http://schemas.microsoft.com/windowsazure localizedValue,omitempty" json:"localizedValue,omitempty"`
	Value          *Value          `xml:"http://schemas.microsoft.com/windowsazure value,omitempty" json:"value,omitempty"`
	XMLName        xml.Name        `xml:"http://schemas.microsoft.com/windowsazure name,omitempty" json:"name,omitempty"`
}

type Properties struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure properties,omitempty" json:"properties,omitempty"`
}

type ResourceId struct {
	Attr_i_nil string   `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"  json:",omitempty"`
	XMLName    xml.Name `xml:"http://schemas.microsoft.com/windowsazure resourceId,omitempty" json:"resourceId,omitempty"`
}

type StartTime struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure startTime,omitempty" json:"startTime,omitempty"`
}

type TimeGrain struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure timeGrain,omitempty" json:"timeGrain,omitempty"`
}

type Timestamp struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure timestamp,omitempty" json:"timestamp,omitempty"`
}

type Total struct {
	Text    float64  `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure total,omitempty" json:"total,omitempty"`
}

type Unit struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure unit,omitempty" json:"unit,omitempty"`
}

type Value struct {
	MetricValueSet *MetricValueSet `xml:"http://schemas.microsoft.com/windowsazure MetricValueSet,omitempty" json:"MetricValueSet,omitempty"`
	Text           string          `xml:",chardata" json:",omitempty"`
	XMLName        xml.Name        `xml:"http://schemas.microsoft.com/windowsazure value,omitempty" json:"value,omitempty"`
}
