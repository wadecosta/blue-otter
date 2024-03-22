function encryptStringAES(inputString, passphrase) {
	var encrypted = CryptoJS.AES.encrypt(inputString, passphrase);
	return encrypted.toString();
}

function decryptStringAES(encryptedString, passphrase) {
	var decrypted = CryptoJS.AES.decrypt(encryptedString, passphrase);
	var plaintext = decrypted.toString(CryptoJS.enc.Utf8);
	return plaintext;
}

var plaintext = "Hola mi amigos";
var passphrase = "p"

var encryptedString = encryptStringAES(plaintext, passphrase);
console.log("Encrypted String:", encryptedString);

var decryptedString = decryptStringAES(encryptedString, passphrase);
console.log("Decrypted String:", decryptedString);
