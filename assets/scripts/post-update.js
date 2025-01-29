import { CheckClick, CheckLength } from "./tools.js";
import { update } from "./post-create-update.js";

const categories = document.querySelectorAll('input[id^="category"]')
const submitBtn = document.getElementById("submit")
const urlUpdate = `${window.location.origin}/updatePost`;

categories.forEach((category) => {
    category.addEventListener('change', () => {
        const categoriesChecked = document.querySelectorAll('input[id^=category]:checked')
        CheckLength(category, categoriesChecked.length)
    })
});

// eventlistner to button submit
submitBtn.addEventListener('click', async () => {
    if (CheckClick() == true) {
        await update(urlUpdate)
    }
});


// check categories of post
// categories.forEach(category => {
//     if (selected.includes(category.value)) {
//         category.setAttribute("checked", "true")
//     }
// });

const selectedSet = new Set(selected);
categories.forEach(category => {
    if (selectedSet.has(category.value)) {
        category.setAttribute("checked", "true");
    }
});