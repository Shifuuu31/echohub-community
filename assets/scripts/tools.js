export { fetchResponse, displayMessages}

const fetchResponse = async (url, obj) => {
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(obj),
    })

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)
    return response.json()
}


const displayMessages =  (messages, redirectUrl, popupMsg) => {
    const errorMsgsDiv = document.getElementById('errorMsgs')
    errorMsgsDiv.innerHTML = ''

    messages.forEach((msg) => {
        const paragraph = document.createElement('p')
        paragraph.textContent = msg
        paragraph.style.color = 'red'
        paragraph.style.fontWeight = 600
        paragraph.style.fontSize = '16px'

        errorMsgsDiv.appendChild(paragraph)
        
        if (msg == 'User Registred successfully!' || msg == 'Login successful!') {
            paragraph.style.color = 'green'
            const overlay = document.getElementById('overlay')
            const goBtn = document.getElementById('gobtn')
            const h2Element = document.querySelector("#popup h2");
            h2Element.textContent = popupMsg
            overlay.classList.add('show')
            goBtn.addEventListener('click', ()=> {
                overlay.classList.remove('show')
                window.location.href = redirectUrl
            })
        }
    })
}

const sleep = ms => new Promise(r => setTimeout(r, ms));
