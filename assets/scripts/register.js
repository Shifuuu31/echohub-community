import { fetchResponse, R_L_Popup } from "./tools.js"

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
            gender: document.getElementById('gender').value,
            password: passwordInput.value,
            rpassword: rPasswordInput.value,
        }
        console.log(newUser)

        try {
            const response = await fetchResponse(`/confirmRegister`, newUser)
            console.log(response)
            if (response.status === 400) {
                console.log("Bad request: Invalid info Or Missing field.")
                displayMsg(response.body.messages)

            } else if (response.status === 200) {
                console.log("Registred successfully" )
                displayMsg(response.body.messages, true)
                R_L_Popup('/login', `${newUser.username.charAt(0).toUpperCase()+ newUser.username.slice(1)}, You're Registred Successfully!`)
            } else {
                console.log("Unexpected response:", response.body)
            }
    
        } catch (error) {
            console.error('Error during login process:', error)
        }

    })
})

