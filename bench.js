
var http = require('http');
var https = require('https');

var count = 0;

var options = {
	hostname: '10.0.0.125',
	port: 443,
	path: '/',
	method: 'GET',
	headers: {'Connection' : 'Close'},
	agent: null,
	rejectUnauthorized: false,
};

var f = function() {
	var req = https.request(options, function(res) {
		res_body = '';
		res.on('data', function(chunk) {
			res_body += chunk;
		});
		res.on('end', function() {
			count++;
			f();
		});
	});
	req.on('error', function(e) {
		console.log('error' + e.message);
		setTimeout(f, 1000);
	});
	req.end();
};

for (var i = 0; i < 10; i++) {
	setTimeout(f, 10);
}

var c = function() {
	console.log('count:' + count);
	count = 0;
	setTimeout(c, 1000);
};

setTimeout(c, 1000);

