var five = require("johnny-five");
var mongoose = require('mongoose');

var board = new five.Board();

// Connect to MongoDB
var uri = 'mongodb://username:password@server:port/db';
db = mongoose.connect(uri);
Schema = mongoose.Schema;

// Define New Data Schema for MongoDB
var userSchema = new Schema({
	temperature : { type: Number, default: 0 },
	date: { type: Date, default: Date.now() }},
	{ collection: "temperatureData" });

// Assign Schema to Model
var userModel = mongoose.model('Temp', userSchema);

// Start Arduino Board
board.on("ready", function() {

	// Define Temperature Sensor Output Pin
	var temp = new five.Temperature({
		pin: "A0",
		controller: "TMP36"
	});

	// Set Interval of 60 Seconds
	setInterval(function(){
		// Output Temp. in Console
		console.log("Temperature: %d", temp.celsius);
		// Check Temperature & add Date
		var test = new userModel({temperature: temp.celsius, date : Date.now()});
		// Save information to MongoDB
		test.save(function (err, test) {
			if (err) {
				console.log("error");
				return console.error(err);
			}
		});
	}, 60000); // Interval of 60 Seconds
});