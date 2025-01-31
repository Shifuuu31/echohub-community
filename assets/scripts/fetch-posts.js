export { desplayPosts }
import { fetchResponse } from "./tools.js"

const desplayPosts = async (category = "All", append = false) => {
    const posts = document.getElementById("posts")
    let DataToFetch = {}

    if (!append) {
        try {
            const response = await fetchResponse(`/confirmLogin`, credentials)

            if (response.status === 401) {
                console.log("Unauthorized: Invalid credentials.")

            } else if (response.status === 200) {
                console.log("Login successful" )
            } else {
                console.log("Unexpected response:", response.body)
            }
            displayMessages(response.body.messages, "/",  `Hello, ${credentials.username.charAt(0).toUpperCase()+ credentials.username.slice(1) }!`)
    
        } catch (error) {
            console.error('Error during login process:', error)
        }


        const maxId = await fetchResponse("/maxId")
        console.log(maxId);
        
        if (maxId.status != 200) {
            console.error("Failed to get maxId")
            return false
        }
        DataToFetch.postID = maxId.body
        DataToFetch.category = category
        posts.innerHTML = ''
    }

    if (DataToFetch.postID) {
        const FetchedPosts = await fetchResponse("/post", DataToFetch)
        if (FetchedPosts) {
            for (let i = 0; i < FetchedPosts.length; i++) {
                const postData = document.createElement('div')
                postData.innerHTML = `          
                <div id="post" post-id="${FetchedPosts[i].PostId}">
                    <div id="user-post-info"><img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
                        <h3>@${FetchedPosts[i].PostUserName} <br><span>${new Date(FetchedPosts[i].PostTime).toUTCString()}</span></h3>
                        <div id="dropdown-content" style="margin-left:auto">
                            <a href="/updatePost?ID=${FetchedPosts[i].PostId}"><img src="/assets/imgs/update.png" style="border:none; border-radius:0px;"> Update Post</a>
                            <hr>
                            <a href="/deletePost?ID=${FetchedPosts[i].PostId}"><img src="/assets/imgs/delete.png" style="border:none; border-radius:0px;"> Delete Post</a>
                        </div>
                    </div>
                    <div id="post-body">
                        <h3 id="post-title">${FetchedPosts[i].PostTitle}</h3>
                        <pre>${wrapLinks(FetchedPosts[i].PostContent)}</pre>
                    </div>
                    <div id="post-categories">
                        <div id="links">
                            ${(FetchedPosts[i].PostCategories || []).map(category => `<li>${category}</li>`).join(' ')}
                        </div>
                        <div id="buttons" >
                            <button><img src="/assets/imgs/like.png" alt="Like"> ${FetchedPosts[i].LikeCount}</button>
                            <button><img src="/assets/imgs/dislike.png" alt="Dislike"> ${FetchedPosts[i].DislikeCount}</button>
                            <button id="commentBtn"><img src="/assets/imgs/comment.png" alt="Comment"> ${FetchedPosts[i].CommentsCount}</button>
                        </div>
                    </div>`
                posts.append(postData)
            }
            DataToFetch.postID = FetchedPosts[FetchedPosts.length - 1].PostId
        } else if (!append) {
            posts.innerHTML = `<h1 style="text-align: center">No posts to display.</h1>`
            return false
        }
    }

    return true
}

function wrapLinks(text) {
    const urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig

    const wrappedText = text.replace(urlRegex, (url) => {
        return `<a href='${url}' target="_blank">${url}</a>`
    })

    return wrappedText
}