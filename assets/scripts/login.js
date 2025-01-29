import { fetchResponse, displayMessages } from "./tools.js"

document.addEventListener('DOMContentLoaded', () => {
    const loginBtn = document.getElementById('loginBtn')
    const passwordInput = document.getElementById("password")
    const passShow = document.getElementById('passShow')

    passShow.addEventListener('click', () => {
        if (passwordInput.type === "password") {
            passwordInput.type = "text"
            passShow.src = '/assets/imgs/visible.png'
        }else {
            passwordInput.type = "password"
            passShow.src = '/assets/imgs/unvisible.png'
        }
    })
    
    loginBtn.addEventListener('click', async (event) => {
        event.preventDefault()

        const credentials = {
            username: document.getElementById('username').value,
            password: passwordInput.value,
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

