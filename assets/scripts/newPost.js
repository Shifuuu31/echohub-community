import { fetchResponse, displayErr, DropDown } from "./tools.js"

document.addEventListener("DOMContentLoaded", () => {
    DropDown()
    CategoriesSelection()
    const submitPost = document.getElementById('submitPost')

    submitPost.addEventListener('click', AddNewPost)
})


const AddNewPost = async () => {
    const newPost = {
        title: document.getElementById('title').value,
        content: document.getElementById('content').value,
        selectedCategories: []
    }
    document.querySelectorAll('input[id^=category]:checked').forEach((selectedCategory) => {
        newPost.selectedCategories.push(selectedCategory.value)
    })

    console.log(newPost)


    const response = await fetchResponse(`/addNewPost`, newPost)

    if (response.status === 401) {
        console.log("Unauthorized: try to login")
    } else if (response.status === 400) {
        console.log(response.body);
        displayErr(response.body.messages)
    } else if (response.status === 200) {
        console.log("post added successfully")
        window.location.href = "/"
    } else {
        console.log("Unexpected response:", response.body)
    }
}


const CategoriesSelection = () => {
    const categories = document.querySelectorAll('input[id^=category]')

    categories.forEach((category) => {
        category.addEventListener('change', () => {
            const checkedCategories = document.querySelectorAll('input[id^=category]:checked')
            if (checkedCategories.length > 3) {
                category.checked = false
                displayErr(['You can only select up to 3 categories'])
            }
        })
    })
}