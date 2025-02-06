import { closePopup, popupBackground } from "./popup.js"
import { displayPosts, DataToFetch } from "./display.js"
import { DropDown } from "./tools.js"


// event listner for filter by category
const CategoriesFilter = () => {
    const categories = document.querySelectorAll("input[id^=category]")
    categories.forEach(category => {
        category.addEventListener('change', async(event) => {
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
            isThrottled = false;

            setTimeout(async() => {
                if (window.innerHeight + window.scrollY >= document.body.offsetHeight) {
                    await displayPosts(DataToFetch.category, true)
                }
                isThrottled = true;
            }, 300);
        }
    });
}




const init = async () => {
    try {
        DropDown()
        await displayPosts()
        CategoriesFilter()
        throttleScroll()
    } catch (error) {
        console.error('Failed to init application:', error)
    }
}

document.addEventListener("DOMContentLoaded", init)// Entry point