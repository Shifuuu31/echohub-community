// Get max ID 
const GetMaxID = async () => {
    try {
        const response = await fetch(`http://${window.location.host}/max-id`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        const maxId = await response.json();
        return maxId;
    } catch (err) {
        console.error("Error fetching maxId:", err);
    }
};

// function to fetch data
const fetchData = async (url, obj) => {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(obj),
        })
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)
        return response.json()

    } catch (e) {
        console.error("Error fetching data:", e.message);
    }
}

const nothing = document.createElement("div");
nothing.style.textAlign = 'center'
const posts = document.getElementById("posts");

const displayPosts = async (category = "All") => {
    const url = `http://${window.location.host}/post`;
    const maxId = await GetMaxID();
    let countPosts = 0

    if (maxId == null) {
        console.error("Failed to get maxId");
        return;
    }
    let DataToFetch = {
        postID: maxId,
        category: category
    };
    posts.innerHTML = ''
    nothing.innerHTML = ''
    for (let count = 1; count <= 10; count++) {
        const postData = document.createElement('div')
        const post = await fetchData(url, DataToFetch);
        if (post) {
            countPosts++
            postData.innerHTML = `
            <div id="post">
            <div class="post-info-1"><img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
            <h3>${post.PostUserName}<br><span>${new Date(post.PostTime).toUTCString()}</span></h3>
            </div>
            <div class="post-info-2">
            <h3>${post.PostTitle}</h3>
            <p>${post.PostContent}</p>
            </div>
            <div class="post-info-3">
            <div class="links">
            ${(post.PostCategories).map(category => `<li>${category}</li>`).join(' ')}
            </div>
            <div class="buttons">
            <button><img src="/assets/imgs/like.png" alt="Like"> ${post.LikeCount}</button>
            <button><img src="/assets/imgs/dislike.png" alt="Dislike"> ${post.DislikeCount}</button>
            <button id="commentBtn"><img src="/assets/imgs/comment.png" alt="Comment"> ${post.CommentsCount}</button>
            </div>
            </div>
            </div>`
            posts.append(postData)
            DataToFetch.postID = post.PostId - 1
            if (DataToFetch.postID == 0) {
                break
            }
        }
    }
    if (countPosts === 0) {
        posts.innerHTML = `<h1> No posts to display.</h1>`
    }
};

const categories = document.querySelectorAll("input[id^=category]");
categories.forEach(category => {
    category.addEventListener('change', (target) => {
        displayPosts(target.target.defaultValue)
    });
});

displayPosts();

