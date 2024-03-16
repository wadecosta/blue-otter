if(sessionStorage.getItem("key") === null) {
	console.log("Key not found");
} else {
	let key = sessionStorage.getItem("key");
	console.log("The key has been passed:" + key);
}
