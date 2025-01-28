import { fetchResponse, displayMessages } from "./tools.js"

document.addEventListener('DOMContentLoaded', () => {
    const loginBtn = document.getElementById('loginBtn')

    loginBtn.addEventListener('click', async (event) => {
        console.log('dsds')
        event.preventDefault()

        const credentials = {
            username: document.getElementById('username').value,
            password: document.getElementById('password').value,
            rememberMe: document.getElementById('remember').checked,
        }

        try {
            const msgs = await fetchResponse(`/confirmLogin`, credentials)
            displayMessages(msgs, '/', `Hello, ${credentials.username}!`)
        } catch (error) {
            console.error('Error fetching response:', error)
        }
    })
})

