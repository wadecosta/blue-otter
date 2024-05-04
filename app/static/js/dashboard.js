function getEventCount() {
	return document.getElementsByClassName('sticky').length;
}


if(sessionStorage.getItem("key") === null) {
	console.log("Key not found");
	window.location.href = "/vault";

} else {
	let key = sessionStorage.getItem("key");

	let items = getEventCount();

	let tempIV = document.getElementById("iv").value;


	for(let i = 0; i < items; i++) {

		let tempT = document.getElementById('T-'+i).innerText;
		let tempD = document.getElementById('D-'+i).innerText;

		let tempET = decryptText(tempT, key, tempIV);
		document.getElementById('T-'+i).innerText = tempET;

		let tempED = decryptText(tempD, key, tempIV);
		document.getElementById('D-'+i).innerText = tempED;
	}
}
