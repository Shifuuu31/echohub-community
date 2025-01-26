export { Popup }

const Popup = () => {
    const postsBtns = document.querySelectorAll("#commentBtn, .post-info-2")
    const popupBackground = document.querySelector(".popup-background")
    const closeButton = document.querySelector("#popup-comment button")

    for (let i = 0; i < postsBtns.length; i++) {
        postsBtns[i].addEventListener("click", OpenPopup)
    }

    closeButton.addEventListener("click", ClosePopup)
    popupBackground.addEventListener("click", (event) => {
        if (event.target === popupBackground) {
            ClosePopup()
        }
    })
}

const OpenPopup = () => {
    document.querySelector(".popup-background").style.display = "block"
    document.getElementById("popup").style.display = "block"
    document.body.style.overflow = "hidden"
}

const ClosePopup = () => {
    document.querySelector(".popup-background").style.display = "none"
    document.getElementById("popup").style.display = "none"
    document.body.style.overflow = "auto"
}
