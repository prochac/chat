createMessage = message => {
    const li = document.createElement("LI");
    li.setAttribute("class", "list-group-item d-flex justify-content-between align-items-center");
    li.append(document.createTextNode(message.text));


    const span = document.createElement("SPAN");
    span.setAttribute("class", "badge badge-secondary badge-pill");
    span.append(document.createTextNode(moment(message.timestamp).fromNow()));
    li.append(span)

    return li
}

showAlert = err => {
    const div = document.createElement("DIV");
    div.setAttribute("class", "w-100 alert alert-danger");
    div.setAttribute("role", "alert");
    div.append(document.createTextNode(err));
    document.getElementById("container").prepend(div)
    setTimeout(() => {
        div.parentNode.removeChild(div);
    }, 5000);
}

let socket = new WebSocket((location.protocol === 'https:' ? "wss://" : "ws://") + window.location.host + "/ws");
socket.onopen = e => console.log("[open] Connection established");
socket.onmessage = e => {
    const msg = JSON.parse(e.data);
    messages.append(createMessage(msg));
    return messages.parentNode.scrollTo(0, messages.scrollHeight);
};
socket.onclose = event => event.wasClean
    ? console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`)
    : showAlert('[close] Connection died');
socket.onerror = err => showAlert(`[error] ${err.message}`);


let messages = document.getElementById("messages");
const message = document.getElementById("message");
const form = document.getElementById("form");

form.addEventListener("submit", ev => {
    ev.preventDefault();
    return fetch("/post-message", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(message.value),
    }).then((resp) => {
        if (!resp.ok) {
            return resp.text().then(text => {
                throw new Error(text);
            })    
        }
        message.value = '';
    })
    .catch(err => showAlert(err));
});

reloadMessages = () => fetch("/get-messages")
    .then(response => response.json())
    .then(json => {
        const newMessages = document.createElement("UL");
        newMessages.setAttribute("class", "list-group");
        newMessages.setAttribute("id", "messages");

        json.forEach(message => newMessages.prepend(createMessage(message)));

        messages.parentNode.replaceChild(newMessages, messages);
        messages = document.getElementById("messages");
    })
    .catch(err => showAlert(err));

document.getElementById("app").onload = () => {
    reloadMessages().then(() => 
   		messages.parentNode.scrollTo(0, messages.scrollHeight));
    setInterval(reloadMessages, 5000);
};

