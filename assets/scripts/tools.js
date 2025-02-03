export { fetchResponse, displayMessages, DropDown, AddPost, AddComment }

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
        // console.error('Error fetching response:', error)
        throw error
    }
}
const DropDown = () => {
    const profilePictureContainer = document.getElementById('profile-picture-container')
    const dropdown = document.getElementById('dropdown')
    if (profilePictureContainer && dropdown) {

        const toggleDropdown = (event) => {
            event.stopPropagation()
            dropdown.classList.toggle('active')
        }

        profilePictureContainer.addEventListener('click', toggleDropdown)
        document.addEventListener('click', () => { dropdown.classList.remove('active') })
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

// add post div to html
const AddPost = (post) => {
    const postData = document.createElement('div')
    postData.innerHTML =
        `<div id="post" post-id="${post.PostId}">
    <div id="user-post-info">
        <section style="display: flex;">
            <img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
            <h3>@${post.PostUserName} <br><span>${timeAgo(post.PostTime)}</span></h3>
        </section>
        ${Username === post.PostUserName ? '<button><img src="/assets/imgs/option.png"></button>' : ''}
    </div>
    ${Username === post.PostUserName ?
            `<div id="post-dropdown">
        <div id="dropdown-content">
            <a href="/updatePost?ID=${post.PostId}"><img src="/assets/imgs/update.png"> Update Post</a>
            <hr>
            <a href="/deletePost?ID=${post.PostId}"><img src="/assets/imgs/delete.png"> Delete Post</a>
        </div>
    </div>` : ''}
    <div id="post-body">
        <h3 id="post-title">${post.PostTitle}</h3>
        <pre>${wrapLinks(post.PostContent)}</pre>
    </div>
    <div id="post-categories">
        <div id="links">
            ${(post.PostCategories || []).map(category => `<li>${category}</li>`).join(' ')}
        </div>
        <div id="buttons">
            <button><img src="/assets/imgs/like.png" alt="Like"> ${post.LikeCount}</button>
            <button><img src="/assets/imgs/dislike.png" alt="Dislike"> ${post.DislikeCount}</button>
            <button id="commentBtn"><img src="/assets/imgs/comment.png" alt="Comment"> ${post.CommentsCount}</button>
        </div>
    </div>
</div>`

    return postData
}

const AddComment = (comment) => {
    const commentDiv = document.createElement("div");
    commentDiv.id = "comment";
    commentDiv.innerHTML = `
        <div id="user-info-and-buttons">
            <div id="user-comment-info">
                <img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
                <h3>${comment.UserName} <br><span>${timeAgo(comment.CreationDate)}</span></h3>
            </div>
        </div>
        <div id="user-comment-info">
            <p>${comment.Content}</p>
        </div>`
    return commentDiv
}

function wrapLinks(text) {
    const urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig

    const wrappedText = text.replace(urlRegex, (url) => {
        return `<a href='${url}' target="_blank">${url}</a>`
    })

    return wrappedText
}

function timeAgo(input) {
    const date = input instanceof Date ? input : new Date(input);
    const formatter = new Intl.RelativeTimeFormat('en');
    const seconds = (Date.now() - date) / 1000;

    const units = [
        ['year', 31536000],
        ['month', 2592000],
        ['week', 604800],
        ['day', 86400],
        ['hour', 3600],
        ['minute', 60],
        ['second', 1]
    ];

    for (const [unit, secondsInUnit] of units) {
        if (Math.abs(seconds) >= secondsInUnit) {
            return formatter.format(-Math.round(seconds / secondsInUnit), unit);
        }
    }

    return 'just now';
}





