var currentTime = new Date();

var hours = currentTime.getHours();
var message = "Good ";

if(hours < 12) {
	message += "Morning 🌇";
}
else if(hours < 17) {
	message += "Afternoon 🌆";
}
else {
	message += "Evening 🌃";
}

const element = document.getElementById("good_message");
element.innerHTML = message;
