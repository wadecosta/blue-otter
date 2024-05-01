function encryptText(title_plaintext, key_value, iv_value) {

        /* Title */
        //let title_plaintext = document.getElementById('title').value;

        //let key_value = sessionStorage.getItem("key");
        let key = CryptoJS.enc.Utf8.parse(key_value); // 256-bit key

        //let iv_value = document.getElementById('iv').value;
        let iv = CryptoJS.enc.Utf8.parse(iv_value); // 128-bit IV

        let encrypted = CryptoJS.AES.encrypt(title_plaintext, key, {
                iv: iv,
                mode: CryptoJS.mode.CBC,
                padding: CryptoJS.pad.Pkcs7
        });

	return encrypted.toString();
}

function encrypt() {
	
	event.preventDefault();
	
	/* Title Encryption */
	let title_plaintext = document.getElementById('title').value;
	let key_value = sessionStorage.getItem("key");
	let iv_value = document.getElementById('iv').value;

	console.log(iv_value.toString());

	let encrypted_title = encryptText(title_plaintext, key_value, iv_value);
	document.getElementById('encryptedTitle').value = encrypted_title.toString();
	console.log(encrypted_title.toString());

	/* Description Encryption */
	let description_plaintext = document.getElementById('description').value;
	let encrypted_description = encryptText(description_plaintext, key_value, iv_value);
	document.getElementById('encryptedDescription').value = encrypted_description.toString();
	console.log(encrypted_description.toString());

	console.log(document.getElementById('encryptedDescription').value);
	console.log(document.getElementById('encryptedTitle').value);

	console.log("Form will be submitted now.");

	console.log(document.forms[0]);

	
	const xhr = new XMLHttpRequest();
	xhr.open("POST", window.location.href);
	xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8")

	const body = JSON.stringify({
		title: encrypted_title.toString(),
		description: encrypted_description.toString()
	});

	xhr.send(body);


	//document.forms[0].submit();
	
	return false;
}
