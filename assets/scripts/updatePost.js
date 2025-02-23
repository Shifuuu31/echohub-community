import { AddOrUpdatePost, CategoriesSelection, DropDown } from "./tools.js"

document.addEventListener("DOMContentLoaded", () => {
  const categories = document.querySelectorAll('input[id^=category]')

  categories.forEach(category => {
    if (selected.includes(category.value)) {
      category.checked = true
    }
  })

  DropDown() 
  CategoriesSelection()
  const submitPost = document.getElementById('submitPost')

  submitPost.addEventListener('click', () => {
    AddOrUpdatePost(`/updatingPost?ID=${document.querySelector('.wraper').id}`)
  })
})
