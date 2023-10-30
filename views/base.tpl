{{ define "base" }}<!DOCTYPE html>
<html lang="en">
<head>
	<link rel="icon" type="image/png" href="{{ .BaseUri }}/favicon.png">
	
	<!-- bootstrap -->
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="{{ .BaseUri }}/stubman/static/css/bootstrap.min.css" />
	<link rel="stylesheet" href="{{ .BaseUri }}/stubman/static/css/bootstrap-theme.min.css" />
	<link rel="stylesheet" href="{{ .BaseUri }}/stubman/static/css/style.css" />

    {{ template "title" . }}
</head>
<body>
    {{ template "sidebar" . }}
	<div class="container">
	<h1>Stubman</h1>
	<ul class="nav nav-tabs">
	  <li role="presentation" {{ if .HomePage }}class="active"{{ end }}><a href="{{ .BaseUri }}/stubman/">Home</a></li>
	  <li role="presentation" {{ if .CreatePage }}class="active"{{ end }}><a href="{{ .BaseUri }}/stubman/create/">New</a></li>
	</ul>
    {{ template "content" . }}
	<footer class="main-footer navbar navbar-fixed-bottom">
	Victor K&copy; 2016
	</footer>
	</div>

	<script src="{{ .BaseUri }}/stubman/static/js/jquery.min.js" crossorigin="anonymous"></script>
	<script src="{{ .BaseUri }}/stubman/static/js/bootstrap.min.js" crossorigin="anonymous"></script>
	{{ template "scripts" . }}
</body>
</html>
{{ end }}
// We define empty blocks for optional content so we don't have to define a block in child templates that don't need them
{{ define "scripts" }}{{ end }}
{{ define "sidebar" }}{{ end }}