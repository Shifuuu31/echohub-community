import { CheckClick, CheckLength } from "./tools.js";



const categories = document.querySelectorAll('input[id^=category]')
const submitBtn = document.getElementById('submit')


categories.forEach((category) => {
    category.addEventListener('change', () => {
        const categoriesChecked = document.querySelectorAll('input[id^=category]:checked')
        CheckLength(category, categoriesChecked.length)
    })
});

submitBtn.addEventListener('click', CheckClick)