<!DOCTYPE html>
<html lang="en">
<head>
	<link rel="icon" type="image/png" href="/favicon.png">
	
	<!-- bootstrap -->
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="{{ .BaseUri }}/stubman/static/css/bootstrap.min.css" />
	<link rel="stylesheet" href="{{ .BaseUri }}/stubman/static/css/bootstrap-theme.min.css" />
	<link rel="stylesheet" href="{{ .BaseUri }}/stubman/static/css/style.css" />

    <title>{{ .Title }}</title>
</head>
<body>
	<div class="container">
		<h1>{{ .Title }}</h1>
		<p class="alert alert-warning" role="alert">
		{{ .Message }}
		</p>
		<p>Go to <a href="{{ .BaseUri }}/stubman/">home page</a></li></p>
	
	<footer class="main-footer navbar navbar-fixed-bottom">
	Victor K&copy; 2016
	</footer>
	</div>
</body>
</html>