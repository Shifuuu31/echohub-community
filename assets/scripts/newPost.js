import { AddOrUpdatePost, CategoriesSelection, DropDown } from "./tools.js"

document.addEventListener("DOMContentLoaded", () => {
    DropDown()
    CategoriesSelection()
    const submitPost = document.getElementById('submitPost')

    submitPost.addEventListener('click', () => {
        AddOrUpdatePost(`/addNewPost`)
    })
})

