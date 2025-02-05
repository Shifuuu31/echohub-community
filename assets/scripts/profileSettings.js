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

        console.log(toUpdate)

        try {
            const response = await fetchResponse(`/updateProfile`, toUpdate)
            console.log(response)
            if (response.status === 400) {
                console.log("Bad request: Invalid info Or Missing field.")
                displayMsg(response.body.messages)
            } else if (response.status === 200) {
                console.log("Profile Updated successfully" )
                displayMsg(response.body.messages, true)
            } else {
                console.log("Unexpected response:", response.body)
            }
    
        } catch (error) {
            console.error('Error during login process:', error)
        }
    })
})

