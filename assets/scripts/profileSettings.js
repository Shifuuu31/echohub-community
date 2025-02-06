import { DropDown, fetchResponse, displayMsg } from "./tools.js"

document.addEventListener("DOMContentLoaded", () => {
    const usernameInput = document.getElementById('username')
    const emailInput = document.getElementById('email')
    const passwordInput = document.getElementById("password")
    const rPasswordInput = document.getElementById("rPassword")
    DropDown()
    const updateBtn = document.getElementById('update')

    updateBtn.addEventListener('click', async() => {
        const toUpdate = {
            username: usernameInput.value,
            email: emailInput.value,
            password: passwordInput.value,
            rpassword: rPasswordInput.value,
            changes: [],
        }
        if (toUpdate.username.length != 0) toUpdate.changes.push('username')
        if (toUpdate.email.length != 0) toUpdate.changes.push('email')
        if (toUpdate.password.length != 0) toUpdate.changes.push('password')

        try {
            const response = await fetchResponse(`/updateProfile`, toUpdate)
            if (response.status === 400) {
                console.log("Bad request: Invalid info Or Missing field.")
                displayMsg(response.body.messages)
            } else if (response.status === 200) {
                console.log("Profile Updated successfully" )
                console.log(response.body)
                if (response.body.extra.includes('username')) {
                    document.getElementById('h-username').firstChild.nodeValue = usernameInput.placeholder = toUpdate.username
                    usernameInput.value = ''
                }
                if (response.body.extra.includes('email')) {
                    document.getElementById('h-email').innerText = emailInput.placeholder = toUpdate.email
                    emailInput.value = ''
                }
                if (response.body.extra.includes('password')) {
                    passwordInput.placeholder = 'enter new password'
                    rPasswordInput.placeholder = 'repeat your new password'
                    rPasswordInput.value = rPasswordInput.value  = ''
                }
                displayMsg(response.body.messages, true)
            } else {
                console.log("Unexpected response:", response.body)
            }
    
        } catch (error) {
            console.error('Error during login process:', error)
        }
    })
})

