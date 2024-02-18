async function encrypt() {
	const key = document.getElementById("password")
	window.crypto.getRan

}

async function encryptString(plainText, key) {
  const textEncoder = new TextEncoder();
  const encodedText = textEncoder.encode(plainText);

  const cryptoKey = await window.crypto.subtle.importKey(
    'raw',
    key,
    { name: 'AES-CBC' },
    false,
    ['encrypt']
  );

  const iv = window.crypto.getRandomValues(new Uint8Array(16));

  const encryptedData = await window.crypto.subtle.encrypt(
    { name: 'AES-CBC', iv: iv },
    cryptoKey,
    encodedText
  );

  const encryptedBytes = new Uint8Array(encryptedData);

  // Combine IV and encrypted data into a single array
  const result = new Uint8Array(iv.length + encryptedBytes.length);
  result.set(iv);
  result.set(encryptedBytes, iv.length);

  return result;
}

// Example usage
(async () => {
  const plainText = "Hello, world!";
  const key = new Uint8Array(32); // 256-bit key
  window.crypto.getRandomValues(key);

  const encryptedData = await encryptString(plainText, key);
  console.log("Encrypted data:", encryptedData);
})();
