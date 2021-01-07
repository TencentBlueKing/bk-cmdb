package md2man

import (
	"testing"

	"github.com/russross/blackfriday/v2"
)

type TestParams struct {
	extensions blackfriday.Extensions
}

func TestEmphasis(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		".nh\n\n.PP\nnothing inline\n",

		"simple *inline* test\n",
		".nh\n\n.PP\nsimple \\fIinline\\fP test\n",

		"*at the* beginning\n",
		".nh\n\n.PP\n\\fIat the\\fP beginning\n",

		"at the *end*\n",
		".nh\n\n.PP\nat the \\fIend\\fP\n",

		"*try two* in *one line*\n",
		".nh\n\n.PP\n\\fItry two\\fP in \\fIone line\\fP\n",

		"over *two\nlines* test\n",
		".nh\n\n.PP\nover \\fItwo\nlines\\fP test\n",

		"odd *number of* markers* here\n",
		".nh\n\n.PP\nodd \\fInumber of\\fP markers* here\n",

		"odd *number\nof* markers* here\n",
		".nh\n\n.PP\nodd \\fInumber\nof\\fP markers* here\n",

		"simple _inline_ test\n",
		".nh\n\n.PP\nsimple \\fIinline\\fP test\n",

		"_at the_ beginning\n",
		".nh\n\n.PP\n\\fIat the\\fP beginning\n",

		"at the _end_\n",
		".nh\n\n.PP\nat the \\fIend\\fP\n",

		"_try two_ in _one line_\n",
		".nh\n\n.PP\n\\fItry two\\fP in \\fIone line\\fP\n",

		"over _two\nlines_ test\n",
		".nh\n\n.PP\nover \\fItwo\nlines\\fP test\n",

		"odd _number of_ markers_ here\n",
		".nh\n\n.PP\nodd \\fInumber of\\fP markers\\_ here\n",

		"odd _number\nof_ markers_ here\n",
		".nh\n\n.PP\nodd \\fInumber\nof\\fP markers\\_ here\n",

		"mix of *markers_\n",
		".nh\n\n.PP\nmix of *markers\\_\n",

		"*What is A\\* algorithm?*\n",
		".nh\n\n.PP\n\\fIWhat is A* algorithm?\\fP\n",
	}
	doTestsInline(t, tests)
}

func TestStrong(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		".nh\n\n.PP\nnothing inline\n",

		"simple **inline** test\n",
		".nh\n\n.PP\nsimple \\fBinline\\fP test\n",

		"**at the** beginning\n",
		".nh\n\n.PP\n\\fBat the\\fP beginning\n",

		"at the **end**\n",
		".nh\n\n.PP\nat the \\fBend\\fP\n",

		"**try two** in **one line**\n",
		".nh\n\n.PP\n\\fBtry two\\fP in \\fBone line\\fP\n",

		"over **two\nlines** test\n",
		".nh\n\n.PP\nover \\fBtwo\nlines\\fP test\n",

		"odd **number of** markers** here\n",
		".nh\n\n.PP\nodd \\fBnumber of\\fP markers** here\n",

		"odd **number\nof** markers** here\n",
		".nh\n\n.PP\nodd \\fBnumber\nof\\fP markers** here\n",

		"simple __inline__ test\n",
		".nh\n\n.PP\nsimple \\fBinline\\fP test\n",

		"__at the__ beginning\n",
		".nh\n\n.PP\n\\fBat the\\fP beginning\n",

		"at the __end__\n",
		".nh\n\n.PP\nat the \\fBend\\fP\n",

		"__try two__ in __one line__\n",
		".nh\n\n.PP\n\\fBtry two\\fP in \\fBone line\\fP\n",

		"over __two\nlines__ test\n",
		".nh\n\n.PP\nover \\fBtwo\nlines\\fP test\n",

		"odd __number of__ markers__ here\n",
		".nh\n\n.PP\nodd \\fBnumber of\\fP markers\\_\\_ here\n",

		"odd __number\nof__ markers__ here\n",
		".nh\n\n.PP\nodd \\fBnumber\nof\\fP markers\\_\\_ here\n",

		"mix of **markers__\n",
		".nh\n\n.PP\nmix of **markers\\_\\_\n",

		"**`/usr`** : this folder is named `usr`\n",
		".nh\n\n.PP\n\\fB\\fB\\fC/usr\\fR\\fP : this folder is named \\fB\\fCusr\\fR\n",

		"**`/usr`** :\n\n this folder is named `usr`\n",
		".nh\n\n.PP\n\\fB\\fB\\fC/usr\\fR\\fP :\n\n.PP\nthis folder is named \\fB\\fCusr\\fR\n",
	}
	doTestsInline(t, tests)
}

