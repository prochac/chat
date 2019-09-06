let socket = new WebSocket((location.protocol === 'https:' ? "wss://" : "ws://") + window.location.host + "/ws");
socket.onopen = e => console.log("[open] Connection established");
socket.onmessage = e => {
    const li = document.createElement("LI");
    li.setAttribute("class", "list-group-item");
    li.appendChild(document.createTextNode(e.data));
    messages.appendChild(li);
    return messages.parentNode.scrollTo(0, messages.scrollHeight);
};
socket.onclose = event => event.wasClean
    ? console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`)
    : console.log('[close] Connection died');
socket.onerror = err => console.log(`[error] ${err.message}`);


let messages = document.getElementById("messages");
const message = document.getElementById("message");
const form = document.getElementById("form");

form.addEventListener("submit", ev => {
    ev.preventDefault();

    const json = JSON.stringify(message.value);
    message.value = '';

    return fetch("/post-message", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: json,
    });
});

reloadMessages = () => fetch("/get-messages")
    .then(response => response.json())
    .then(json => {
        const newMessages = document.createElement("UL");
        newMessages.setAttribute("class", "list-group");
        newMessages.setAttribute("id", "messages");

        json.forEach(message => {

            const li = document.createElement("LI");
            li.setAttribute("class", "list-group-item");
            li.appendChild(document.createTextNode(message));
            return newMessages.appendChild(li);
        });

        messages.parentNode.replaceChild(newMessages, messages);
        messages = document.getElementById("messages");
    });

document.getElementById("app").onload = () => {
    reloadMessages().then(() => 
   		messages.parentNode.scrollTo(0, messages.scrollHeight));
    setInterval(reloadMessages, 5000);
};

