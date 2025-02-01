export { desplayPosts, DataToFetch }
import { fetchResponse } from "./tools.js"
let DataToFetch = {}

const desplayPosts = async (category = "All", scroll = false) => {
    const posts = document.getElementById("posts")

    if (!scroll) {
        try {
            const response = await fetchResponse(`/maxId`)

            if (response.status === 200) {
                if (response.body == 0) {
                    posts.innerHTML = `<h1 style="text-align: center">No posts to display1</h1>`
                    return false
                }
                DataToFetch.postID = response.body
                DataToFetch.category = category
                posts.innerHTML = ''
            } else {
                console.log("Unexpected response:", response.body)
            }

        } catch (error) {
            console.error('Error during fetching maxId:', error)
        }
    }
    let FetchedPosts = []

    try {
        const response = await fetchResponse(`/posts`, DataToFetch)
        if (response.status === 200) {
            if (response.body.type == 'client') {
                posts.innerHTML = `<h1 style="text-align: center">No posts to display2</h1>`
                return false
            }
            console.log("Posts Fetched succefully")
            FetchedPosts = response.body
        }else if (response.status === 400) {
            console.log("Bad Request", response.status, response.body.message)
        } else {
            console.log("Unexpected response:", response.body)
        }

    } catch (error) {
        console.error('Error during fetching Posts:', error)
    }

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
            DataToFetch.postID = FetchedPosts[FetchedPosts.length - 1].PostId - 1
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