func TestEmphasisMix(t *testing.T) {
	var tests = []string{
		"***triple emphasis***\n",
		".nh\n\n.PP\n\\fB\\fItriple emphasis\\fP\\fP\n",

		"***triple\nemphasis***\n",
		".nh\n\n.PP\n\\fB\\fItriple\nemphasis\\fP\\fP\n",

		"___triple emphasis___\n",
		".nh\n\n.PP\n\\fB\\fItriple emphasis\\fP\\fP\n",

		"***triple emphasis___\n",
		".nh\n\n.PP\n***triple emphasis\\_\\_\\_\n",

		"*__triple emphasis__*\n",
		".nh\n\n.PP\n\\fI\\fBtriple emphasis\\fP\\fP\n",

		"__*triple emphasis*__\n",
		".nh\n\n.PP\n\\fB\\fItriple emphasis\\fP\\fP\n",

		"**improper *nesting** is* bad\n",
		".nh\n\n.PP\n\\fBimproper *nesting\\fP is* bad\n",

		"*improper **nesting* is** bad\n",
		".nh\n\n.PP\n*improper \\fBnesting* is\\fP bad\n",
	}
	doTestsInline(t, tests)
}

func TestCodeSpan(t *testing.T) {
	var tests = []string{
		"`source code`\n",
		".nh\n\n.PP\n\\fB\\fCsource code\\fR\n",

		"` source code with spaces `\n",
		".nh\n\n.PP\n\\fB\\fCsource code with spaces\\fR\n",

		"` source code with spaces `not here\n",
		".nh\n\n.PP\n\\fB\\fCsource code with spaces\\fRnot here\n",

		"a `single marker\n",
		".nh\n\n.PP\na `single marker\n",

		"a single multi-tick marker with ``` no text\n",
		".nh\n\n.PP\na single multi\\-tick marker with ``` no text\n",

		"markers with ` ` a space\n",
		".nh\n\n.PP\nmarkers with  a space\n",

		"`source code` and a `stray\n",
		".nh\n\n.PP\n\\fB\\fCsource code\\fR and a `stray\n",

		"`source *with* _awkward characters_ in it`\n",
		".nh\n\n.PP\n\\fB\\fCsource *with* \\_awkward characters\\_ in it\\fR\n",

		"`split over\ntwo lines`\n",
		".nh\n\n.PP\n\\fB\\fCsplit over\ntwo lines\\fR\n",

		"```multiple ticks``` for the marker\n",
		".nh\n\n.PP\n\\fB\\fCmultiple ticks\\fR for the marker\n",

		"```multiple ticks `with` ticks inside```\n",
		".nh\n\n.PP\n\\fB\\fCmultiple ticks `with` ticks inside\\fR\n",
	}
	doTestsInline(t, tests)
}

func TestListLists(t *testing.T) {
	var tests = []string{
		"\n\n**[grpc]**\n: Section for gRPC socket listener settings. Contains three properties:\n - **address** (Default: \"/run/containerd/containerd.sock\")\n - **uid** (Default: 0)\n - **gid** (Default: 0)",
		".nh\n\n.TP\n\\fB[grpc]\\fP\nSection for gRPC socket listener settings. Contains three properties:\n.RS\n.IP \\(bu 2\n\\fBaddress\\fP (Default: \"/run/containerd/containerd.sock\")\n.IP \\(bu 2\n\\fBuid\\fP (Default: 0)\n.IP \\(bu 2\n\\fBgid\\fP (Default: 0)\n\n.RE\n\n",
	}
	doTestsParam(t, tests, TestParams{blackfriday.DefinitionLists})
}

