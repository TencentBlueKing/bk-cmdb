package reporter

import (
	"fmt"
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
	TestCases  []HtmlTestCase
	Name       string
	TotalNum   int
	FailedNum  int
	SuccessNum int
	State      string
	RunTime    float64
	Url        string
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
		filename:        filename,
		generateSummary: generateSummary,
		reportUrl:       reportUrl,
	}
}

func (reporter *HtmlReporter) SpecSuiteWillBegin(ginkgoConfig config.GinkgoConfigType, summary *types.SuiteSummary) {
	reporter.suite = HtmlTestSuite{
		Name:      summary.SuiteDescription,
		TestCases: []HtmlTestCase{},
	}
	reporter.suite.Url = reporter.reportUrl + reporter.filename
	dir = reporter.filename[:strings.LastIndex(reporter.filename, "/")]
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
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
	var detail string
	if setupSummary.State == types.SpecStatePassed {
		detail = setupSummary.CapturedOutput
	} else {
		detail = fmt.Sprintf("%s<br><br>%s", setupSummary.Failure.Location.String(), setupSummary.Failure.Message)
	}
	testCase.Detail = detail
	reporter.suite.TestCases = append(reporter.suite.TestCases, testCase)
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
	var detail string
	if specSummary.State == types.SpecStatePassed {
		detail = specSummary.CapturedOutput
	} else {
		detail = fmt.Sprintf("%s<br><br>%s", specSummary.Failure.Location.String(), specSummary.Failure.Message)
	}
	testCase.Detail = detail
	reporter.suite.TestCases = append(reporter.suite.TestCases, testCase)
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
	}
	defer file.Close()
	t, err := template.New("HtmlTemplate").Parse(HtmlTemplate)
	if nil != err {
		fmt.Printf("Failed to parse Html template: %s\n\t%s", reporter.filename, err.Error())
	}
	err = t.Execute(file, reporter.suite)
	if err != nil {
		fmt.Printf("Failed to generate Html report file: %s\n\t%s", reporter.filename, err.Error())
	}

	if reporter.generateSummary {
		reporter.generateSummaryHtml()
	}
}

func (reporter *HtmlReporter) generateSummaryHtml() {
	if _, err := os.Stat(dir + "/summary.html"); err != nil && os.IsNotExist(err) {
		summary, err := os.Create(dir + "/summary.html")
		if err != nil {
			fmt.Printf("Failed to create summary Html report file\n\t%s", err.Error())
		}
		defer summary.Close()
		temp, err := template.New("SummaryHtmlTemplate").Parse(SummaryHtmlTemplate)
		if nil != err {
			fmt.Printf("Failed to parse summary Html template\n\t%s", err.Error())
		}
		err = temp.Execute(summary, reporter)
		if err != nil {
			fmt.Printf("Failed to generate summary Html report file\n\t%s", err.Error())
		}
	}

	summary, err := os.OpenFile(dir+"/summary.html", os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Failed to open summary Html report file\n\t%s", err.Error())
	}
	defer summary.Close()
	summary.Seek(-26, 2)
	temp, err := template.New("SummaryTemplate").Parse(SummaryTemplate)
	if nil != err {
		fmt.Printf("Failed to parse summary Html template\n\t%s", err.Error())
	}
	err = temp.Execute(summary, reporter.suite)
	if err != nil {
		fmt.Printf("Failed to generate summary Html report file\n\t%s", err.Error())
	}
}
