function initEdit() {
	var reqHeader = $('#header-request-tpl').html()
	var reqHeadersList = $('#request-headers')
	var reqUriEl = $('#model-request_uri')
	var reqBodyEl = $('#model-request_body')
	var curlPreviewEl = $('#curl-preview')

	$('#request-add-header').click(function(e) {
		e.preventDefault();
		reqHeadersList.append(reqHeader)
		updateCurlPreview()
	});
	
	var respHeader = $('#header-response-tpl').html()
	var respList = $('#response-headers')
	$('#response-add-header').click(function(e) {
		e.preventDefault();
		respList.append(respHeader)
	});
	
	var method = $('#model-request_method')
	var req_container = $('#request_body_container')
	if (method.val() == 'GET' || method.val() == 'ANY') {
		req_container.find('textarea').attr('disabled', 'disabled')
	} else {
		req_container.find('textarea').removeAttr('disabled')
	}
	method.change(function(e) {
		e.preventDefault();
		if (method.val() == 'GET' || method.val() == 'ANY') {
			req_container.find('textarea').attr('disabled', 'disabled')
		} else {
			req_container.find('textarea').removeAttr('disabled')
		}

		updateCurlPreview()
	});

	reqUriEl.change(function() {
		updateCurlPreview()
	})
	
	$('.decoder').click(function(e) {
		e.preventDefault();
		var el = $(this).closest('.form-group').find('textarea');
		try {
			var parsed = JSON.parse(el.val().trim())
			el.val(JSON.stringify(parsed, null, '  '))	
		} catch(e) {
			console.log(e)
		}

		updateCurlPreview()
	});
	
	$('.encoder').click(function(e) {
		e.preventDefault();
		var el = $(this).closest('.form-group').find('textarea');
		try {
			var parsed = JSON.parse(el.val().trim())
			el.val(JSON.stringify(parsed))
			updateCurlPreview()
		} catch(e) {
			console.log(e)
		}
	});
	
	$('.headers-group').on('click', '.btn-del-header', function(e) {
		e.preventDefault();
		$(this).parent().remove()

		updateCurlPreview()
	});

	$('.headers-group').on('change', 'request[headers][]', function(e) {
		updateCurlPreview()
	});

	$('#curl-copy-btn').click(function() {
		updateCurlPreview();

		var el = document.createElement('textarea')
		el.value = curlPreviewEl.text()

		// Select the text field
		el.select();
		if (el.setSelectionRange) {
			el.setSelectionRange(0, 99999); // For mobile devices
		}
	  
		 // Copy the text inside the text field
		navigator.clipboard.writeText(el.value);

		var currEl = this.querySelector('.b-label')

		var savedText = currEl.innerText;
		currEl.innerText = 'Copied to clipboard!'
		this.disabled = true;
		this.classList.replace("btn-secondary", "btn-success")

		// rollback name
		setTimeout(
			() => { 
				currEl.innerText = savedText; 
				this.disabled = false;
				this.classList.replace("btn-success", "btn-secondary")
			}, 
			2000
		);
	});

	updateCurlPreview();

	function updateCurlPreview() {
		var params = {
			url: location.protocol + '//' + location.host + reqUriEl.val(),
			method: method.val() === 'ANY' || !method.val() ? 'GET' : method.val(),
			headers: {}
		}

		const options = {
			include: true // include response headers
		}

		var hasContentType = false;
		$('.headers-group').find('input[name="request[headers][]"]').each(function (offset, el) {
			var header = this.value.trim()
			if (header) { // skip empty
				var index = header.indexOf(':');
				if (index !== -1) {
					var headerName = header.slice(0, index).trim().toLowerCase()
					if (headerName === 'content-type') {
						hasContentType = true;
					}

					params.headers[header.slice(0, index).trim()] = header.slice(index + 1).trim()
				}
			}
		});
		
		if (!reqBodyEl.disabled) {
			try {
				if (!hasContentType) {
					params.headers['Content-Type'] = 'application/json' // force JSON content type
				}

				var parsed 
				
				if (!reqBodyEl.val()) {
					params.body = '{}'
				} else {
					parsed = JSON.parse(reqBodyEl.val().trim())
					params.body = JSON.stringify(parsed)
				}
			} catch(e) {
				console.log("Body parse error:", e)
			}
		}

		if (location.protocol === 'https:') {
			options.insecure = true
		}
		var curlText = window['curl-generator'].CurlGenerator(params, options)

		curlPreviewEl.text(curlText, options);
	}
}
