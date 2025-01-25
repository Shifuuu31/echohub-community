// const Credentials = {
//     username: 'yassine',
//     password: 'rF8$FjBHz@m48&b',
//     rememberMe: false
// }
const newUser = {
    username: "ecde",
    email: "ecd.e@ed.de",
    password: "e2dededeDde.e",
    rpassword: "e2dededeDde.e",
}

const fetchResponse = (url, obj) => {
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(obj)
    })
        .then(response => response.json())
        .then(data => {
            console.log('Response Data:', data)
        })
        .catch(error => console.error('Error:', error))
}

fetchResponse('http://localhost:8080/confirmRegister', newUser)
// fetchResponse('http://localhost:8080/confirmLogin', Credentials)