func TestLineBreak(t *testing.T) {
	var tests = []string{
		"this line  \nhas a break\n",
		".nh\n\n.PP\nthis line\n.br\nhas a break\n",

		"this line \ndoes not\n",
		".nh\n\n.PP\nthis line\ndoes not\n",

		"this line\\\ndoes not\n",
		".nh\n\n.PP\nthis line\\\\\ndoes not\n",

		"this line\\ \ndoes not\n",
		".nh\n\n.PP\nthis line\\\\\ndoes not\n",

		"this has an   \nextra space\n",
		".nh\n\n.PP\nthis has an\n.br\nextra space\n",
	}
	doTestsInline(t, tests)

	tests = []string{
		"this line  \nhas a break\n",
		".nh\n\n.PP\nthis line\n.br\nhas a break\n",

		"this line \ndoes not\n",
		".nh\n\n.PP\nthis line\ndoes not\n",

		"this line\\\nhas a break\n",
		".nh\n\n.PP\nthis line\n.br\nhas a break\n",

		"this line\\ \ndoes not\n",
		".nh\n\n.PP\nthis line\\\\\ndoes not\n",

		"this has an   \nextra space\n",
		".nh\n\n.PP\nthis has an\n.br\nextra space\n",
	}
	doTestsInlineParam(t, tests, TestParams{
		extensions: blackfriday.BackslashLineBreak})
}

func TestTable(t *testing.T) {
	var tests = []string{
		`
| Animal               | Color         |
| --------------| --- |
| elephant        | Gray. The elephant is very gray.  |
| wombat     | No idea.      |
| zebra        | Sometimes black and sometimes white, depending on the stripe.     |
| robin | red. |
`,
		`.nh

.TS
allbox;
l l 
l l .
\fB\fCAnimal\fR	\fB\fCColor\fR
elephant	T{
Gray. The elephant is very gray.
T}
wombat	No idea.
zebra	T{
Sometimes black and sometimes white, depending on the stripe.
T}
robin	red.
.TE
`,
	}
	doTestsInlineParam(t, tests, TestParams{blackfriday.Tables})
}

func TestLinks(t *testing.T) {
	var tests = []string{
		"See [docs](https://docs.docker.com/) for\nmore",
		".nh\n\n.PP\nSee docs\n\\[la]https://docs.docker.com/\\[ra] for\nmore\n",
	}
	doTestsInline(t, tests)
}

func execRecoverableTestSuite(t *testing.T, tests []string, params TestParams, suite func(candidate *string)) {
	// Catch and report panics. This is useful when running 'go test -v' on
	// the integration server. When developing, though, crash dump is often
	// preferable, so recovery can be easily turned off with doRecover = false.
	var candidate string
	const doRecover = true
	if doRecover {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("\npanic while processing [%#v]: %s\n", candidate, err)
			}
		}()
	}
	suite(&candidate)
}

func runMarkdown(input string, params TestParams) string {
	renderer := NewRoffRenderer()
	return string(blackfriday.Run([]byte(input), blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(params.extensions)))
}

func doTestsParam(t *testing.T, tests []string, params TestParams) {
	execRecoverableTestSuite(t, tests, params, func(candidate *string) {
		for i := 0; i+1 < len(tests); i += 2 {
			input := tests[i]
			*candidate = input
			expected := tests[i+1]
			actual := runMarkdown(*candidate, params)
			if actual != expected {
				t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
					*candidate, expected, actual)
			}

			// now test every substring to stress test bounds checking
			if !testing.Short() {
				for start := 0; start < len(input); start++ {
					for end := start + 1; end <= len(input); end++ {
						*candidate = input[start:end]
						runMarkdown(*candidate, params)
					}
				}
			}
		}
	})
}

func doTestsInline(t *testing.T, tests []string) {
	doTestsInlineParam(t, tests, TestParams{})
}

func doTestsInlineParam(t *testing.T, tests []string, params TestParams) {
	params.extensions |= blackfriday.Strikethrough
	doTestsParam(t, tests, params)
}
