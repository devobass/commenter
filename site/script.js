var toInsert = []
var result

async function getData() {
	const url = "/fetch";
	try {
		const response = await fetch(url);
		if (!response.ok) {
			toInsert = `${response.status}, unable to retrieve comments`
			throw new Error(`Response status: ${response.status}`);
		}

		result = await response.json();
	} catch (error) {
		console.error(error.message);
	}
}

function appendCommentToSite(comment) {
	const unix = new Date(comment.time * 1000);
	const date = unix.toLocaleString();

	const p = document.createElement('p');
	const user = document.createElement('span');
	const msg = document.createElement('span');
	const time = document.createElement('span');

	p.className = "cmt"
	user.className = "username"
	msg.className = "message"
	time.className = "time"

	user.textContent = `${comment.username}`
	msg.textContent = `${comment.message}`
	time.textContent = date;

	p.appendChild(user);
	p.appendChild(time);
	p.appendChild(msg);

	toInsert.push(p);
}

(async () => {
	await getData();
	if (result.length == 0) {
		const board = document.getElementById("board").innerHTML = "<p>such emptiness.</p>";
		return
	}
	result.reverse().forEach(appendCommentToSite)
	const board = document.getElementById("board");
	board.replaceChildren();
	
	board.append(...toInsert);
})();
