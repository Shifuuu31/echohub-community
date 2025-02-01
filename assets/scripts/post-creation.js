import { CheckClick, CheckLength } from "./tools.js";
import { update } from "./post-create-update.js";

const categories = document.querySelectorAll('input[id^=category]')
const submitBtn = document.getElementById('submit')
const urlcreate = `${window.location.origin}/createPost`;


categories.forEach((category) => {
    category.addEventListener('change', () => {
        const categoriesChecked = document.querySelectorAll('input[id^=category]:checked')
        CheckLength(category, categoriesChecked.length)
    })
});

// submitBtn.addEventListener('click', CheckClick)
submitBtn.addEventListener('click', async () => {
    if (CheckClick() == true) {
        await update(urlcreate)
    }
});