function getEventCount() {
	return document.getElementsByTagName('p').length;
}


if(sessionStorage.getItem("key") === null) {
	console.log("Key not found");
	window.location.href = "/vault";

} else {
	let key = sessionStorage.getItem("key");
	console.log("The key has been passed:" + key);

	let items = getEventCount();
	console.log("Total Item:" + items);
}
