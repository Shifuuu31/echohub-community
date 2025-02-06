import {openPopup} from './popup.js'

export { fetchResponse, R_L_Popup, DropDown, AddPost, AddComment, displayMsg, AddOrUpdatePost, CategoriesSelection }

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
        console.error(error)
    }
}

const DropDown = () => {
    const profilePictureContainer = document.getElementById('profile-picture-container')
    const dropdown = document.getElementById('dropdown')
    if (profilePictureContainer && dropdown) {

        const toggleDropdown = (event) => {
            console.log("pD")
            event.stopPropagation()
            dropdown.classList.toggle('active')
        }

        profilePictureContainer.addEventListener('click', toggleDropdown)
        document.addEventListener('click', () => { dropdown.classList.remove('active') })
    }
}

const R_L_Popup = (redirectUrl, popupMsg) => {
    const overlay = document.getElementById('overlay')
    const goBtn = document.getElementById('gobtn')
    const h2Element = document.querySelector("#popup h2")
    h2Element.textContent = popupMsg
    overlay.classList.add('show')
    goBtn.addEventListener('click', () => {
        overlay.classList.remove('show')
        window.location.href = redirectUrl
    })
}

// add post div to html
const AddPost = (post) => {
    post.Content = wrapLinks(post.Content)
    let splittedContent = post.Content
    let moreContent = ''
    const flag = (post.Content.length > 180)
    if (flag) {
        splittedContent = post.Content.slice(0, 150)
        moreContent = post.Content.slice(150)
    }

    const postData = document.createElement('div')
    postData.innerHTML =
        `<div id="post" post-id="${post.ID}">
            <div id="user-post-info">
                <section style="display: flex">
                    <img src="${post.ProfileImg}" alt="User Avatar" loading="lazy">
                    <h3>@${post.UserName} <br><span>${timeAgo(post.CreatedAt)}</span></h3>
                </section>
                ${Username === post.UserName ? '<button class="options-btn"><img src="/assets/imgs/option.png"></button>' : ''}
            </div>
            ${Username === post.UserName ?
            `<div id="post-dropdown" class="post-dropdown${post.ID}">
                <div id="dropdown-content">
                    <a href="/updatePost?ID=${post.ID}"><img src="/assets/imgs/update.png"> Update Post</a>
                    <hr>
                    <a href="/deletePost?ID=${post.ID}"><img src="/assets/imgs/delete.png"> Delete Post</a>
                </div>
            </div>` : ''}
            <div id="post-body">
                <h3 id="post-title">${post.Title}</h3>
                <pre>${splittedContent}<span id="moreContent">${moreContent}</span><span ${flag ? `id="dots">...</span>` : ''}</pre>
                ${flag ? `<button id="moreBtn">more</button>` : ''}
            </div>
            <div id="post-categories">
                <div id="buttons">
                    <button><img src="/assets/imgs/like.png" alt="Like"> ${post.LikeCount}</button>
                    <button><img src="/assets/imgs/dislike.png" alt="Dislike"> ${post.DislikeCount}</button>
                    <button id="commentBtn"><img src="/assets/imgs/comment.png" alt="Comment"> ${post.CommentsCount}</button>
                </div>
                <div id="links">
                    ${(post.Categories || []).map(category => `<li>${category}</li>`).join(' ')}
                </div>
            </div>
        </div>`

        
        if (Username === post.UserName) {
            const optionsButton = postData.querySelector('.options-btn')
        const dropdown = postData.querySelector(`.post-dropdown${post.ID}`)

        optionsButton.addEventListener('click', (event) => {
            event.stopPropagation();
            dropdown.classList.toggle('active')
        })

        document.addEventListener('click', (event) => {
            if (!dropdown.contains(event.target) && !optionsButton.contains(event.target)) {
                dropdown.classList.remove('active')
            }
        })
    }


    
    const cmntBtn = postData.querySelector("#commentBtn")
    console.log('AH',cmntBtn);
    if (cmntBtn) {
        // cmntBtn.removeEventListener("click", openPopup)
        cmntBtn.addEventListener("click", openPopup)

    } 

    const btn = postData.querySelector("#moreBtn")
    const dots = postData.querySelector("#dots")
    
    const moreText = postData.querySelector("#moreContent")
    if (btn && dots) {
        btn.addEventListener("click", () => {
            if (dots.style.display === "none") {
                dots.style.display = "inline"
                btn.innerHTML = "more"
                moreText.style.display = "none"
            } else {
                dots.style.display = "none"
                btn.innerHTML = "less"
                moreText.style.display = "inline"
            }
        })
    }


    
    return postData
}


