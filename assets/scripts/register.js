document.addEventListener('DOMContentLoaded', () => {
    const loginButton = document.getElementById('registerBtn')

    loginButton.addEventListener('click', async (event) => {
        event.preventDefault()

        const newUser = {
            username: document.getElementById('username').value,
            email: document.getElementById('email').value,
            password: document.getElementById('password').value,
            rpassword: document.getElementById('rPassword').value,
        }
        console.log(newUser.username)
        console.log(newUser.email)
        console.log(newUser.password)
        console.log(newUser.rpassword)

        try {
            const msgs = await fetchResponse(`http://${window.location.href}/confirmRegister`, newUser)
            console.log(msgs)

            displayMessages(msgs)
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
const sleep = ms => new Promise(r => setTimeout(r, ms));

const displayMessages = (messages) => {
    const errorMsgsDiv = document.getElementById('errorMsgs')
    errorMsgsDiv.innerHTML = ''

    messages.forEach((msg) => {
        const paragraph = document.createElement('p')
        paragraph.textContent = msg
        paragraph.style.color = 'red'
        paragraph.style.fontWeight = 600
        paragraph.style.fontSize = '16px'

        errorMsgsDiv.appendChild(paragraph)
        if (msg == ['User Registred successfully!']) {
            paragraph.style.color = 'green'
            sleep(1000)
            window.location.href = '/login'

        }
    })
}
