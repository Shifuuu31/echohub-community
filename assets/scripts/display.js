export { displayPosts, DataToFetch }
import { fetchResponse, AddPost } from "./tools.js"
let DataToFetch = {}

const displayPosts = async (category = "All", scroll = false) => {
    const postsContainer = document.getElementById("posts")
    const postMsg = document.getElementById("postMsg")
    const Skeleton = document.getElementById("post-placeholder")

    if (!scroll) {
        // get max id
        try {
            const response = await fetchResponse(`/maxId`)
            if (response.status === 200) {
                if (response.body == 0) {
                    console.log('No posts to display')
                    postMsg.innerHTML = `<h1>No posts to display</h1>`
                    Skeleton.style.display = 'none'
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
            postMsg.innerHTML = `<h1>No posts to display</h1>`
            Skeleton.style.display = 'none'
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
            postsContainer.append(AddPost(FetchedPosts[i]))
        }
        if (FetchedPosts.length < 10) {
            console.log('No more posts to display')
            postMsg.innerHTML = `<h1>No more posts to display</h1>`
            Skeleton.style.display = 'none'
            return false
        } else {
            // send last post fetched id for scroll
            DataToFetch.start = FetchedPosts[FetchedPosts.length - 1].PostId - 1
        }
    } else {
        console.log(`${scroll ? 'No more posts to display' : 'No posts to display'}`);
        postMsg.innerHTML = `<h1>${scroll ? 'No more posts to display' : 'No posts to display'}</h1>`
        Skeleton.style.display = 'none'
        return false;
    }

    return true
}

