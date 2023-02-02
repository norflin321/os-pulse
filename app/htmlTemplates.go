package main

const baseHtml = `
	<!doctype html>
	<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width,initial-scale=1">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<title>os_pulse</title>
			<link href="index.css" rel="stylesheet" />
		</head>
	<body>
		<div class="container">%s</div>
		<a class="github-icon" href="https://github.com/norflin321">
			<img src="assets/github.svg"/>
		</a>
	</body>
	</html>
`

const columnHtml = `
<div class="column-head">
	<h1>%s</h1>
	<a class="goto-icon" href="%s">
		<img src="assets/goto.svg"/>
	</a>
</div>
<div class="items">
	%s
</div>
`

const githubItemHtml = `
<div class="item">
	<a class="title" href="%s">%s</a>
	<div class="desc">%s</div>
	<div class="info">
		%s
		<div class="stars">
			<img class="icon" src="assets/star.svg"/>
			<div class="text">%s</div>
		</div>
		<div class="forks">
			<img class="icon" src="assets/fork.svg"/>
			<div class="text">%s</div>
		</div>
		<div class="stars-today">
			<img class="icon" src="assets/star.svg"/>
			<div class="text">%s today</div>
		</div>
	</div>
</div>
`

const hackerNewsItemHtml = `
<div class="item">
	<a class="title" href="%s">%s</a>
	<div class="info">%s</div>
</div>
`
