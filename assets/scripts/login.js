document.addEventListener('DOMContentLoaded', () => {
    const loginButton = document.getElementById('loginButton')

    loginButton.addEventListener('click', async (event) => {
        event.preventDefault()

        const credentials = {
            username: document.getElementById('username').value,
            password: document.getElementById('password').value,
            rememberMe: document.getElementById('remember').checked,
        }
        console.log(credentials.rememberMe)

        try {
            const msgs = await fetchResponse('http://localhost:8080/confirmLogin', credentials)
            console.log(msgs)

            displayMessages(msgs, credentials.username)
        } catch (error) {
            console.error('Error fetching response:', error)
        }
    })
})

const fetchResponse = async (url, obj) => {
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(obj),
    })

    // console.log(JSON.stringify(obj))
    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)
    return response.json()
}

const displayMessages = (messages, username) => {
    const errorMsgsDiv = document.getElementById('errorMsgs')
    errorMsgsDiv.innerHTML = ''

    messages.forEach((msg) => {
        const paragraph = document.createElement('p')
        paragraph.textContent = msg
        paragraph.style.color = 'red'
        paragraph.style.fontWeight = 600
        paragraph.style.fontSize = '16px'

        errorMsgsDiv.appendChild(paragraph)
        if (msg == ['Login successful!']) {
            paragraph.style.color = 'green'
            const overlay = document.getElementById('overlay')
            const closeBtn = document.getElementById('gobtn')
            const h2Element = document.querySelector("#popup h2");
            h2Element.textContent = `Hello, ${username}!`;

            overlay.classList.add('show')

            closeBtn.addEventListener('click', (/*event */)=> {
                // event.preventDefault() //add event target
                overlay.classList.remove('show')
            })
            // window.location.href = 'http://localhost:8080/'

        }
    })
}
