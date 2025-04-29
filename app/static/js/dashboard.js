if(sessionStorage.getItem("key") === null) {
	console.log("Key not found");
	window.location.href = "/vault";

} else {
	let key = sessionStorage.getItem("key");

	let stickiesLength = document.getElementsByClassName('sticky').length;

	let tempIV = document.getElementById("iv").value;

	let TotalBankAccountDollar = 0;


	for(let i = 0; i < stickiesLength; i++) {

		let tempT = document.getElementById('T-'+i).innerText;
		let tempD = document.getElementById('D-'+i).innerText;

		let tempET = decryptText(tempT, key, tempIV);
		document.getElementById('T-'+i).innerText = tempET;

		let tempED = decryptText(tempD, key, tempIV);
		document.getElementById('D-'+i).innerText = tempED;
	}

	/* Decrypt Bank Account Amount */
	let bankAccountLength = document.getElementsByClassName('bank_account').length;
	for(let i = 0; i < bankAccountLength; i++) {
		let tempAccount = document.getElementById('A-'+i).innerText;
		let tempAccountDecrypted = decryptText(tempAccount, key, tempIV);
		
		/* Add amount to TotalAmount */
		TotalBankAccountDollar += parseFloat(tempAccountDecrypted);

		document.getElementById('A-'+i).innerText = tempAccountDecrypted;
	}

	document.getElementById("BankAccountTotal").innerText = "Total Balance $" + TotalBankAccountDollar.toFixed(2);


	/* Decrypt CD */
	let CDLength = document.getElementsByClassName('CD').length;
	for(let i = 0; i < CDLength; i++) {
		let tempDeposit = document.getElementById('Deposit-'+i).innerText;
		let tempDepositDecrypted = decryptText(tempDeposit, key, tempIV);
		document.getElementById('Deposit-'+i).innerText = tempDepositDecrypted;

		let tempTerm = document.getElementById('Term-'+i).innerText;
                let tempTermDecrypted = decryptText(tempTerm, key, tempIV);
                document.getElementById('Term-'+i).innerText = tempTermDecrypted;

		let tempApy = document.getElementById('Apy-'+i).innerText;
                let tempApyDecrypted = decryptText(tempApy, key, tempIV);
                document.getElementById('Apy-'+i).innerText = tempApyDecrypted;



		let finalAmount = calculateCD(tempDepositDecrypted, tempApyDecrypted, tempTermDecrypted, 1);
		console.log(finalAmount);
		document.getElementById('FinalAmount-'+i).innerText = finalAmount;
	}

	let cardsLength = document.getElementsByClassName('my_card').length;

	for(let i = 0; i < cardsLength; i++) {
		
		/* Card Header */
		let tempCardBank = document.getElementById('Card-CardBank-'+i).innerText;
		let tempCardBankEncrypted = decryptText(tempCardBank, key, tempIV);
		document.getElementById('Card-CardBank-'+i).innerText = tempCardBankEncrypted;

		let tempCardName = document.getElementById('Card-CardName-'+i).innerText;
		let tempCardNameEncrypted = decryptText(tempCardName, key, tempIV);
		document.getElementById('Card-CardName-'+i).innerText = tempCardNameEncrypted;
		
		let tempBalance = document.getElementById('Card-Balance-'+i).innerText;
		let tempBalanceEncrypted = decryptText(tempBalance, key, tempIV);
		document.getElementById('Card-Balance-'+i).innerText = "Current Balance: $" + tempBalanceEncrypted;

		let tempDueDate = document.getElementById('Card-DueDate-'+i).innerText;
		let tempDueDateEncrypted = decryptText(tempDueDate, key, tempIV);
		document.getElementById('Card-DueDate-'+i).innerText = "Due on " + tempDueDateEncrypted;
	}
}
