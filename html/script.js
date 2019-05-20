const url = 'https://fnpqoxz7ll.execute-api.us-east-1.amazonaws.com/main/asdf'
document.querySelector('#button-submit-url').addEventListener('click', () => {
    const submission = document.querySelector('#submission-url').value
    fetch(url, {
        method: 'POST',
        mode: 'no-cors',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ uri: submission }),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(response.statusText);
            }
            response.json();
        })
        .then(json => completed(json, true))
        .catch(error => completed(error, false));
});

function completed(json, succeeded) {
    let alertColor = succeeded ? 'alert-success' : 'alert-danger';
    let alert = document.querySelector('#alert-result');
    alert.classList.add(alertColor);
    alert.innerHTML = json;
    alert.hidden = false;
}