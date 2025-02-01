export { fetchResponse, displayMessages, CheckClick, CheckLength }

const fetchResponse = async (url, obj = {}) => {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(obj),
        })
        const responseBody = await response.json()
        
        return { status: response.status, body: responseBody }
    } catch (error) {
        console.error('Error fetching response:', error)
        throw error
    }
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


const CheckClick = () => {
    const categoriesChecked = document.querySelectorAll('input[id^=category]:checked')
    if (categoriesChecked.length == 0) {
        alert(`Select at least one category`)
        return false;
    }
    return true;
}

const MAX_CATEGORIES = 3
const CheckLength = (category, checkedLength) => {
    if (checkedLength > MAX_CATEGORIES) {
        category.checked = false
        alert(`You can only select up to ${MAX_CATEGORIES} categories.`)
    }
}

