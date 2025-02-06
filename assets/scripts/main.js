import { closePopup, popupBackground } from "./popup.js"
import { displayPosts, DataToFetch } from "./display.js"
import { DropDown } from "./tools.js"

// let attachPopupListeners = null

const initPosts = async () => {
    // const attachPopupListeners = Popup()
    await displayPosts()
    // if (postsAdded) {
    //     attachPopupListeners()
    // }
}

const infiniteScroll = async () => {
    // const attachPopupListeners = Popup()
    const isAtBottom = window.innerHeight + window.scrollY >= document.body.offsetHeight

    if (isAtBottom) {
        await displayPosts(DataToFetch.category, true)

        // if (newPostsAdded) {
        //     attachPopupListeners()
        // }
    }
}


const CategoriesFilter = async (event) => {
    // const attachPopupListeners = Popup()
    await displayPosts(event.target.defaultValue)

    // if (postsLoaded) {
    //     attachPopupListeners()
    // }
}

// event listner for sort by category
const setupCategoryListeners = () => {
    const categories = document.querySelectorAll("input[id^=category]")
    categories.forEach(category => {
        category.addEventListener('change', CategoriesFilter)
    })
}
let isThrottled = true;
// event listner for scroll
const throttleScroll = () => {
    window.addEventListener('scroll', () => {
        if (isThrottled) {
            isThrottled = false;

            setTimeout(() => {
                infiniteScroll();
                isThrottled = true;
            }, 300);
        }
    });
}


if (popupBackground) {
        popupBackground.addEventListener("click", closePopup)
    }

const init = async () => {
    try {
        DropDown()
        await initPosts()
        setupCategoryListeners()
        throttleScroll()
    } catch (error) {
        console.error('Failed to init application:', error)
    }
}

document.addEventListener("DOMContentLoaded", init)// Entry point