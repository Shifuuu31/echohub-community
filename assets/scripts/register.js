import { fetchResponse, displayMessages } from "./tools.js"


document.addEventListener('DOMContentLoaded', () => {
    const registerBtn = document.getElementById('registerBtn')
    const passwordInput = document.getElementById("password")
    const rPasswordInput = document.getElementById("rPassword")
    const passShow = document.getElementById('passShow')
    const rPassShow = document.getElementById('rPassShow')

    passShow.addEventListener('click', () => {
        if (passwordInput.type === "password") {
            passwordInput.type = "text"
            passShow.src = '/assets/imgs/visible.png'
        }else {
            passwordInput.type = "password"
            passShow.src = '/assets/imgs/unvisible.png'
        }
    })
    rPassShow.addEventListener('click', () => {
        if (rPasswordInput.type === "password") {
            rPasswordInput.type = "text"
            rPassShow.src = '/assets/imgs/visible.png'
        }else {
            rPasswordInput.type = "password"
            rPassShow.src = '/assets/imgs/unvisible.png'
        }
    })

    registerBtn.addEventListener('click', async (event) => {
        event.preventDefault()

        const newUser = {
            username: document.getElementById('username').value,
            email: document.getElementById('email').value,
            password: passwordInput.value,
            rpassword: rPasswordInput.value,
        }

        try {
            const msgs = await fetchResponse(`/confirmRegister`, newUser)
            displayMessages(msgs, '/login', `${newUser.username}, You're Registred Successfully!`)
        } catch (error) {
            console.error('Error fetching response:', error)
        }
    })
})