const AddComment = (comment) => {
    const commentDiv = document.createElement("div")
    commentDiv.id = "comment"
    commentDiv.innerHTML = `
    <div id="user-info-and-buttons">
    <div id="user-comment-info">
                <img src="${comment.ProfileImg}" alt="User Avatar" loading="lazy">
                <h3>@${comment.UserName} <br><span>${timeAgo(comment.CreationDate)}</span></h3>
            </div>
            <div id="buttons">
                <button><img src="/assets/imgs/like.png" alt="Like"> 12</button>
                <button><img src="/assets/imgs/dislike.png" alt="Dislike"> 4123</button>
            </div>
        </div>
        <div id="user-comment-info">
            <p>${comment.Content}</p>
        </div>`
    return commentDiv
}



const AddOrUpdatePost = async (url) => {
    console.log(url)
    const newPost = {
        title: document.getElementById('title').value,
        content: document.getElementById('content').value,
        selectedCategories: []
    }
    document.querySelectorAll('input[id^=category]:checked').forEach((selectedCategory) => {
        newPost.selectedCategories.push(selectedCategory.value)
    })

    // console.log(newPost)

    const response = await fetchResponse(url, newPost)
    console.log(response.body);
    if (response.status === 401) {
        console.log("Unauthorized: try to login")
    } else if (response.status === 400) {
        console.log(response.body);
        displayMsg(response.body.messages)
    } else if (response.status === 200) {
        console.log("post added successfully")
        window.location.href = "/"
    } else {
        console.log("Unexpected response:", response.body)
    }

}


const CategoriesSelection = () => {
    let categories = document.querySelectorAll('input[id^=category]')

    categories.forEach((category) => {
        category.addEventListener('change', () => {
            const checkedCategories = document.querySelectorAll('input[id^=category]:checked')
            if (checkedCategories.length > 3) {
                category.checked = false
                displayMsg(['You can only select up to 3 categories'])
            }
        })
    })
}

const displayMsg = (msgs, success = false) => {
    const errPopups = document.querySelectorAll('.errPopup')
    errPopups.forEach(popup => popup.remove())

    const baseTop = 100
    const gap = 50

    msgs.forEach((msg, index) => {
        const errPopup = document.createElement("div")
        errPopup.id = `errPopup-${index}`
        errPopup.classList.add("errPopup")

        if (success == true) errPopup.style.backgroundColor = '#02bf08' // lightgreen

        errPopup.innerHTML = `
        <span class="close-btn">&times;</span>
        ${msg}
      `

        errPopup.style.top = `${baseTop + index * gap}px`

        document.body.appendChild(errPopup);

        const closeBtn = errPopup.querySelector('.close-btn')
        if (closeBtn) {
            closeBtn.addEventListener('click', () => {
                errPopup.remove()
            });
        }
    });
};


function wrapLinks(text) {
    const urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.]*[-A-Z0-9+&@#\/%=~_|])/ig

    const wrappedText = text.replace(urlRegex, (url) => {
        return `<a href='${url}' target="_blank">${url}</a>`
    })

    return wrappedText
}

function timeAgo(input) {
    const date = input instanceof Date ? input : new Date(input)
    const formatter = new Intl.RelativeTimeFormat('en')
    const seconds = (Date.now() - date) / 1000

    const units = [
        ['year', 31536000],
        ['month', 2592000],
        ['week', 604800],
        ['day', 86400],
        ['hour', 3600],
        ['minute', 60],
        ['second', 1]
    ]

    for (const [unit, secondsInUnit] of units) {
        if (Math.abs(seconds) >= secondsInUnit) {
            return formatter.format(-Math.round(seconds / secondsInUnit), unit)
        }
    }

    return 'just now'
}





