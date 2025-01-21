const MAX_CATEGORIES = 3
const categories = document.querySelectorAll('input[name="categoryElement"]')

for (let category of categories) {
    category.addEventListener('change', () => {
        const categories = document.querySelectorAll('input[name="categoryElement"]:checked')

        if (categories.length > MAX_CATEGORIES) {
            categories.checked = false
            alert(`You can only select up to ${MAX_CATEGORIES} categories.`)
        }
    })
}