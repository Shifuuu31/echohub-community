export { Popup }

const Popup = () => {
    const popup = document.getElementById("popup")
    const popupBackground = document.getElementById("popup-background")
    const closeButton = document.querySelector(".close")

    const attachEventListeners = () => {
        const postsBtns = document.querySelectorAll("#commentBtn, #post-title")
        // console.log('popup post', popupPost)
        
        postsBtns.forEach(postBtn => {
            postBtn.removeEventListener("click", openPopup)
            postBtn.addEventListener("click", (event) => {
                
                openPopup(event)
            })
        })
    }
    
    const openPopup = (event) => {
        if (popup && popupBackground) {
            popupBackground.style.display = popup.style.display = "block"
            let popupPost = document.querySelector("#popup #post")
            const targetedPost = event.target.closest('#post')
                // console.log(targetedPost)
                popupPost.replaceWith(targetedPost.cloneNode(true))
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
