const selected = {{.Categories_selected }};
const categories = document.querySelectorAll('[id^="category_"]');
categories.forEach(category => {
    if (selected.includes(category.value)) {
        category.setAttribute("checked", "true")
    }
});