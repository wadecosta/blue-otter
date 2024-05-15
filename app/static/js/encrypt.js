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

function encrypt_sticky() {
	
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

	return false;
}

function encrypt_card() {
	
	event.preventDefault();

	/* General Encryption Values */
	let key_value = sessionStorage.getItem("key");
        let iv_value = document.getElementById('iv').value;

	/* Card Bank Encryption */
	let card_bank_plaintext = document.getElementById('card_bank').value;
	let encrypted_card_bank = encryptText(card_bank_plaintext, key_value, iv_value);
	//document.getElementById('encryptedCardBank').value = encrypted_card_bank.toString();

	/* Card Name Encryption */
	let card_name_plaintext = document.getElementById('card_name').value;
        let encrypted_card_name = encryptText(card_name_plaintext, key_value, iv_value);
        //document.getElementById('encryptedCardName').value = encrypted_card_name.toString();

	/* Balance Encryption */
	let balance_plaintext = document.getElementById('balance').value;
        let encrypted_balance = encryptText(balance_plaintext, key_value, iv_value);
        //document.getElementById('encryptedBalance').value = encrypted_balance.toString();

	/* Due Date Encryption */
	let due_date_plaintext = document.getElementById('due_date').value;
        let encrypted_due_date= encryptText(due_date_plaintext, key_value, iv_value);
        //document.getElementById('encryptedDueDate').value = encrypted_due_date.toString();
	
	const xhr = new XMLHttpRequest();
	xhr.open("POST", window.location.href);
	xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8")

	const body = JSON.stringify({
		card_bank: encrypted_card_bank.toString(),
		card_name: encrypted_card_name.toString(),
		balance: encrypted_balance.toString(),
		due_date: encrypted_due_date.toString()
	});

	xhr.send(body);
	return false;
}
