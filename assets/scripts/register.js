import { fetchResponse, displayMessages } from "./tools.js"


document.addEventListener('DOMContentLoaded', () => {
    const registerBtn = document.getElementById('registerBtn')

    registerBtn.addEventListener('click', async (event) => {
        event.preventDefault()

        const newUser = {
            username: document.getElementById('username').value,
            email: document.getElementById('email').value,
            password: document.getElementById('password').value,
            rpassword: document.getElementById('rPassword').value,
        }

        try {
            const msgs = await fetchResponse(`/confirmRegister`, newUser)
            displayMessages(msgs, '/login', `${newUser.username}, You're Registred Successfully!`)
        } catch (error) {
            console.error('Error fetching response:', error)
        }
    })
})

