
var api = api || function() {

	var status = function status() {
		return request("/status");
	};

	var request = function request(url) {
		return $.ajax({
	     	url: url,
	    	dataType: 'json',
	    	type: 'GET',
	      	contentType: "application/json",
	    });
	};

	return {
		status: status,
	};
}();