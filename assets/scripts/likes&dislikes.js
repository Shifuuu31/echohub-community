export { setupLikeDislikeListner }
import { fetchResponse } from "./tools.js";

// like/dislike when click on them
async function handleLikeDislike(button, isLike) {
    const parent = button.parentElement
    const entityId = parent.getAttribute('data-entity-id');
    const entityType = parent.getAttribute('data-entity-type');
    const lastReaction = parent.getAttribute('data-reaction');

    const DataToFetch = {
        entityId: parseInt(entityId),
        entityType: entityType,
        liked: isLike,
    };

    try {
        const response = await fetchResponse(`/like-dislike`, DataToFetch)        
        if (response.status == 200) {
            console.log(`${DataToFetch.entityType} liked | disliked succefully`);
            const allEntities = document.querySelectorAll(`[data-entity-id="${entityId}"]`);

            allEntities.forEach(entity => {
                const likeBtn = entity.querySelector(".like-btn")
                const dislikeBtn = entity.querySelector(".dislike-btn")

            });
        } else {
            console.log("Unexpected response:")
            return
        }

    } catch (error) {
        console.error('Error:', error);
    }
}
const setupLikeDislikeListner = () => {
    const likeButtons = document.querySelectorAll('.like-btn');
    const dislikeButtons = document.querySelectorAll('.dislike-btn');
    likeButtons.forEach(button => {
        button.addEventListener('click', () => {
            handleLikeDislike(button, true)
        });
    });
    dislikeButtons.forEach(button => {
        button.addEventListener('click', () => {
            handleLikeDislike(button, false)
        });
    });
}
