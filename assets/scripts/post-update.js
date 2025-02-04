import { AddUpdatePost, CategoriesSelection } from "./newPost.js";
import { DropDown } from "./tools.js"

document.addEventListener("DOMContentLoaded", () => {
    DropDown()
    CategoriesSelection()
    const submitPost = document.getElementById('submitPost')

    submitPost.addEventListener('click', () => {
        AddUpdatePost(`/updatePost`)
    })
})

const categories = document.querySelectorAll('input[id^=category]');
// check categories of post
categories.forEach(category => {
    if (selected.includes(category.value)) {
        category.setAttribute("checked", "true")
        const label = document.querySelector(`label[for="${category.id}"]`);
        label.style.color = 'white';
        label.style.backgroundColor = '#8552ff';
    }
});