const url = 'http://shortcode-dev-1630588222.us-east-1.elb.amazonaws.com'
document.querySelector('#button-submit-url').addEventListener('click', () => {
    const submission = document.querySelector('#submission-url').value
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ uri: submission }),
    })
        .then(response => response.json())
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