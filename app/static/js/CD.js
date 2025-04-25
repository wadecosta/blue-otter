function postDeleteCD(button_id) {

	(fetch('/delCD', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ button_id : button_id })
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('Network response was not ok');
			}
			return response.json();
		})
		.then(data => {
			location.reload(true);
		})
		.catch(error => {
			console.error('There was a problem with the fetch operation:', error)
		})
	);
}

function calculateCD(principal, annualRate, months, compoundFrequency) {

	console.log("principal :", principal, " annualRate :", annualRate, " months:", months, " compoundFrequency:", compoundFrequency);

	/* Convert annual interest rate to a decimal */
	const rate = annualRate / 100;

	console.log("rate :", rate);

	/* Convert months into years */
	const years = months / 12;

	console.log("years :", years);

	/* Calculate the total number of compounding periods */
	const periods = compoundFrequency * years;

	console.log("periods :", periods);

	/* Calculate the interest rate per compounding period */
	const ratePerPeriod = rate / compoundFrequency;

	console.log("ratePerPeriod :", ratePerPeriod);
	
	/* Calculate the future value using the compound interest formula */
	const futureValue = principal * Math.pow(1 + ratePerPeriod, periods);

	console.log("futureValue", futureValue);

	/* Return the future value rounded to 2 decimal places */
	return futureValue.toFixed(2);
}
