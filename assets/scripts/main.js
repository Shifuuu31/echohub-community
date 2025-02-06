import { closePopup, popupBackground } from "./popup.js"
import { displayPosts, DataToFetch } from "./display.js"
import { handleLikeDislike } from "./likes&dislikes.js"
import { DropDown } from "./tools.js"
export { setupLikeDislikeListner }

window.addEventListener("load", () => {
    setTimeout(() => {
        document.querySelector(".loader-container").style.display = "none";
        document.getElementById("navBar").style.display = "flex";
        document.getElementById("main").style.display = "block";
    }, 0)
})

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

// event listner for filter by category
const CategoriesFilter = () => {
    const categories = document.querySelectorAll("input[id^=category]")
    categories.forEach(category => {
        category.addEventListener('change', async (event) => {
            await displayPosts(event.target.defaultValue)

        })
    })
}

const ulCategories = document.getElementById("categories");

ulCategories.addEventListener("wheel", (event) => {
    event.preventDefault()
    ulCategories.scrollLeft += event.deltaY
})

if (popupBackground) {
    popupBackground.addEventListener("click", closePopup)
}

let isThrottled = true;
// event listner for scroll
const throttleScroll = () => {
    window.addEventListener('scroll', () => {
        if (isThrottled) {
            isThrottled = false

            setTimeout(async () => {
                if (window.innerHeight + window.scrollY >= document.body.offsetHeight) {
                    await displayPosts(DataToFetch.category, true)
                }
                isThrottled = true
            }, 300)
        }
    })
}

const init = async () => {
    try {
        DropDown()
        await displayPosts()
        CategoriesFilter()
        throttleScroll()
        setupLikeDislikeListner()
    } catch (error) {
        console.error('Failed to init application:', error)
    }
}

document.addEventListener("DOMContentLoaded", init)// Entry point