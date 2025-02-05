import { fetchResponse, R_L_Popup, displayErr } from "./tools.js"

document.addEventListener('DOMContentLoaded', () => {
    const loginBtn = document.getElementById('loginBtn')
    const passwordInput = document.getElementById("password")
    const passShow = document.getElementById('passShow')

    passShow.addEventListener('click', () => {
        if (passwordInput.type === "password") {
            passwordInput.type = "text"
            passShow.src = '/assets/imgs/visible.png'
        } else {
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
                displayErr(response.body.messages)
            } else if (response.status === 200) {
                console.log("Login successful")
                R_L_Popup("/", `Hello, ${credentials.username.charAt(0).toUpperCase() + credentials.username.slice(1)}!`)
            } else {
                console.log("Unexpected response:", response.body)
            }

        } catch (error) {
            console.error('Error during login process:', error)
        }
    })
})

