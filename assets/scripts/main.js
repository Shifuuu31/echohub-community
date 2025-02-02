import { Popup } from "./popup.js"
import { displayPosts, DataToFetch } from "./display.js"
import { DropDown } from "./tools.js"

let attachPopupListeners = null

const initializePosts = async () => {
    attachPopupListeners = Popup()
    const postsAdded = await displayPosts()
    if (postsAdded) {
        attachPopupListeners()
    }
}

const handleScroll = async () => {
    const isAtBottom = window.innerHeight + window.scrollY >= document.body.offsetHeight

    if (isAtBottom) {
        const newPostsAdded = await displayPosts(DataToFetch.category, true)

        if (newPostsAdded) {
            attachPopupListeners()
        }
    }
}


const handleCategoryChange = async (event) => {
    const postsLoaded = await displayPosts(event.target.defaultValue)

    if (postsLoaded) {
        attachPopupListeners()
    }
}

// event listner for sort by category
const setupCategoryListeners = () => {
    const categories = document.querySelectorAll("input[id^=category]")
    categories.forEach(category => {
        category.addEventListener('change', handleCategoryChange)
    })
}
let isThrottled = true;
// event listner for scroll
const setupScrollListener = () => {
    window.addEventListener('scroll', () => {
        if (isThrottled) {
            isThrottled = false;

            setTimeout(() => {
                handleScroll();
                isThrottled = true;
            }, 100);
        }
    });
}

const initialize = async () => {
    try {
        DropDown()
        await initializePosts()
        setupCategoryListeners()
        setupScrollListener()
    } catch (error) {
        console.error('Failed to initialize application:', error)
    }
}

document.addEventListener("DOMContentLoaded", initialize)// Entry point