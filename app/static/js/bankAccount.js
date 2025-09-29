function postDeleteBankAccount(button_id) {
	
	(fetch('/delBankAccount', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ button_id: button_id })
		})
		.then(response => {
			if(!response.ok) {
				throw new Error('Network response was not okk');
			}
			return response.json();
		})
		.then(data => {
			console.log(data);
			location.reload(true);
		})
		.catch(error => {
			console.error('There was a problem with the fetch operation', error);
		})
	);
}
