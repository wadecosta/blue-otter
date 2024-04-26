function decryptText(encrypted_text, key_value, iv_value) {
	let key = CryptoJS.enc.Utf8.parse(key_value);
	let iv = CryptoJS.enc.Utf8.parse(iv_value);

	let decrypted = CryptoJS.AES.decrypt(encrypted_text, key, {
		iv: iv,
		mode: CryptoJS.mode.CBC,
		padding: CryptoJS.pad.Pkcs7
	});

	return decrypted.toString(CryptoJS.enc.Utf8);
}
