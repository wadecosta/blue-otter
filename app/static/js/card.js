function postDeleteCard(button_id) {
                        
	(fetch('/delCard' , {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ button_id: button_id })
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('Network response was not ok');
			}
			return response.json();
		})
		.then(data => {
			console.log(data);
		})
		.catch(error => {
			console.error('There was a problem with the fetch operation:', error);
	}));

	console.log(button_id)
}
