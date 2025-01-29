export { update }

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

const update = async (url) => {
    let DataToFetch = {
        id: document.getElementsByClassName("wraper")[0].id,
        title: document.getElementById("post-title").value,
        content: document.getElementById("content").value,
        categories: [],
    };

    document.querySelectorAll("input[id^=category]:checked").forEach(category => { DataToFetch.categories.push(category.value) })
    await fetchData(url, DataToFetch)
    window.location.href = "/";
}