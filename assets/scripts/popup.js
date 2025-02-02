import { fetchResponse, AddComment } from "./tools.js"
export { Popup }


const commentsSection = document.getElementById("comments-section");
const Popup = () => {
    const popup = document.getElementById("popup")
    const popupBackground = document.getElementById("popup-background")
    const closeButton = document.querySelector(".close")

    const attachEventListeners = () => {
        const postsBtns = document.querySelectorAll("#commentBtn, #post-title")
        postsBtns.forEach(postBtn => {
            postBtn.removeEventListener("click", openPopup)
            postBtn.addEventListener("click", (event) => {

                openPopup(event)
            })
        })
    }


    const openPopup = async (event) => {
        let popupPost = document.querySelector("#popup #post")
        if (popup && popupBackground) {
            popupBackground.style.display = popup.style.display = "block"
            const targetedPost = event.target.closest("#post")
            const postID =targetedPost.getAttribute("post-id")
            popupPost.replaceWith(targetedPost.cloneNode(true))
            const cmntGrp = document.getElementById('comment-group')
            cmntGrp.innerHTML = `<textarea placeholder="Type a comment..." type="text" id="comment-field"></textarea>
                        <button class="new-comment" id="${postID}"><i class="fas fa-paper-plane"></i></button>`

            const newCmntBtn = document.getElementById(`${postID}`)

            newCmntBtn.addEventListener('click', async () => {
                const cmntTxt = document.getElementById('comment-field')
                const newCmnt = {
                    postid: postID,
                    userid: 1, // to be handled
                    content: cmntTxt.value,
                }
                cmntTxt.value = ''

                // await createComment(newCmnt)

            })

            await displayComments(postID);

        }
    }

    const displayComments = async (postid) => {
        commentsSection.innerHTML = ''
        let comments
        console.log(typeof postid)
        console.log(postid);
        try {
            const response = await fetchResponse(`/comments`, { ID: postid })
            if (response.status === 200) {
                console.log("comments recieved succesfully" )
                console.log(response.body)
                comments = response.body

            } else {
                console.log("Unexpected response:", response.body)
            }
    
        } catch (error) {
            console.error('Error during login process:', error)
        }
        comments.forEach(comment => {
            commentsSection.appendChild(AddComment(comment));
        });

    }

    const createComment = async (newCmnt) => {
        console.log(newCmnt.postid);
        await fetchResponse('/createComment', newCmnt)
        // await displayComments(newCmnt.postid)
    }


    const closePopup = (event) => {
        if (event.target === popupBackground || event.target === closeButton) {
            popupBackground.style.display = popup.style.display = "none"

        }
    }

    if (popupBackground) {
        popupBackground.addEventListener("click", closePopup)
    }

    return attachEventListeners
}


