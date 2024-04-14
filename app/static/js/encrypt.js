//function encryptText() {

	/* Title */
//	let title_plaintext = document.getElementById('title').value;

//	let key_value = sessionStorage.getItem("key");
//	let key = CryptoJS.enc.Utf8.parse(key_value); // 256-bit key

//	let iv_value = document.getElementById('iv').value;
//	let iv = CryptoJS.enc.Utf8.parse(iv_value); // 128-bit IV

//	let encrypted = CryptoJS.AES.encrypt(title_plaintext, key, {
//		iv: iv,
//		mode: CryptoJS.mode.CBC,
//		padding: CryptoJS.pad.Pkcs7
//	});

//	document.getElementById('encryptedTitle').value = encrypted.toString();
//	console.log(encrypted.toString());
//}

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
	/* Title Encryption */
	let title_plaintext = document.getElementById('title').value;
	let key_value = sessionStorage.getItem("key");
	let iv_value = document.getElementById('iv').value;

	console.log(iv_value);

	let encrypted_title = encryptText(title_plaintext, key_value, iv_value);
	document.getElementById('encryptedTitle').value = encrypted_title.toString();
	console.log(encrypted_title.toString());
}
