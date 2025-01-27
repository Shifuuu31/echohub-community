import { CheckClick, CheckLength } from "./tools.js";

const submitBtn = document.getElementById("submit")
const categories = document.querySelectorAll('input[id^="category"]')


const fetchData = async (url, obj) => {
    try {
        await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(obj),
        })
    } catch (e) {
        console.error("Error fetching data:", e.message);
    }
}

const update = async () => {
    const url = `${window.location.origin}/postUpdate`;

    let DataToFetch = {
        id: document.getElementsByClassName("wraper")[0].id,
        title: document.getElementById("post-title").value,
        content: document.getElementById("content").value,
        categories: [],
    };
    document.querySelectorAll("input[id^=category]:checked").forEach(category => { DataToFetch.categories.push(category.name) })

    console.log(typeof DataToFetch.id);

    await fetchData(url, DataToFetch)
    window.location.href = "/";
}


categories.forEach((category) => {
    category.addEventListener('change', () => {
        const categoriesChecked = document.querySelectorAll('input[id^=category]:checked')
        console.log('hihi', categoriesChecked);
        CheckLength(category, categoriesChecked.length)
    })
});

// eventlistner to button submit
submitBtn.addEventListener('click', async () => {
    // console.log(CheckClick() )
// console.log('hooo');

if (CheckClick() == true) {
        // console.log('hiiii');
        await update()

    }
});

// check categories of post
const selectedCategories = document.querySelectorAll('[id^=category]');
selectedCategories.forEach(category => {
    if (selected.includes(category.name)) {
        category.setAttribute("checked", "true")
    }
});