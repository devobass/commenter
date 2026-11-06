var toInsert = []
var result

const button = document.getElementById('comment-button');
button.addEventListener('click', async (e) => {
	 e.preventDefault();
	const username = document.getElementById("username").value;
	const message = document.getElementById("message").value; 

	const payload = { username, message };

	try {
		const response = await fetch("/comment", {
			method: "POST",
			body: JSON.stringify(payload),
			headers: { "Content-Type": "application/json; charset=UTF-8" }
		});

		if (response.ok) {
			window.location.reload();
		} else {
			console.error("Server returned an error status:", response.status);
		}
	} catch (error) {
		console.error("Network error occurred:", error);
	}
});

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
	const id = document.createElement('id');

	p.className = "cmt"
	user.className = "username"
	msg.className = "message"
	time.className = "time"
	id.className = "msgid"

	user.textContent = `${comment.username}`
	msg.textContent = `${comment.message}`
	id.textContent = `${comment.id}`
	time.textContent = date;

	p.appendChild(user);
	p.appendChild(time);
	p.appendChild(msg);
	p.appendChild(id);

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
