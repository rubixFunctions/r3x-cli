const r3x = require('@rubixfunctions/r3x-js-sdk')

let schema
r3x.execute(function(){
	let response = {'message' : 'Hello r3x function'}
	return response 
}, schema)