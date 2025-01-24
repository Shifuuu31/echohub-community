const GetMaxID = async () => {
    try {
        const response = await fetch("http://localhost:7788/max-id");
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        const maxId = await response.json();
        return maxId;
    } catch (err) {
        console.error("Error fetching maxId:", err);
    }
};

const fetchData = async (url, obj) => {
    console.log(obj.start);
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(obj),
    })

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)
    return response.json()
}

const displayPosts = async (category = "All", changeCategory = false) => {
    const url = `http://localhost:7788/post`;
    const maxId = await GetMaxID();
    const wrapper = document.querySelector(".wraper");

    if (maxId == null) {
        console.error("Failed to get maxId");
        return;
    }
    let credentials = {
        start: maxId,
        category: category
    };

    try {
        if (changeCategory) {
            const posts = document.querySelectorAll("div[id=post]");
            posts.forEach(post => wrapper.removeChild(post));
        }
        for (let count = 1; count <= 10; count++) {
            const post = await fetchData(url, credentials);
            if (post) {
                console.log(post);
                const Post = document.createElement("div");
                Post.setAttribute("id", "post");
                Post.innerHTML = `
			<div class="post-info-1"><img src="/assets/imgs/avatar.png" alt="User Avatar" loading="lazy">
            <h3>${post.PostUserName}<br><span>${post.PostTime}</span></h3>
			</div>
			<div class="post-info-2">
            <h3>${post.PostTitle}</h3>
            <p>${post.PostContent}</p>
			</div>
			<div class="post-info-3">
            <div class="links">
            ${post.PostCategories.map(category => `<li>${category}</li>`)}
            </div>
            <div class="buttons">
            <button><img src="/assets/imgs/like.png" alt="Like"> ${post.LikeCount}</button>
            <button><img src="/assets/imgs/dislike.png" alt="Dislike"> ${post.DislikeCount}</button>
            <button><img src="/assets/imgs/comment.png" alt="Comment"> ${post.CommentsCount}</button>
            </div>
			</div>`

                wrapper.appendChild(Post);
                credentials.start = post.PostId - 1
            } else {
                console.log("No posts to display.");
            }
        }
    } catch (e) {
        console.error("Error fetching data:", e.message);
    }
};

const categories = document.querySelectorAll("input[id^=category]");
categories.forEach(category => {
    category.addEventListener('change', () => {
        const categoryChecked = document.querySelectorAll('input[id^=category]:checked');
        displayPosts(categoryChecked[0].value, true)
    });
});

displayPosts();

