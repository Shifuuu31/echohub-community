export { CheckClick, CheckLength }

const CheckClick = () => {
    const categoriesChecked = document.querySelectorAll('input[id^=category]:checked')
    if (categoriesChecked.length == 0) {
        alert(`Select at least one category`)
        return false;
    }
    return true;
}

const MAX_CATEGORIES = 3
const CheckLength = (category, checkedLength) => {
    if (checkedLength > MAX_CATEGORIES) {
        category.checked = false
        alert(`You can only select up to ${MAX_CATEGORIES} categories.`)
    }
}

