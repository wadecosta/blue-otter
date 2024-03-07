/* Create AES key */
async function generateAESKey(passphrase) {
  // Convert passphrase to ArrayBuffer
  const encoder = new TextEncoder();
  const passphraseBuffer = encoder.encode(passphrase);

  // Generate salt (random bytes)
  const salt = crypto.getRandomValues(new Uint8Array(16));

  // Derive key using PBKDF2
  const derivedKey = await crypto.subtle.importKey(
    "raw",
    passphraseBuffer,
    { name: "PBKDF2" },
    false,
    ["deriveKey"]
  ).then(baseKey => crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      salt: salt,
      iterations: 1000000,
      hash: "SHA-256"
    },
    baseKey,
    { name: "AES-GCM", length: 256 },
    true,
    ["encrypt", "decrypt"]
  ));

  // Export the derived key as ArrayBuffer
  const aesKeyArrayBuffer = await crypto.subtle.exportKey("raw", derivedKey);

  // Convert ArrayBuffer to hex string
  const aesKeyHex = Array.from(new Uint8Array(aesKeyArrayBuffer))
    .map(byte => byte.toString(16).padStart(2, '0')).join('');

  return aesKeyHex;
}






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


// Example usage
/*const passphrase = "YourStrongPassphraseHere";

generateAESKey(passphrase)
  .then(aesKeyHex => {
    console.log("AES Key:", aesKeyHex);
  })
  .catch(error => {
    console.error("Error:", error);
  });*/
