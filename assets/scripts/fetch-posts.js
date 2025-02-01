export { displayPosts, DataToFetch }
import { fetchResponse, addPost } from "./tools.js"
let DataToFetch = {}

const displayPosts = async (category = "All", scroll = false) => {
    const postsContainer = document.getElementById("posts")
    const postMsg = document.getElementById("postMsg")

    if (!scroll) {
        // get max id
        try {
            const response = await fetchResponse(`/maxId`)
            if (response.status === 200) {
                if (response.body == 0) {
                    console.log('No posts to display')
                    postMsg.innerHTML = `<h1 style="text-align: center">No posts to display</h1>`
                    return false
                }
                DataToFetch.start = response.body
                DataToFetch.category = category
                postsContainer.innerHTML = ''
                postMsg.innerHTML = ''
            } else {
                console.log("Unexpected response:", response.body)
                return false
            }
        } catch (error) {
            console.error('Error during fetching maxId:', error)
        }
    }

    let FetchedPosts = []
    // get posts
    try {
        const response = await fetchResponse(`/posts`, DataToFetch)
        if (response.status === 200) {
            console.log("Posts Fetched succefully")
            FetchedPosts = response.body
        } else if (response.status === 100) {
            console.log('No posts to display')
            postMsg.innerHTML = `<h1 style="text-align: center">No posts to display</h1>`
            return false
        } else if (response.status === 400) {
            console.log("Bad Request", response.status, response.body.message)
            return false
        } else {
            console.log("Unexpected response:", response.body)
            return false
        }
    } catch (error) {
        console.error('Error during fetching Posts:', error)
    }

    // check if there is posts
    if (FetchedPosts.length > 0) {
        for (let i = 0; i < FetchedPosts.length; i++) {
            let postData = addPost(FetchedPosts[i])
            postsContainer.append(postData)
        }
        if (FetchedPosts.length < 10) {
            console.log('No more posts to display')
            postMsg.innerHTML = `<h1 style="text-align: center">No more posts to display</h1>`
            return false
        } else {
            // send last post fetched id for scroll
            DataToFetch.start = FetchedPosts[FetchedPosts.length - 1].PostId - 1
        }
    }  else {
        postMsg.innerHTML = `<h1 style="text-align: center">${scroll ? 'No more posts to display' : 'No posts to display'}</h1>`;
        return false;
    }

    return true
}

