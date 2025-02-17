import { fetchResponse, AddComment, displayMsg } from "./tools.js"
// import {setupLikeDislikeListner } from "./likes&dislikes.js"
export { openPopup, closePopup, popupBackground }

const popup = document.getElementById("popup")
const popupBackground = document.getElementById("popup-background")
const closeButton = document.querySelector(".close")
const commentsSection = document.getElementById("comments-section");


const openPopup = async (event) => {

    if (popup && popupBackground) {
        popupBackground.style.display = popup.style.display = "block"

        const targetedPost = event.target.closest("#post")
        const postID = targetedPost.getAttribute("post-id")
        const cmntGrp = document.getElementById('comment-group')
        if (cmntGrp) {
            cmntGrp.innerHTML = `<textarea placeholder="Type a comment..." type="text" id="comment-field" maxlength="1000"></textarea>
                        <button class="new-comment" id="${postID}"><i class="fas fa-paper-plane"></i></button>`

            const newCmntBtn = document.getElementById(`${postID}`)

            newCmntBtn.addEventListener('click', async () => {
                const cmntField = document.getElementById('comment-field')

                const created = await createComment({
                    postid: postID,
                    content: cmntField.value,
                })
                if (created == true) {
                    const postCmntBtn = targetedPost.querySelector('#commentBtn')
                    postCmntBtn.childNodes[1].nodeValue = parseInt(postCmntBtn.childNodes[1].nodeValue.trim(), 10) + 1
                    cmntField.value = ''
                }
            })
        }

        await displayComments(postID);
        // setupLikeDislikeListner()
    }
}

const closePopup = (event) => {
    if (event.target === popupBackground || event.target === closeButton) {
        popupBackground.style.display = popup.style.display = "none"
    }
}

const createComment = async (newCmnt) => {
    try {
        const response = await fetchResponse(`/createComment`, newCmnt)
        if (response.status === 200) {
            await displayComments(newCmnt.postid)
            return true
        } else if (response.status === 400) {
            console.log(response.body)
            displayMsg([response.body])

        } else {
            console.log("Unexpected response:", response.body)
        }

        return false
    } catch (error) {
        console.error('Error during login process:', error)
    }
    return false
}


const displayComments = async (postid) => {
    commentsSection.innerHTML = ''


    let comments
    try {
        const response = await fetchResponse(`/comments`, { ID: postid })
        if (response.status === 200) {
            console.log("comments recieved succesfully")
            comments = response.body

        } else {
            console.log("Unexpected response:", response.body)
        }

    } catch (error) {
        console.error('Error during login process:', error)
    }

    if (comments.length > 0) {
        comments.forEach(comment => {
            console.log(comment);
            commentsSection.appendChild(AddComment(comment))
        })
    } else {
        commentsSection.innerHTML = `<div id="availabilityMsg" style="margin:20%"><h3>No comments yet</h3></div>`
    }
}
