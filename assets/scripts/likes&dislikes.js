export { handleLikeDislike }
import { displayMsg, fetchResponse } from "./tools.js";

// like/dislike when click on them
async function handleLikeDislike(postData, isLike) {
    const btns = postData.querySelector("#buttons")
    const entityId = btns.getAttribute("entity-id")
    const entityType = btns.getAttribute("entity-type")

    const DataToFetch = {
        entityId: parseInt(entityId),
        entityType: entityType,
        liked: isLike,
    }

        const response = await fetchResponse(`/like-dislike`, DataToFetch)
        if (response.status == 200) {
            const likeBtn = btns.querySelector('.like-btn')
            const dislikeBtn = btns.querySelector('.dislike-btn')
            likeBtn.childNodes[1].textContent = response.body.extra[0]
            dislikeBtn.childNodes[1].textContent = response.body.extra[1]
            switch (response.body.messages[0]) {
                case "liked":
                    likeBtn.childNodes[0].src = '/assets/imgs/live-like.png'
                    dislikeBtn.childNodes[0].src = '/assets/imgs/dislike.png'
                    break;

                case "disliked":
                    likeBtn.childNodes[0].src = '/assets/imgs/like.png'
                    dislikeBtn.childNodes[0].src = '/assets/imgs/live-dislike.png'
                    break;

                default:
                    likeBtn.childNodes[0].src = '/assets/imgs/like.png'
                    dislikeBtn.childNodes[0].src = '/assets/imgs/dislike.png'
            }

        }
        else if (response.status == 403) {
            console.log('hshhs')
            console.log(response)
            displayMsg([response.body])
        }
        else if (response.status == 400) {
            displayMsg(['oops!'])
        } else {
            console.log("Unexpected response:")
            return
        }

}