import { DropDown } from "./tools.js"
import { AddOrUpdatePost, CategoriesSelection } from "./newPost.js"

document.addEventListener("DOMContentLoaded", () => {
  const categories = document.querySelectorAll('input[id^=category]')

  categories.forEach(category => {
    if (selected.includes(category.value)) {
      category.checked = true
    }
  })

  DropDown() // there problem on that dropdown not working only on this page
  CategoriesSelection()
  const submitPost = document.getElementById('submitPost')

  submitPost.addEventListener('click', () => {
    AddOrUpdatePost(`/updatePost`)
  })
})
