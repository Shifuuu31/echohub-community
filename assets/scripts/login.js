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
            const response = await fetchResponse(`/confirmLogin`, credentials)

            if (response.status === 401) {
                console.log("Unauthorized: Invalid credentials.")

            } else if (response.status === 200) {
                console.log("Login successful" )
            } else {
                console.log("Unexpected response:", response.body)
            }
            displayMessages(response.body.messages, "/",  `Hello, ${credentials.username.charAt(0).toUpperCase()+ credentials.username.slice(1) }!`)
    
        } catch (error) {
            console.error('Error during login process:', error)
        }
    })
})

