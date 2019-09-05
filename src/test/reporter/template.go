package reporter

const (
	HtmlTemplate = `
<html>
	<head>
		<title>{{.Name}}</title>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
	</head>
	<body>
		<div>
			<h1>{{.Name}}</h1>
			<p><strong>RunTime: </strong>{{.RunTime}}s</p>
			<p><strong>Total Tests: </strong>{{.TotalNum}}</p>
			<p><strong>Failed Tests: </strong>{{.FailedNum}}</p>
		</div>
		<table border="1" width="90%">
			<tr>
				<th>Name</th>
				<th>State</th>
				<th>RunTime</th>
				<th>Detail</th>
			</tr>
			{{ range .TestCases }}
				{{ if ne .State "Passed" }}
					<tr style='color:red;'>
				{{ else }}
					<tr>
				{{ end }}
				<td>{{.Name}}</td>
				<td>{{.State}}</td>
				<td>{{.RunTime}}s</td>
				<td>{{.Detail}}</td>
			</tr>
			{{ end }}
		</table>
	</body>
</html>
`
	SummaryHtmlTemplate = `
<html>
<head>
    <title>Summary of Test Results</title>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
	<style type="text/css">
		/* gridtable */
		table.gridtable {
			font-family: verdana,arial,sans-serif;
			font-size:12px;
			color:#333333;
			border-width: 1px;
			border-color: #666666;
			border-collapse: collapse;
		}
		table.gridtable th {
			border-width: 1px;
			padding: 8px;
			border-style: solid;
			border-color: #666666;
			background-color: #dedede;
		}
		table.gridtable td {
			border-width: 1px;
			padding: 8px;
			border-style: solid;
			border-color: #666666;
			background-color: #ffffff;
		}
	/* /gridtable */
	</style>
</head>
<body>
    <h2>Summary of Test Results</h2>
	<table class="gridtable">
		<tr>
			<th>Test Name</th>
			<th>Test Status</th>
			<th>Total Check Num</th>
			<th>Successful Check Num</th>
			<th>Failed Check Num</th>
			<th>Html Url</th>
		</tr>
	</table>
</body>
</html>
`

	SummaryTemplate = `
		<tr>
			<td>{{.Name}}</td>
			<td>{{.State}}</td>
			<td>{{.TotalNum}}</td>
			<td>{{.SuccessNum}}</td>
			<td>{{.FailedNum}}</td>
			<td><a href='{{.Url}}'>html</a></td>
		</tr>
	</table>
</body>
</html>
`
)
