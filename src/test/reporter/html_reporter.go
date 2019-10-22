package reporter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

type HtmlReporter struct {
	filename        string
	suite           HtmlTestSuite
	generateSummary bool
	reportUrl       string
}

type HtmlTestSuite struct {
	FailedTestCases []HtmlTestCase
	OtherTestCases  []HtmlTestCase
	Name            string
	TotalNum        int
	FailedNum       int
	SuccessNum      int
	State           string
	RunTime         float64
	Url             string
}

type HtmlTestCase struct {
	Name    string
	State   string
	Detail  string
	RunTime float64
}

var dir string

func NewHtmlReporter(filename string, reportUrl string, generateSummary bool) *HtmlReporter {
	return &HtmlReporter{
		filename:        strings.Trim(filename, "/"),
		generateSummary: generateSummary,
		reportUrl:       reportUrl,
	}
}

func (reporter *HtmlReporter) SpecSuiteWillBegin(ginkgoConfig config.GinkgoConfigType, summary *types.SuiteSummary) {
	reporter.suite = HtmlTestSuite{
		Name:            summary.SuiteDescription,
		FailedTestCases: []HtmlTestCase{},
		OtherTestCases:  []HtmlTestCase{},
	}
	reporter.suite.Url = reporter.reportUrl + reporter.filename
	dir = reporter.filename[:strings.LastIndex(reporter.filename, "/")]
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			reporter.filename = reporter.filename[strings.LastIndex(reporter.filename, "/")+1:]
		}
	}
}

func (reporter *HtmlReporter) SpecWillRun(specSummary *types.SpecSummary) {
}

func (reporter *HtmlReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {
	reporter.handleSetupSummary("BeforeSuite", setupSummary)
}

func (reporter *HtmlReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {
	reporter.handleSetupSummary("AfterSuite", setupSummary)
}

func (reporter *HtmlReporter) handleSetupSummary(name string, setupSummary *types.SetupSummary) {
	testCase := HtmlTestCase{
		Name:    name,
		RunTime: setupSummary.RunTime.Seconds(),
		State:   reporter.transformState(setupSummary.State),
	}
	if setupSummary.State == types.SpecStateFailed {
		testCase.Detail = fmt.Sprintf("%s<br><br>%s", setupSummary.Failure.Location.String(), setupSummary.Failure.Message)
		reporter.suite.FailedTestCases = append(reporter.suite.FailedTestCases, testCase)
	} else {
		testCase.Detail = setupSummary.CapturedOutput
		reporter.suite.OtherTestCases = append(reporter.suite.OtherTestCases, testCase)
	}
}

func (reporter *HtmlReporter) transformState(state types.SpecState) string {
	switch state {
	case types.SpecStateInvalid:
		return "Invalid"
	case types.SpecStatePending:
		return "Pending"
	case types.SpecStateSkipped:
		return "Skipped"
	case types.SpecStatePassed:
		return "Passed"
	case types.SpecStateFailed:
		return "Failed"
	case types.SpecStatePanicked:
		return "Panicked"
	case types.SpecStateTimedOut:
		return "TimedOut"
	default:
		return ""
	}
}

func (reporter *HtmlReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	testCase := HtmlTestCase{
		Name:    strings.Join(specSummary.ComponentTexts[1:], "/"),
		RunTime: specSummary.RunTime.Seconds(),
		State:   reporter.transformState(specSummary.State),
	}
	if specSummary.State == types.SpecStateFailed {
		testCase.Detail = fmt.Sprintf("%s<br><br>%s", specSummary.Failure.Location.String(), specSummary.Failure.Message)
		reporter.suite.FailedTestCases = append(reporter.suite.FailedTestCases, testCase)
	} else {
		testCase.Detail = specSummary.CapturedOutput
		reporter.suite.OtherTestCases = append(reporter.suite.OtherTestCases, testCase)
	}
}

func (reporter *HtmlReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	reporter.suite.TotalNum = summary.NumberOfSpecsThatWillBeRun
	reporter.suite.RunTime = summary.RunTime.Seconds()
	reporter.suite.FailedNum = summary.NumberOfFailedSpecs
	reporter.suite.SuccessNum = reporter.suite.TotalNum - reporter.suite.FailedNum
	if summary.SuiteSucceeded {
		reporter.suite.State = "Passed"
	} else {
		reporter.suite.State = "Failed"
	}
	file, err := os.Create(reporter.filename)
	if err != nil {
		fmt.Printf("Failed to create Html report file: %s\n\t%s", reporter.filename, err.Error())
		return
	}
	defer file.Close()
	t, err := template.New("HtmlTemplate").Parse(HtmlTemplate)
	if nil != err {
		fmt.Printf("Failed to parse Html template: %s\n\t%s", reporter.filename, err.Error())
		return
	}
	err = t.Execute(file, reporter.suite)
	if err != nil {
		fmt.Printf("Failed to generate Html report file: %s\n\t%s", reporter.filename, err.Error())
		return
	}

	if reporter.generateSummary {
		reporter.generateSummaryHtml()
	}
}

func (reporter *HtmlReporter) generateSummaryHtml() {
	var buf bytes.Buffer
	if _, err := os.Stat(dir + "/summary.html"); err != nil && os.IsNotExist(err) {
		summary, err := os.Create(dir + "/summary.html")
		if err != nil {
			fmt.Printf("Failed to create summary Html report file\n\t%s", err.Error())
			return
		}
		defer summary.Close()
		temp, err := template.New("SummaryHtmlTemplate").Parse(SummaryHtmlTemplate)
		if nil != err {
			fmt.Printf("Failed to parse summary Html template\n\t%s", err.Error())
			return
		}
		err = temp.Execute(summary, reporter)
		if err != nil {
			fmt.Printf("Failed to generate summary Html report file\n\t%s", err.Error())
			return
		}
	}

	summary, err := ioutil.ReadFile(dir + "/summary.html")
	if err != nil {
		fmt.Printf("Failed to open summary Html report file\n\t%s", err.Error())
		return
	}

	sum := string(summary)
	index := strings.Index(sum, "</table>")
	temp, err := template.New("SummaryTemplate").Parse(SummaryTemplate)
	if nil != err {
		fmt.Printf("Failed to parse summary Html template\n\t%s", err.Error())
		return
	}
	err = temp.Execute(&buf, reporter.suite)
	if err != nil {
		fmt.Printf("Failed to generate summary Html report file\n\t%s", err.Error())
		return
	}
	sum = sum[:index] + buf.String() + sum[index:]
	buf.Reset()

	index = strings.Index(sum, "</body>")
	temp, err = template.New("FailedTemplate").Parse(FailedTemplate)
	if nil != err {
		fmt.Printf("Failed to parse summary Html template\n\t%s", err.Error())
		return
	}
	err = temp.Execute(&buf, reporter.suite)
	if err != nil {
		fmt.Printf("Failed to generate summary Html report file\n\t%s", err.Error())
		return
	}
	sum = sum[:index] + buf.String() + sum[index:]

	err = ioutil.WriteFile(dir+"/summary.html", []byte(sum), 0666)
	if err != nil {
		fmt.Printf("Failed to write summary Html report file\n\t%s", err.Error())
		return
	}
}
