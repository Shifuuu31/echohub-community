export { Popup }

const Popup = () => {
    const popup = document.getElementById("popup")
    const popupBackground = document.getElementById("popup-background")
    const closeButton = document.querySelector(".close")

    const attachEventListeners = () => {
        const postsBtns = document.querySelectorAll("#commentBtn, #post-title")
        let popupPost = document.querySelector("#popup #post")
        console.log('popup post', popupPost)

        postsBtns.forEach(postBtn => {
            postBtn.removeEventListener("click", openPopup)
            postBtn.addEventListener("click", (event) => {
                let targetedPost = event.target.closest('#post')
                console.log(targetedPost)
                popupPost.replaceWith(targetedPost)
                openPopup()
            })
        })
    }

    const openPopup = () => {
        if (popup && popupBackground) {
            popupBackground.style.display = popup.style.display = "block"

        }
    }

    const closePopup = (event) => {
        if (event.target === popupBackground || event.target === closeButton) {
            popupBackground.style.display = popup.style.display = "none"

        }
    }

    if (popupBackground) {
        popupBackground.addEventListener("click", closePopup)
    }

    return attachEventListeners
}
