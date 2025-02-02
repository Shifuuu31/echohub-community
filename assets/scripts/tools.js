export { fetchResponse, displayMessages, CheckClick, CheckLength, addPost }

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

const displayMessages = (messages, redirectUrl, popupMsg) => {
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
            goBtn.addEventListener('click', () => {
                overlay.classList.remove('show')
                window.location.href = redirectUrl
            })
        }
    })
}

// const sleep = ms => new Promise(r => setTimeout(r, ms));

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

// add post div to html
const addPost = (post) => {
    const postData = document.createElement('div')
    postData.innerHTML = `          
                <div id="post" post-id="${post.PostId}">
                    <div id="user-post-info"><img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
                        <h3>@${post.PostUserName} <br><span>${new Date(post.PostTime).toUTCString()}</span></h3>
                        <div id="dropdown-content" style="margin-left:auto">
                            <a href="/updatePost?ID=${post.PostId}"><img src="/assets/imgs/update.png" style="border:none; border-radius:0px;"> Update Post</a>
                            <hr>
                            <a href="/deletePost?ID=${post.PostId}"><img src="/assets/imgs/delete.png" style="border:none; border-radius:0px;"> Delete Post</a>
                        </div>
                    </div>
                    <div id="post-body">
                        <h3 id="post-title">${post.PostTitle}</h3>
                        <pre>${wrapLinks(post.PostContent)}</pre>
                    </div>
                    <div id="post-categories">
                        <div id="links">
                            ${(post.PostCategories || []).map(category => `<li>${category}</li>`).join(' ')}
                        </div>
                        <div id="buttons" >
                            <button><img src="/assets/imgs/like.png" alt="Like"> ${post.LikeCount}</button>
                            <button><img src="/assets/imgs/dislike.png" alt="Dislike"> ${post.DislikeCount}</button>
                            <button id="commentBtn"><img src="/assets/imgs/comment.png" alt="Comment"> ${post.CommentsCount}</button>
                        </div>
                    </div>`

    return postData
}

function wrapLinks(text) {
    const urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig

    const wrappedText = text.replace(urlRegex, (url) => {
        return `<a href='${url}' target="_blank">${url}</a>`
    })

    return wrappedText
}

// add Skeleton
