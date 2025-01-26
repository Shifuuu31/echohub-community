export { Popup }

const popup = document.getElementById("popup")
const postsBtns = document.querySelectorAll("#commentBtn, .post-info-2")
const popupBackground = document.querySelector(".popup-background")
const closeButton = document.querySelector(".close")

const Popup = () => {
    // open popup
    postsBtns.forEach(postBtn => postBtn.addEventListener("click", () => {
        popupBackground.style.display = popup.style.display = "block"
    }))
    
    // close popup
    popupBackground.addEventListener("click", (event) => {
        if (event.target === popupBackground || event.target === closeButton) {
            popupBackground.style.display = popup.style.display = "none"

        }
    })
}