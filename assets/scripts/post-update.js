// console.log('jgjugjhgff');
console.log(selected);
const selectedCategories = document.querySelectorAll('[id^="category"]');
selectedCategories.forEach(category => {
    console.log(category.name);
    if (selected.includes(category.name)) {
        category.setAttribute("checked", "true")
    }
});