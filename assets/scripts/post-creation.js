const MAX_CATEGORIES = 3
const categories = document.querySelectorAll('input[name="categoryElement"]')
const submitBtn = document.getElementById('submit')

for (let category of categories) {
    category.addEventListener('change', () => {
        const categoriesChecked = document.querySelectorAll('input[name="categoryElement"]:checked')
        if (categoriesChecked.length > MAX_CATEGORIES) {
            category.checked = false
            alert(`You can only select up to ${MAX_CATEGORIES} categories.`)
        }
    })
}

submitBtn.addEventListener('click', () => {
    const categoriesChecked = document.querySelectorAll('input[name="categoryElement"]:checked')
    if (categoriesChecked.length == 0) {
        alert(`Select at list one category`)
    }
})