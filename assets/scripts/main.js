import { Popup } from "./popup.js"
import { desplayPosts, DataToFetch } from "./fetch-posts.js"

let isLoading = false
let attachPopupListeners = null
// let DataToFetch = {}

const initializePosts = async () => {
    attachPopupListeners = Popup()
    const postsAdded = await desplayPosts()
    
    if (postsAdded) {
        attachPopupListeners()
    }
}

const handleScroll = async () => {
    const isAtBottom = window.innerHeight + window.scrollY >= document.body.offsetHeight
    
    if (isAtBottom && !isLoading) {
        isLoading = true
        const newPostsAdded = await desplayPosts(DataToFetch.category, true)
        
        if (newPostsAdded) {
            attachPopupListeners()
        }
        
        isLoading = false
    }
}

const handleCategoryChange = async (event) => {
    const postsLoaded = await desplayPosts(event.target.defaultValue)
    
    if (postsLoaded) {
        attachPopupListeners()
    }
}

const setupCategoryListeners = () => {
    const categories = document.querySelectorAll("input[id^=category]")
    categories.forEach(category => {
        category.addEventListener('change', handleCategoryChange)
    })
}

const setupScrollListener = () => {
    window.addEventListener('scroll', handleScroll)
}

const initialize = async () => {
    try {
        await initializePosts()
        setupCategoryListeners()
        setupScrollListener()
    } catch (error) {
        console.error('Failed to initialize application:', error)
    }
}

document.addEventListener("DOMContentLoaded", initialize)// Entry point